package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Galdoba/ffstuff/app/mfline/ump"
)

type ShellCompiler struct {
	mediaFile string
	srtFile   string
	season    string
	episode   string
	langs     map[string]string
	layout    map[string]string
}

func NewCompiler(files []string) (*ShellCompiler, error) {
	sc := ShellCompiler{}
	sc.langs = make(map[string]string)
	sc.layout = make(map[string]string)
	for _, file := range files {
		parts := strings.Split(file, ".")
		ext := parts[len(parts)-1]
		switch ext {
		case "mp4":
			re := regexp.MustCompile(`(s[0-9]{1,}e[0-9]{1,})`)
			epiTag := re.FindString(file)
			if epiTag == "" {
				return nil, fmt.Errorf("no data on episode")
			}
			mp := ump.NewProfile()
			if err := mp.ConsumeFile(in_dir + file); err != nil {
				return nil, fmt.Errorf("scan: %v", err)
			}
			sc.mediaFile = file
			sc.season, sc.episode = seasonAndEpisode(epiTag)
			sc.layout = streamLayouts(mp)
			sc.langs = streamLangs(mp)
		case "srt":
			sc.srtFile = file
		}
	}
	return &sc, nil
}

func streamLayouts(mp *ump.MediaProfile) map[string]string {
	layouts := make(map[string]string)
	i := 0
	for _, stream := range mp.Streams {
		if stream.Codec_type != "audio" {
			continue
		}
		streamCode := fmt.Sprintf("[0:a:%v]", i)
		layouts[streamCode] = stream.Channel_layout
		i++
	}
	return layouts
}

func streamLangs(mp *ump.MediaProfile) map[string]string {
	langs := make(map[string]string)
	i := 0
	for _, stream := range mp.Streams {
		if stream.Codec_type != "audio" {
			continue
		}
		streamCode := fmt.Sprintf("[0:a:%v]", i)
		langs[streamCode] = stream.Tags["language"]
		i++
	}
	return langs
}

func seasonAndEpisode(tag string) (string, string) {
	parts := strings.Split(tag, "e")
	if len(parts) == 1 {
		return "", ""
	}
	parts[0] = strings.TrimPrefix(parts[0], "s")
	return parts[0], parts[1]
}
