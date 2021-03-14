package fldr

import (
	"fmt"
	"os"
	"strings"
)

type folder struct {
	Label          string
	Dyanmic        bool
	AddressFormula string
}

type Folder interface {
	Address() string
	Make() error
}

/*
fldr.New().SetName("Name").SetDynamic()

fldr.New(
	f.Name("Name"),
	f.Dynamic(),
)




*/

const (
	Label          = 0
	Dynamic        = 1
	AddressFormula = 2
)

func New(fName string, instructions ...folder) Folder {
	f := folder{}
	f.Label = fName
	for _, instr := range instructions {
		if instr.AddressFormula != "" {
			f.AddressFormula = instr.AddressFormula
		}
		f.Dyanmic = instr.Dyanmic
	}
	return &f
}

func Set(field int, val interface{}) folder {
	f := folder{}
	switch field {
	default:
		return f
	case Label:
		switch val.(type) {
		case string:
			f.Label = val.(string)
		}
	case AddressFormula:
		switch val.(type) {
		case string:
			f.AddressFormula = val.(string)
		}
	case Dynamic:
		switch val.(type) {
		case bool:
			f.Dyanmic = val.(bool)
		}
	}
	return f
}

func (f *folder) Make() error {
	info, err := os.Stat(f.AddressFormula)
	if err == nil && info.IsDir() {
		return nil
	}
	prts := strings.Split(f.AddressFormula, "\\")
	for i := 1; i <= len(prts); i++ {
		tempPath := strings.Join(prts[0:i], "\\")
		_, err := os.Stat(tempPath)
		if os.IsExist(err) {
			continue
		}
		err = os.Mkdir(tempPath, 0700)
		if tempPath == f.AddressFormula {
			fmt.Println("Directory created:", tempPath)
		}
	}
	return nil
}

// func (f *folder) SetAddress() string {
// 	return f.AddressFormula
// }

func (f *folder) Address() string {
	return f.AddressFormula
}
