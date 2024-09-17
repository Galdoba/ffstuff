package bridge

import (
	"fmt"
	"strings"

	. "github.com/Galdoba/ffstuff/app/aue/internal/define"
	source "github.com/Galdoba/ffstuff/app/aue/internal/files/sourcefile"
	target "github.com/Galdoba/ffstuff/app/aue/internal/files/targetfile"
	"github.com/Galdoba/ffstuff/app/aue/internal/media"
	"github.com/Galdoba/ffstuff/app/aue/internal/metainfo"
	"github.com/Galdoba/ffstuff/pkg/logman"
)

type fileBridge struct {
	sourceFiles       []*source.SourceFile
	metaInfo          metainfo.Collection
	targetFiles       []*target.TargetFile
	targetsModified   []target.TargetFile
	outputBase        string
	destinationPrefix string
}

func New() *fileBridge {
	br := fileBridge{}
	//br.sourceProfiles = make(map[string]*ump.MediaProfile)
	br.metaInfo = metainfo.NewCollection()
	return &br
}

func (br *fileBridge) Connect(sources []*source.SourceFile, targets []*target.TargetFile) error {
	for _, err := range []error{
		br.connectSources(sources...),
		br.connectTargets(targets...),
		br.sealConnections(),
	} {
		if err != nil {
			return fmt.Errorf("source/target connection failed: %v", err)
		}
	}
	return nil
}

func (br *fileBridge) connectSources(sfs ...*source.SourceFile) error {
	sourceMap := make(map[string]int)
	for _, source := range sfs {
		// if err := source.FillProfile(); err != nil {
		// 	return fmt.Errorf("source connecton: %v")
		// }
		if err := source.Validate(); err != nil {
			return fmt.Errorf("source validation: %v")
		}

		br.sourceFiles = append(br.sourceFiles, source)
		sourceMap[source.Name()]++
	}
	//check repeated connection
	for _, v := range sourceMap {
		if v != 1 {
			return fmt.Errorf("source connections not unique")
		}
	}
	return nil
}

func (br *fileBridge) connectTargets(targets ...*target.TargetFile) error {
	targetMap := make(map[string]int)
	for _, target := range targets {
		br.targetFiles = append(br.targetFiles, target)
		targetMap[target.ClaimedGoal]++
	}
	for _, v := range targetMap {
		if v != 1 {
			return fmt.Errorf("target connections not unique")
		}
	}
	return nil
}

func (br *fileBridge) sealConnections() error {
	for _, source := range br.sourceFiles {
		metas := metainfo.Parse(source.Name())
		metas = append(metas, audioMeta(source)...)
		if err := br.metaInfo.Add(metas...); err != nil {
			return fmt.Errorf("sealing failed: %v", err)
		}

	}
	br.defineProjectName()
	br.calculateEditDir()
	err := br.updateTargets()
	if err != nil {
		return fmt.Errorf("sealing failed: %v", err)
	}

	return nil
}

func (br *fileBridge) defineProjectName() {
	br.outputBase = br.metaInfo.Show(META_Base)
	if br.metaInfo.Show(META_Season) != "" {
		br.outputBase += "_s" + br.metaInfo.Show(META_Season)
	}
	if br.metaInfo.Show(META_Episode) != "" {
		br.outputBase += "_" + br.metaInfo.Show(META_Episode)
	}
	if br.metaInfo.Show(META_PRT) != "" {
		br.outputBase += "_" + br.metaInfo.Show(META_PRT)
	}
}

func (br *fileBridge) updateTargets() error {
	streamTypesByName := source.MapByStreamTypes(br.sourceFiles)
	for i, target := range br.targetFiles {
		target.UsedSourcesNames = filterSources(target, streamTypesByName)
		if target.ExpectedName != "" {
			return fmt.Errorf("target was modified before: %v", target)
		}
		target.ExpectedName += br.ProjectBase()
		switch target.ClaimedGoal {
		case PURPOSE_Output_Video:
			target.ExpectedName += "_HD.mp4"
		case PURPOSE_Output_Audio1:
			target.ExpectedName += "_AUDIO" + br.metaInfo.Show(META_Audio_Lang_0) + br.metaInfo.Show(META_Audio_Layout_0) + ".m4a"
		case PURPOSE_Output_Audio2:
			target.ExpectedName += "_AUDIO" + br.metaInfo.Show(META_Audio_Lang_1) + br.metaInfo.Show(META_Audio_Layout_1) + ".m4a"
		case PURPOSE_Output_Subs:
			target.ExpectedName += ".srt"
		}
		br.targetFiles[i] = target
	}
	return nil
}

func (br *fileBridge) calculateEditDir() error {
	if len(br.targetFiles) == 0 {
		logman.Warn("edit path omited: no target files received")
		return nil
	}
	br.destinationPrefix = br.metaInfo.Show(META_Base) + "_s" + br.metaInfo.Show(META_Season) + "/"
	return nil
}

func filterSources(target *target.TargetFile, streamTypesByName map[string][]int) []string {
	filteredSources := []string{}

	switch target.ClaimedGoal {
	case PURPOSE_Output_Video:
		for name, nums := range streamTypesByName {
			if nums[0] == 1 {
				filteredSources = append(filteredSources, name)
			}
		}
	case PURPOSE_Output_Audio1:
		for name, nums := range streamTypesByName {
			if nums[1] >= 1 {
				filteredSources = append(filteredSources, name)
			}
		}
	case PURPOSE_Output_Audio2:
		for name, nums := range streamTypesByName {
			if nums[1] == 2 {
				filteredSources = append(filteredSources, name)
			}
		}
	case PURPOSE_Output_Subs:
		for name, nums := range streamTypesByName {
			if nums[2] == 1 {
				filteredSources = append(filteredSources, name)
			}
		}
	}
	return filteredSources
}

func audioMeta(source *source.SourceFile) []metainfo.Meta {
	metas := []metainfo.Meta{}
	audioStreams := media.AudioStreams(source.Profile())
	for i, stream := range audioStreams {
		langKey := fmt.Sprintf("%v%v", META_Audio_Lang_, i)
		langVal := strings.ToUpper(stream.Tags[STREAM_TAG_LANGUAGE])
		metas = append(metas, metainfo.NewMeta(langKey, langVal))

		lauoutKey := fmt.Sprintf("%v%v", META_Audio_Layout_, i)
		layoutVal := normalizeLayout(stream.Channel_layout)
		metas = append(metas, metainfo.NewMeta(lauoutKey, layoutVal))
	}
	return metas
}

func normalizeLayout(layout string) string {
	switch layout {
	default:
		if strings.Contains(layout, "5") && strings.Contains(layout, "1") {
			return "51"
		}
		return layout
	case "stereo":
		return "20"
	}
}

func (br *fileBridge) Sources() []*source.SourceFile {
	return br.sourceFiles
}

func (br *fileBridge) Targets() []*target.TargetFile {
	return br.targetFiles
}

func (br *fileBridge) ProjectBase() string {
	return br.outputBase
}

func (br *fileBridge) DestinationPrefix() string {
	return br.destinationPrefix
}
