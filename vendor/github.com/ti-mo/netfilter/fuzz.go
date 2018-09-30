//+build gofuzz

package netfilter

func Fuzz(data []byte) int {

	attrs, err := UnmarshalAttributes(data)
	if err != nil {
		if attrs != nil {
			panic("attrs != nil on error")
		}
		return 0
	}

	_, err = MarshalAttributes(attrs)
	if err != nil {
		panic("error re-marshaling Attributes (!?)")
	}

	return 1
}
