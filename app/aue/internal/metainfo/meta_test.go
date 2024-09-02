package metainfo

import (
	"testing"

	"github.com/Galdoba/ffstuff/app/aue/internal/define"
)

func TestAdd(t *testing.T) {
	col := NewCollection()
	base1 := NewMeta(define.META_Base, "Industry_s03e03_PRT240830124300_SER_05052_18")
	base2 := NewMeta(define.META_Base, "Industry_s03e03_PRT240830124300_SER_05052_18RUS")
	err := col.Add(base1, base2)
	if err != nil {
		t.Errorf("err: %v", err)
	}
}
