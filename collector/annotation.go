// Copyright 2015 The Prometheus Authors
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

package collector

type Annotator interface {
	Netdev(name string) Labels
}

type Label struct {
	Key   string
	Value string
}

func SetAnnotator(a Annotator) {
	annotator = a
}

type Labels []Label

func (l Labels) keys() []string {
	ret := make([]string, 0, len(l))

	for _, la := range l {
		ret = append(ret, la.Key)
	}

	return ret
}

func (l Labels) values() []string {
	ret := make([]string, 0, len(l))

	for _, la := range l {
		ret = append(ret, la.Value)
	}

	return ret
}

type dummyAnnotator struct{}

func (d *dummyAnnotator) Netdev(name string) Labels {
	return nil
}
