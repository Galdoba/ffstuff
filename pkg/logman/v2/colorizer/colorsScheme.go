package colorizer

type colorSchema struct {
	color256 map[colorKey]uint8
}

func New(customKeys ...colorData) *colorSchema {
	c := colorSchema{}
	c.color256 = make(map[colorKey]uint8)
	// c.color256 = defaultColorKey()
	for _, custom := range customKeys {
		c.color256[custom.key] = custom.val
	}
	return &c
}

func DefaultScheme() *colorSchema {
	c := colorSchema{}
	c.color256 = make(map[colorKey]uint8)
	c.color256 = defaultColorKey()
	return &c
}

func (c *colorSchema) WithColors(colors ...colorData) *colorSchema {
	for _, custom := range colors {
		c.color256[custom.key] = custom.val
	}
	return c
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

func defaultColorKey() map[colorKey]uint8 {
	colMap := make(map[colorKey]uint8)
	colMap[fgKey("base")] = 7
	colMap[bgKey("base")] = 0

	colMap[fgKey("string")] = 208
	colMap[fgKey("byte")] = 95
	colMap[fgKey("rune")] = 95
	colMap[fgKey("int")] = 120
	colMap[fgKey("int8")] = 120
	colMap[fgKey("int16")] = 120
	colMap[fgKey("int32")] = 120
	colMap[fgKey("int64")] = 120
	colMap[fgKey("float32")] = 9
	colMap[fgKey("float64")] = 9
	colMap[fgKey("bool")] = 12

	colMap[fgKey("struct")] = 221
	colMap[fgKey("slice")] = 14 //248
	colMap[fgKey("interface")] = 2
	colMap[fgKey("nil")] = 12
	colMap[fgKey("map")] = 14  //207
	colMap[fgKey("ptr")] = 221 //207
	colMap[fgKey("func")] = 36 //207
	colMap[fgKey("chan")] = 2  //207

	colMap[fldKey("fatal")] = 52
	colMap[fldKey("error")] = 88
	colMap[fldKey("warn")] = 184
	colMap[fldKey("info")] = 255
	colMap[fldKey("debug")] = 152
	colMap[fldKey("trace")] = 230
	return colMap
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
