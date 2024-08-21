package bridge

import (
	"fmt"

	source "github.com/Galdoba/ffstuff/app/aue/internal/files/sourcefile"
	"github.com/Galdoba/ffstuff/app/aue/internal/metainfo"
	"github.com/Galdoba/ffstuff/pkg/ump"
)

type fileBridge struct {
	sourceNames    []string
	sourceProfiles map[string]*ump.MediaProfile
	metaInfo       metainfo.Collection
	targetNames    []string
}

func New() *fileBridge {
	br := fileBridge{}
	br.sourceProfiles = make(map[string]*ump.MediaProfile)
	br.metaInfo = metainfo.NewCollection()
	return &br
}

func (br *fileBridge) ConnectSources(sfs ...source.SourceFile) error {
	for _, source := range sfs {
		br.sourceNames = append(br.sourceNames, source.Name())
		br.sourceProfiles[source.Name()] = source.Profile()
		metas := metainfo.Parse(source.Name())
		if err := br.metaInfo.Add(metas...); err != nil {
			return fmt.Errorf("source connection failed: %v", err)
		}

	}
	//check repeated connection
	return nil
}
