package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Galdoba/ffstuff/app/mfline/ump"
	"github.com/Galdoba/ffstuff/app/prodamed/config"
)

type ShellCompiler struct {
	nameBase      string
	mediaFile     string
	srtFile       string
	season        string
	episode       string
	langs         map[string]string
	layout        map[string]string
	root          map[string]string
	shellFileName string
	text          string
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
			sc.nameBase = headOfString(file, "--")
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
	sc.shellFileName = fmt.Sprintf("%v_s%ve%v.sh", sc.nameBase, sc.season, sc.episode)
	return &sc, nil
}

func (sc *ShellCompiler) injectLogistics(cfg *config.Config) {
	sc.root = cfg.Option.PATH
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
		langValue := stream.Tags["language"]
		if langValue == "" {
			langValue = "zzz"
		}
		langs[streamCode] = langValue
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
