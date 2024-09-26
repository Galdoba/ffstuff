package copyprocess

import (
	"fmt"

	"github.com/Galdoba/ffstuff/app/grabber/internal/origin"
)

const (
	Status_Pending = 0
)

type copyActionState struct {
	sourcePaths    []origin.Origin
	destination    string
	markerExt      string
	creationErrors []error
	deleteAll      bool
	deleteMarker   bool
}

func NewCopyAction(opts ...Option) *copyActionState {
	cas := copyActionState{}
	// cas.source = src

	settings := defaultCopyOptions()
	for _, modify := range opts {
		modify(&settings)
	}
	cas.destination = settings.destination
	cas.markerExt = settings.markerExt

	for _, src := range settings.sourcePaths {
		fmt.Println("---")
		originFile := origin.New(src)

		cas.sourcePaths = append(cas.sourcePaths, originFile)
	}
	return &cas
}

func (cas *copyActionState) Start() error {
	fmt.Println("me COPY start!!")
	fmt.Println("grab order:")
	for i, src := range cas.sourcePaths {
		fmt.Println(i+1, src)

	}
	// if err := cas.validatePaths(); err != nil {
	// 	return err
	// }

	fmt.Println("me COPY done!!")
	return nil
}

/*
create
*/
