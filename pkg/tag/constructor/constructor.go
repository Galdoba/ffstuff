package constructor

import "fmt"

type TagAssembler interface {
	Assemble(...interface{}) (string, string, error)
}

type tagAssembler struct {
	assemble AssembleFunc
}

func (c *tagAssembler) Assemble(input ...interface{}) (string, string, error) {
	return c.assemble(input...)
}

type AssembleFunc func(...interface{}) (string, string, error)

func NewAssembleFunc(fn func(...interface{}) (string, string, error)) AssembleFunc {
	return AssembleFunc(fn)
}

func Default() *tagAssembler {
	ta := tagAssembler{}
	ta.assemble = defaultAssemblerFunc
	return &ta
}

func defaultAssemblerFunc(input ...interface{}) (string, string, error) {
	output := make([]*string, 2)
	if len(input) != 2 {
		return *output[0], *output[1], fmt.Errorf("default constructor expects 2 values as input")
	}
	for i, data := range input {
		if err := assertInput(data, output[i]); err != nil {
			return *output[0], *output[1], err
		}
	}
	return *output[0], *output[1], nil
}

func assertInput(input interface{}, val *string) error {
	switch input.(type) {
	default:
		return fmt.Errorf("value %v is not a string", input)
	case string:
		value := input.(string)
		val = &value
	}
	return nil
}

func WithAssemblerFunc(afn AssembleFunc) *tagAssembler {
	ta := tagAssembler{}
	ta.assemble = afn
	return &ta
}
