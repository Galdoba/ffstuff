package colorizer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gookit/color"
)

// Colorize - takes any argument and return it's string representation in color.
func (c *colorSchema) Colorize(arg interface{}) string {
	return colorize(c, arg)
}

func colorize(c *colorSchema, arg interface{}) string {
	s := ""
	argVal := reflect.ValueOf(arg)
	flds := constructFields(argVal)
	for i, fl := range flds {
		if i == len(flds)-1 {
			fl.text = strings.TrimSpace(fl.text)
		}
		fl.text = strings.ReplaceAll(fl.text, "interface {}", "interface{}")
		fl.fg = colorToField(c, fl.fType, FG_KEY)
		fl.bg = colorToField(c, fl.fType, BG_KEY)
		s += color.S256(fl.fg, fl.bg).Sprint(fl.text)
	}
	return s
}

type coloredField struct {
	text  string
	fType string
	fg    uint8
	bg    uint8
}

func constructFields(argVal reflect.Value) []*coloredField {
	cFlds := []*coloredField{}
	cFlds = append(cFlds, cField(argVal))
	kind := argVal.Kind().String()
	switch kind {
	default:
		cFlds[0].text = "{"
		for i := 0; i < argVal.NumField(); i++ {
			fld := argVal.Field(i)
			vals := constructFields(fld)
			cFlds = append(cFlds, vals...)
		}
		cFlds[len(cFlds)-1].text = strings.TrimSuffix(cFlds[len(cFlds)-1].text, " ")
		cFlds = append(cFlds, &coloredField{"} ", "struct", 0, 0})
	case "invalid":
		cFlds[0] = nilField()
	case "ptr":
		switch argVal.IsNil() {
		case true:
			cFlds[0] = nilField()
		case false:
			cFlds[0].text = fmt.Sprintf("&{")
			for i := 0; i < argVal.Elem().NumField(); i++ {
				cFlds = append(cFlds, constructFields(argVal.Elem().Field(i))...)
			}
			cFlds[len(cFlds)-1].text = strings.TrimSuffix(cFlds[len(cFlds)-1].text, " ")
			cFlds = append(cFlds, &coloredField{"} ", kind, 0, 0})
		}
	case "func":
		switch argVal.IsNil() {
		case true:
			cFlds[0] = nilField()
		case false:
		}
	case "chan":
		switch argVal.IsNil() {
		case true:
			cFlds[0] = nilField()
		case false:
			cFlds[0].text = fmt.Sprintf("%v", argVal)
			cFlds[0].fType = kind
		}
	case "interface":
		interfaceType := argVal.Type().String()
		switch interfaceType {
		case "interface {}":
			switch argVal.IsNil() {
			case true:
				cFlds[0] = nilField()
			case false:
				cFlds = cFlds[1:]
				cFlds = append(cFlds, constructFields(argVal.Elem())...)
			}
		default:
			cFlds[0].text = fmt.Sprintf("<%v>", argVal)
		}
	case "slice":
		sliceType := strings.TrimPrefix(argVal.Type().String(), "[]")
		switch argVal.IsNil() {
		case true:
			cFlds[0].text = fmt.Sprintf("[] ")
			cFlds[0].fType = sliceType
		case false:
			cFlds[0].fType = sliceType
			for i := 0; i < argVal.Len(); i++ {
				rng := argVal.Index(i)
				slFld := constructFields(rng)
				for _, fld := range slFld {
					fld.fType = sliceType
					if strings.HasSuffix(fld.text, "{") {
						fld.text = "{"
					}
				}
				cFlds = append(cFlds, slFld...)
			}
			cFlds[len(cFlds)-1].text = strings.TrimSuffix(cFlds[len(cFlds)-1].text, " ")
			cFlds = append(cFlds, &coloredField{"]", sliceType, 0, 0})
		}
	case "map":
		switch argVal.IsNil() {
		case true:
			cFlds[0] = nilField()
		case false:
			iter := argVal.MapRange()
			for iter.Next() {
				k := iter.Key()
				v := iter.Value()
				cFlds = append(cFlds, constructFields(k)...)
				cFlds[len(cFlds)-1].text = strings.TrimSuffix(cFlds[len(cFlds)-1].text, " ")
				cFlds = append(cFlds, &coloredField{":", "map", 0, 0})
				cFlds = append(cFlds, constructFields(v)...)
			}
			cFlds[len(cFlds)-1].text = strings.TrimSuffix(cFlds[len(cFlds)-1].text, " ")
		}
		cFlds = append(cFlds, &coloredField{"] ", "map", 0, 0})
	case "string", "bool",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		cFlds[0].text += " "
	}
	return cFlds
}

func nilField() *coloredField {
	return &coloredField{"<nil> ", "nil", 0, 0}
}

func cField(arg reflect.Value) *coloredField {
	kind := fmt.Sprintf("%v", arg.Kind())
	text := fmt.Sprintf("%v", arg)
	switch kind {
	case "slice":
		text = "["
	case "struct":
		text = "{"
	case "map":
		text = "map["
	}
	return &coloredField{text, kind, 0, 0}
}

func (c *colorSchema) getColor(key colorKey) uint8 {
	if v, ok := c.color256[key]; ok {
		return v
	}
	switch key.keytype {
	case FG_KEY:
		return 7
	case BG_KEY:
		return 0
	}
	return 10
}

func colorToField(c *colorSchema, kind, keyType string) uint8 {
	switch kind {
	case "string", "bool",
		"int", "int8", "int16", "int32", "int64",
		"Int",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return c.getColor(NewKey(keyType, kind))
	case "struct", "map", "slice", "interface", "ptr", "func", "chan", "nil":
		return c.getColor(NewKey(keyType, kind))
	}
	switch keyType {
	case FG_KEY:
		return c.getColor(NewKey(FG_KEY, "base"))
	}
	return c.getColor(NewKey(BG_KEY, "base"))
}
