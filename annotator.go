package main

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/prometheus/node_exporter/collector"
	"gopkg.in/yaml.v2"
)

type Annotator struct {
	cfg *AnnotatorCfg
}

type AnnotatorCfg struct {
	Netdev []AnnotatorCfgNetdev `yaml:"netdev"`
}

type AnnotatorCfgNetdev struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels"`
}

func NewAnnotator(fp string) (*Annotator, error) {
	cfg, err := getAnnotatorConfig(fp)
	if err != nil {
		return nil, err
	}

	return &Annotator{
		cfg: cfg,
	}, nil
}

func getAnnotatorConfig(fp string) (*AnnotatorCfg, error) {
	cfg := &AnnotatorCfg{}

	fc, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to read annotations config")
	}

	err = yaml.Unmarshal(fc, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Unmarshal failed")
	}

	return cfg, nil
}

func (a *Annotator) Netdev(name string) collector.Labels {
	if a.cfg.Netdev == nil {
		return nil
	}

	for _, x := range a.cfg.Netdev {
		if name != x.Name {
			continue
		}

		return labelMapToLabelSlice(x.Labels)
	}

	return nil
}

func labelMapToLabelSlice(m map[string]string) collector.Labels {
	ret := make(collector.Labels, 0, len(m))

	for k, v := range m {
		ret = append(ret, collector.Label{
			Key:   k,
			Value: v,
		})
	}

	return ret
}
