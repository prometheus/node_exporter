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
	"go/build"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{}

	var matchCmd = &cobra.Command{
		Use:   "match [file]",
		Short: "Check whether the file matches the context.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// For debugging purposes, allow overriding these.
			goos, found := os.LookupEnv("GOHOSTOS")
			if !found {
				goos = runtime.GOOS
			}
			goarch, found := os.LookupEnv("GOARCH")
			if !found {
				goarch = runtime.GOARCH
			}
			ctx := build.Context{
				GOOS:   goos,
				GOARCH: goarch,
			}
			abs, err := filepath.Abs(args[0])
			if err != nil {
				panic(err)
			}
			match, err := ctx.MatchFile(filepath.Dir(abs), filepath.Base(abs))
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
			if match {
				os.Exit(0)
			}
			os.Exit(1)
		},
	}

	rootCmd.AddCommand(matchCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
