package colorizer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gookit/color"
)

type colorScheme struct {
	color256 map[colorKey]uint8
}

func New(customKeys ...colorData) *colorScheme {
	c := colorScheme{}
	c.color256 = make(map[colorKey]uint8)
	// c.color256 = defaultColorKey()
	for _, custom := range customKeys {
		c.color256[custom.key] = custom.val
	}
	return &c
}

func DefaultScheme() *colorScheme {
	c := colorScheme{}
	c.color256 = make(map[colorKey]uint8)
	c.color256 = defaultColorKey()
	return &c
}

func (c *colorScheme) WithAdditional(colors ...colorData) *colorScheme {
	for _, custom := range colors {
		c.color256[custom.key] = custom.val
	}
	return c
}

type colorKey struct {
	keytype string //field/fg/bg
	value   string
}

const (
	FIELD_KEY = "field"
	FG_KEY    = "fg"
	BG_KEY    = "bg"
)

// fieldKey - дает цвет покраски для поля (например '[error]' - может быть полностью красным в независимости от типа переменной)
func fldKey(val string) colorKey {
	return colorKey{
		keytype: FIELD_KEY,
		value:   val,
	}
}

func fgKey(val string) colorKey {
	return colorKey{
		keytype: FG_KEY,
		value:   val,
	}
}

func bgKey(val string) colorKey {
	return colorKey{
		keytype: BG_KEY,
		value:   val,
	}
}

func NewKey(keyType, value string) colorKey {
	return colorKey{
		keytype: keyType,
		value:   value,
	}
}

func defaultColorKey() map[colorKey]uint8 {
	colMap := make(map[colorKey]uint8)
	colMap[fgKey("string")] = 208
	colMap[fgKey("byte")] = 95
	colMap[fgKey("rune")] = 95
	colMap[fgKey("int")] = 255
	colMap[fgKey("int8")] = 255
	colMap[fgKey("int16")] = 255
	colMap[fgKey("int32")] = 255
	colMap[fgKey("int64")] = 255
	colMap[fgKey("float32")] = 203
	colMap[fgKey("float64")] = 203
	colMap[fgKey("bool")] = 38
	colMap[fgKey("struct")] = 246

	colMap[fgKey("slice")] = 246
	colMap[fldKey("fatal")] = 52
	colMap[fldKey("error")] = 88
	colMap[fldKey("warn")] = 184
	colMap[fldKey("debug")] = 152
	colMap[fldKey("trace")] = 230
	return colMap
}

type colorData struct {
	key colorKey
	val uint8
}

func CustomColor(key colorKey, color256Value uint8) colorData {
	return colorData{
		key: key,
		val: color256Value,
	}
}

func (c *colorScheme) Colorize(arg interface{}) string {
	return colorizeComplex(c, arg)
}

type coloredField struct {
	text  string
	fType string
	fg    uint8
	bg    uint8
}

var depth int

func constructFields(argVal reflect.Value) []*coloredField {
	depth++
	if depth > 50 {
		panic(1)
	}
	cFlds := []*coloredField{}

	cFlds = append(cFlds, cField(argVal))
	kind := fmt.Sprintf("%v", argVal.Kind())
	switch kind {
	default:
		fmt.Println(depth, "default", argVal.Type().String(), kind, argVal.Type())
		fmt.Println("REPEAT", argVal.Type().String())
		for i := 1; i <= argVal.NumField(); i++ {
			fld := argVal.Field(i - 1)
			cFlds = append(cFlds, constructFields(fld)...)
		}
	case "slice":
		for i := 0; i < argVal.Len(); i++ {
			rng := argVal.Index(i)
			cFlds = append(cFlds, constructFields(rng)...)
		}
	case "string", "bool",
		"int", "int8", "int16", "int32", "int64",
		"Int",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
	}
	return cFlds
}

func cField(arg reflect.Value) *coloredField {
	kind := fmt.Sprintf("%v", arg.Kind())
	text := fmt.Sprintf("%v", arg)

	switch kind {
	case "slice":
		text = arg.Type().String() + "{"
	case "struct":
		text = arg.Type().String() + "{"
		fmt.Println(kind, text, "+++++++")
	}

	return &coloredField{text, kind, 0, 0}
}

func colorizeComplex(c *colorScheme, arg interface{}) string {
	colorVals := []uint8{}

	argVal := reflect.ValueOf(arg)
	flds := constructFields(argVal)
	for i, fl := range flds {
		fmt.Println(i, fl)
	}
	argType := argVal.Type()
	argTypeString := argType.String()
	kind := fmt.Sprintf("%v", argVal.Kind())
	s := ""
	switch kind {
	default:
		//fmt.Println("kind:", argVal.Kind(), argVal)
		fg := colorOfkind(c, argVal, FG_KEY)
		bg := colorOfkind(c, argVal, BG_KEY)
		s := colorizeBasic(fmt.Sprintf("%v", argVal), fg, bg)
		return s
	case "slice":
		fmt.Println("kind SLICE:", argVal.Kind(), argTypeString)
		colorVals = append(colorVals, colorOfkind(c, argVal, FG_KEY))
		colorVals = append(colorVals, colorOfkind(c, argVal, BG_KEY))
		col := color.S256(colorVals...)
		s = col.Sprintf("%v{", argTypeString)
		for i := 0; i < argVal.Len(); i++ {
			rng := argVal.Index(i)
			//	fmt.Println(rng.Type())
			fmt.Println("kind", i, rng.Kind(), "val", fmt.Sprintf("%v", rng))
			fg := colorOfkind(c, rng, FG_KEY)
			bg := colorOfkind(c, rng, BG_KEY)

			s += colorizeBasic(fmt.Sprintf("%v", rng), fg, bg) + " "
			//s += c.Colorize(rng.Type()) + " "
		}
		s = strings.TrimSpace(s) + col.Sprintf("}")
		return s
	case "struct":

		fmt.Println("kind STRUCT:", argVal.Kind(), argVal)
		// colorVals = append(colorVals, colorOfkind(c, argVal, FG_KEY))
		// colorVals = append(colorVals, colorOfkind(c, argVal, BG_KEY))
		// col := color.S256(colorVals...)
		// s = col.Sprintf("%v{", argTypeString)

		for i := 1; i <= argType.NumField(); i++ {
			fld := argVal.Field(i - 1)
			fmt.Println("kind", i, fld.Kind(), "val", fmt.Sprintf("%v", fld))
			// kind := fmt.Sprintf("%v", fld.Kind())
			// switch kind {
			// case "slice":
			// 	fmt.Println("me slice", fld, fld.Kind())
			// 	sl := c.Colorize(fld)
			// 	fmt.Println(sl)
			// 	s += c.Colorize(fld)
			// default:
			// 	//		s += color.C256(c.color256[fgKey(fld.Type().String())]).Sprintf("%v", fld)
			// 	fg := colorOfkind(c, argVal, FG_KEY)
			// 	bg := colorOfkind(c, argVal, BG_KEY)
			// 	s += colorizeBasic(fmt.Sprintf("%v", argVal), fg, bg)
			// }

			//s += c.Colorize(fld)
			if i < argType.NumField() {
				s += " "
			}
		}
		// s = strings.TrimSpace(s) + col.Sprintf("}")
		return s

	}

	return "do not"
}

func (c *colorScheme) getColor(key colorKey) uint8 {
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

func typeValueGroup(valType reflect.Type) string {
	if valType.PkgPath() != "" {
		return "struct"
	}
	if strings.HasPrefix(valType.String(), "[]") {
		return "slice"
	}
	if strings.HasPrefix(valType.String(), "map[") {
		return "map"
	}
	switch valType.String() {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "bool", "string", "float32", "float64":
		return "primitive"
	}
	return "huh?"
}

type str1 struct {
	s   string
	i   int
	sl  []float64
	str str2
}

type str2 struct {
	b  bool
	fl float64
}

func testStr() str1 {
	return str1{
		s:  "string 1",
		i:  10,
		sl: []float64{15.1, 16.2},
		str: str2{
			b:  false,
			fl: 3.14,
		},
	}
}

func colorOfkind(c *colorScheme, val reflect.Value, keyType string) uint8 {
	kind := fmt.Sprintf("%v", val.Kind())
	switch kind {
	case "string", "bool",
		"int", "int8", "int16", "int32", "int64",
		"Int",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return c.getColor(NewKey(keyType, kind))
	case "struct":

		return c.getColor(NewKey(keyType, kind))
	case "slice":
		//fmt.Println("slice color", c.getColor(NewKey(keyType, kind)))
		return c.getColor(NewKey(keyType, kind))
	}
	switch keyType {
	case FG_KEY:
		return 7
	}
	return 0
}

func colorizeBasic(text string, fg, bg uint8) string {
	return color.S256(fg, bg).Sprintf(text)
}
