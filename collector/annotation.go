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
