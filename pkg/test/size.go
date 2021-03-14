package main

import (
	"github.com/Galdoba/ffstuff/fldr"
	"github.com/Galdoba/utils"
)

func main() {
	inFolder := fldr.New("InFolder",
		fldr.Set(fldr.AddressFormula, "d:\\SENDER\\DONE_"+utils.DateStamp()),
		fldr.Set(fldr.Dynamic, true),
	)
	inFolder.Make()
}
