// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type metric struct {
	name      string
	help      string
	subsystem string
	file      string
	labels    []string
}

type collectorInfo struct {
	name        string
	files       []string
	metrics     []metric
	flags       []flagInfo
	dataSources []string
	platforms   []string
}

type flagInfo struct {
	name        string
	description string
	defaultVal  string
}

func main() {
	collectorDir := "collector"
	fset := token.NewFileSet()
	collectors := make(map[string]*collectorInfo)

	absCollectorDir, _ := filepath.Abs(collectorDir)
	fmt.Fprintf(os.Stderr, "Searching in %s\n", absCollectorDir)

	err := filepath.Walk(absCollectorDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", path, err)
			return nil
		}

		baseName := filepath.Base(path)
		collectorName := strings.Split(baseName, "_")[0]
		collectorName = strings.TrimSuffix(collectorName, ".go")

		if _, ok := collectors[collectorName]; !ok {
			collectors[collectorName] = &collectorInfo{
				name: collectorName,
			}
		}
		c := collectors[collectorName]
		c.files = append(c.files, path)

		// Extract metadata from file comments
		for _, cg := range f.Comments {
			for _, cmt := range cg.List {
				text := strings.TrimPrefix(cmt.Text, "//")
				text = strings.TrimSpace(text)
				if strings.HasPrefix(text, "Data Sources:") {
					sources := strings.TrimPrefix(text, "Data Sources:")
					for _, s := range strings.Split(sources, ",") {
						c.dataSources = append(c.dataSources, strings.TrimSpace(s))
					}
				}
				if strings.HasPrefix(text, "Platforms:") {
					platforms := strings.TrimPrefix(text, "Platforms:")
					for _, p := range strings.Split(platforms, ",") {
						c.platforms = append(c.platforms, strings.TrimSpace(p))
					}
				}
				if strings.HasPrefix(text, "Metric:") {
					mstr := strings.TrimPrefix(text, "Metric:")
					parts := strings.SplitN(mstr, "-", 2)
					if len(parts) == 2 {
						c.metrics = append(c.metrics, metric{
							name: strings.TrimSpace(parts[0]),
							help: strings.TrimSpace(parts[1]),
							file: path,
						})
					}
				}
			}
		}

		// First pass: find all constant strings that might be subsystems
		subsystems := make(map[string]string)
		ast.Inspect(f, func(n ast.Node) bool {
			if spec, ok := n.(*ast.ValueSpec); ok {
				for i, name := range spec.Names {
					if strings.HasSuffix(name.Name, "Subsystem") || name.Name == "subsystem" {
						if len(spec.Values) > i {
							if lit, ok := spec.Values[i].(*ast.BasicLit); ok {
								subsystems[name.Name] = strings.Trim(lit.Value, "\"")
							}
						}
					}
				}
			}
			return true
		})

		// Second pass: find NewDesc and kingpin.Flag calls
		ast.Inspect(f, func(n ast.Node) bool {
			// Handle ValueSpec for kingpin flags defined at package level
			if vs, ok := n.(*ast.ValueSpec); ok {
				for _, val := range vs.Values {
					f := extractFlag(val)
					if f != nil {
						c.flags = append(c.flags, *f)
					}
				}
			}

			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if sel.Sel.Name != "NewDesc" {
				return true
			}

			if len(call.Args) < 2 {
				return true
			}

			m := metric{file: path}

			// Extract Help string
			if lit, ok := call.Args[1].(*ast.BasicLit); ok {
				m.help = strings.Trim(lit.Value, "\"")
			} else if ident, ok := call.Args[1].(*ast.Ident); ok {
				// Search for the identifier's value
				ast.Inspect(f, func(n2 ast.Node) bool {
					if spec, ok := n2.(*ast.ValueSpec); ok {
						for i, name := range spec.Names {
							if name.Name == ident.Name {
								if len(spec.Values) > i {
									if lit, ok := spec.Values[i].(*ast.BasicLit); ok {
										m.help = strings.Trim(lit.Value, "\"")
									}
								}
							}
						}
					}
					return true
				})
			}

			// Extract Name
			if nameCall, ok := call.Args[0].(*ast.CallExpr); ok {
				// Case: prometheus.BuildFQName(namespace, subsystem, "name")
				if len(nameCall.Args) >= 3 {
					subsystem := ""
					if lit, ok := nameCall.Args[1].(*ast.BasicLit); ok {
						subsystem = strings.Trim(lit.Value, "\"")
					} else if ident, ok := nameCall.Args[1].(*ast.Ident); ok {
						subsystem = subsystems[ident.Name]
					} else if sel, ok := nameCall.Args[1].(*ast.SelectorExpr); ok {
						subsystem = sel.Sel.Name
					}

					if lit, ok := nameCall.Args[2].(*ast.BasicLit); ok {
						namePart := strings.Trim(lit.Value, "\"")
						if subsystem != "" {
							m.name = fmt.Sprintf("node_%s_%s", subsystem, namePart)
						} else {
							m.name = fmt.Sprintf("node_%s", namePart)
						}
					}
				}
			} else if lit, ok := call.Args[0].(*ast.BasicLit); ok {
				m.name = strings.Trim(lit.Value, "\"")
			}

			// Extract Labels
			if len(call.Args) >= 3 {
				if comp, ok := call.Args[2].(*ast.CompositeLit); ok {
					for _, elt := range comp.Elts {
						if lit, ok := elt.(*ast.BasicLit); ok {
							m.labels = append(m.labels, strings.Trim(lit.Value, "\""))
						}
					}
				}
			}

			if m.name != "" && m.help != "" {
				c.metrics = append(c.metrics, m)
			}

			return true
		})

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// De-duplicate everything
	for _, c := range collectors {
		c.dataSources = uniqueStrings(c.dataSources)
		c.platforms = uniqueStrings(c.platforms)
		c.flags = uniqueFlags(c.flags)
		c.metrics = uniqueMetrics(c.metrics)
		sort.Slice(c.metrics, func(i, j int) bool {
			return c.metrics[i].name < c.metrics[j].name
		})
	}

	// Output global METRICS.md
	metricsFile := filepath.Join("docs", "METRICS.md")
	mf, err := os.Create(metricsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", metricsFile, err)
		os.Exit(1)
	}
	defer mf.Close()

	fmt.Fprintln(mf, "# Node Exporter Metrics")
	fmt.Fprintln(mf, "\nThis file is auto-generated by `tools/doc_generator/main.go`.")
	fmt.Fprintln(mf, "\n| Metric | Description | Collector |")
	fmt.Fprintln(mf, "| --- | --- | --- |")

	keys := make([]string, 0, len(collectors))
	for k := range collectors {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		c := collectors[k]
		for _, m := range c.metrics {
			fmt.Fprintf(mf, "| %s | %s | %s |\n", m.name, m.help, c.name)
		}
	}

	// Ensure docs/collectors directory exists
	os.MkdirAll("docs/collectors", 0755)

	// Output per-collector files
	for _, k := range keys {
		c := collectors[k]
		filename := filepath.Join("docs", "collectors", c.name+".md")
		f, err := os.Create(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", filename, err)
			continue
		}

		fmt.Fprintf(f, "# %s collector\n\n", c.name)
		fmt.Fprintf(f, "The %s collector exposes metrics about %s.\n\n", c.name, c.name)

		if len(c.platforms) > 0 {
			fmt.Fprintf(f, "## Supported Platforms\n\n- %s\n\n", strings.Join(c.platforms, "\n- "))
		}

		if len(c.dataSources) > 0 {
			fmt.Fprintf(f, "## Data Sources\n\n- %s\n\n", strings.Join(c.dataSources, "\n- "))
		}

		if len(c.flags) > 0 {
			fmt.Fprintf(f, "## Configuration Flags\n\n| Flag | Description | Default |\n| --- | --- | --- |\n")
			for _, fl := range c.flags {
				fmt.Fprintf(f, "| %s | %s | %s |\n", fl.name, fl.description, fl.defaultVal)
			}
			fmt.Fprintf(f, "\n")
		}

		if len(c.metrics) > 0 {
			fmt.Fprintf(f, "## Metrics\n\n| Metric | Description | Labels |\n| --- | --- | --- |\n")
			for _, m := range c.metrics {
				labels := "n/a"
				if len(m.labels) > 0 {
					labels = strings.Join(m.labels, ", ")
				}
				fmt.Fprintf(f, "| %s | %s | %s |\n", m.name, m.help, labels)
			}
		}

		f.Close()
	}
}

func extractFlag(expr ast.Expr) *flagInfo {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	if sel.Sel.Name == "Flag" {
		if x, ok := sel.X.(*ast.Ident); ok && x.Name == "kingpin" {
			if len(call.Args) >= 2 {
				name, _ := call.Args[0].(*ast.BasicLit)
				desc, _ := call.Args[1].(*ast.BasicLit)
				if name != nil && desc != nil {
					return &flagInfo{
						name:        strings.Trim(name.Value, "\""),
						description: strings.Trim(desc.Value, "\""),
					}
				}
			}
		}
		return nil
	}

	// Recurse into the chain
	f := extractFlag(sel.X)
	if f != nil {
		if sel.Sel.Name == "Default" && len(call.Args) > 0 {
			if d, ok := call.Args[0].(*ast.BasicLit); ok {
				f.defaultVal = strings.Trim(d.Value, "\"")
			}
		}
	}
	return f
}

func uniqueStrings(s []string) []string {
	m := make(map[string]bool)
	res := []string{}
	for _, v := range s {
		if !m[v] {
			m[v] = true
			res = append(res, v)
		}
	}
	sort.Strings(res)
	return res
}

func uniqueFlags(flags []flagInfo) []flagInfo {
	m := make(map[string]bool)
	res := []flagInfo{}
	for _, f := range flags {
		if !m[f.name] {
			m[f.name] = true
			res = append(res, f)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].name < res[j].name
	})
	return res
}

func uniqueMetrics(metrics []metric) []metric {
	m := make(map[string]bool)
	res := []metric{}
	for _, mt := range metrics {
		if !m[mt.name] {
			m[mt.name] = true
			res = append(res, mt)
		}
	}
	return res
}
