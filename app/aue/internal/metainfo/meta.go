package metainfo

import (
	"fmt"
	"strings"

	key "github.com/Galdoba/ffstuff/app/aue/internal/define"
)

type Meta struct {
	Key   string
	Value string
}

func NewMeta(k, v string) Meta {
	return Meta{k, v}
}

type Collection interface {
	Add(...Meta) error
	Show(string) string
}

type metaCollection struct {
	metas map[string]Meta
}

func NewCollection() *metaCollection {
	mc := metaCollection{}
	mc.metas = make(map[string]Meta)
	return &mc
}

func (mc *metaCollection) Add(metas ...Meta) error {
	for _, incoming := range metas {
		if err := assertAcceptable(incoming.Key); err != nil {
			return err
		}
		stored, ok := mc.metas[incoming.Key]
		switch ok {
		case true:
			if strings.HasPrefix(incoming.Value, stored.Value) || strings.HasPrefix(stored.Value, incoming.Value) {
				continue
			}
			if stored.Value != incoming.Value {
				return fmt.Errorf("conflicting data '%v': [%v] != [%v]", incoming.Key, incoming.Value, stored.Value)
			}
		case false:
			mc.metas[incoming.Key] = incoming
		}
	}
	return nil
}

func assertAcceptable(metaKey string) error {
	switch metaKey {
	default:
		return fmt.Errorf("incoming data bad key '%v'", metaKey)
	case key.META_Base,
		key.META_Season,
		key.META_Episode,
		key.META_PRT,
		key.META_Audio_Layout_ + "0",
		key.META_Audio_Layout_ + "1",
		key.META_Audio_Lang_ + "0",
		key.META_Audio_Lang_ + "1":
	}
	return nil
}

func (mc *metaCollection) Show(metaKey string) string {
	if err := assertAcceptable(metaKey); err != nil {
		return ""
	}
	if _, ok := mc.metas[metaKey]; !ok {
		return ""
	}
	return mc.metas[metaKey].Value
}
