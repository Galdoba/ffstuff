package copyprocess

import "fmt"

const (
	Status_Pending = 0
)

type copyActionState struct {
	sourcePaths []string
	destination string
}

func NewCopyAction(opts ...Option) *copyActionState {
	cas := copyActionState{}
	// cas.source = src

	settings := defaultCopyOptions()
	for _, modify := range opts {
		modify(&settings)
	}
	cas.sourcePaths = settings.sourcePaths
	cas.destination = settings.destination

	return &cas
}

func (cas *copyActionState) Start() error {
	fmt.Println("me COPY start!!")
	fmt.Println(cas)
	// if err := cas.validatePaths(); err != nil {
	// 	return err
	// }

	fmt.Println("me COPY done!!")
	return nil
}

/*
create
*/
