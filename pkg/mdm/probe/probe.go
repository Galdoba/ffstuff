package probe

import (
	"fmt"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/malashin/ffinfo"
	"gopkg.in/AlecAivazis/survey.v1"
)

const (
	MediaTypeTrailer4K = "Trailer 4K"
	MediaTypeTrailerHD = "Trailer HD"
	MediaTypeTrailerSD = "Trailer SD"
	MediaTypeFilm4K    = "Film 4K"
	MediaTypeFilmHD    = "Film HD"
	MediaTypeFilmSD    = "Film SD"
	MediaPureSound     = "Sound"
)

type mediaFileReport struct {
	filename  string
	name      string
	f         *ffinfo.File
	data      string
	mediaType string
	vData     []VideoData
	aData     []AudioData
}

type issue struct {
	level string
	text  string
}

type VideoData struct {
	fps        string
	dimentions dimentions
	sar        string
	dar        string
	issues     []string
}

func (vd *VideoData) String() string {
	str := vd.dimentions.String()
	if vd.dar+vd.sar != "" {
		str += " ["
		if vd.sar != "" {
			str += "SAR " + vd.sar
		}
		if vd.dar != "" && vd.sar != "" {
			str += " "
		}
		if vd.dar != "" {
			str += "DAR " + vd.dar
		}
		str += "]"
	}
	str += " " + vd.fps
	return str
}

type dimentions struct {
	width  int
	height int
}

type AudioData struct {
	chanLayout string
	chanNum    int
	sampleRate int
	language   string
	fcMapValue string
}

func (ad *AudioData) ChanLayout() string {
	return ad.chanLayout
}

func (d *dimentions) String() string {
	return fmt.Sprintf("%vx%v", d.width, d.height)
}

/*
для входящего трейлера важно:
колво видеостримов+
расширение+
фпс+
колво аудиостримов+
тип звука

*/

func MediaFileReport(path, mediaType string) (*mediaFileReport, error) {
	report := mediaFileReport{}
	f, e := ffinfo.Probe(path)
	if e != nil {
		return &report, e
	}

	report.filename = f.Format.Filename
	report.mediaType = mediaType
	report.data = f.String()
	com, err := command.New(command.CommandLineArguments(fmt.Sprintf("ffprobe -i %v", path)),
		command.Set(command.BUFFER_ON),
		command.Set(command.TERMINAL_OFF),
	)
	if err != nil {
		return &report, err
	}
	com.Run()
	report.data += "\n" + com.StdOut() + "\n" + com.StdErr()
	allStreams := f.Format.NbStreams

	for st := 0; st < allStreams; st++ {
		stream := f.Streams[st]
		switch stream.CodecType {
		default:
			fmt.Println("DEBUG: unimplemented or unknown stream type: stream", st)
		case "video":
			vid := VideoData{}

			switch stream.RFrameRate {
			default:
				vid.fps = stream.RFrameRate + " (unknown)"
			case "2997/125", "24000/1001", "24/1", "25/1":
				vid.fps = stream.RFrameRate
			}
			vid.dimentions = dimentions{stream.Width, stream.Height}
			//vid.issues = dimentionIssue(vid.dimentions, targetDimentions(mr.mediaType))
			vid.sar = stream.SampleAspectRatio
			vid.dar = stream.DisplayAspectRatio
			report.vData = append(report.vData, vid)
		case "audio":
			aud := AudioData{}
			aud.chanNum = stream.Channels
			aud.sampleRate = stream.Channels
			aud.chanLayout = stream.ChannelLayout
			aud.language = stream.Tags.Language
			report.aData = append(report.aData, aud)
		}

	}
	//fmt.Println(f.String())
	//fmt.Println(f.Format.Filename)
	fmt.Println("------------")
	fmt.Println(report)

	return &report, nil
}

func targetDimentions(mType string) dimentions {
	switch mType {
	default:
		return dimentions{1, 1}
	case MediaTypeFilmHD, MediaTypeTrailerHD:
		return dimentions{1920, 1080}
	case MediaTypeFilm4K, MediaTypeTrailer4K:
		return dimentions{1920, 1080}

	}
}

func (inR mediaFileReport) String() string {
	str := fmt.Sprintf("File: %v\n", inR.filename)
	for i := 0; i < len(inR.vData); i++ {
		if i == 0 {
			str += fmt.Sprintf("Video:\n")
		}
		str += fmt.Sprintf(" Stream %v: %v", i, inR.vData[i].String())
		str += fmt.Sprintf("\n")
	}
	for i := 0; i < len(inR.aData); i++ {
		if i == 0 {
			str += fmt.Sprintf("Audio:\n")
		}
		str += fmt.Sprintf(" Stream %v: %v", i, inR.aData[i].String())
		str += fmt.Sprintf("\n")

	}
	issues := []string{}
	for _, vid := range inR.vData {
		for _, is := range vid.issues {
			if is != "" {
				issues = append(issues, is)
			}
		}
	}
	if len(issues) > 0 {
		str += fmt.Sprintln("ISSUES:")
		for _, iss := range issues {
			str += fmt.Sprintf("%v\n", iss)
		}
	}
	return str
}

func SelectAudio0(mr *mediaFileReport) []string {
	ans := []string{}
	opt := []string{}
	for i, as := range mr.aData {
		opt = append(opt, fmt.Sprintf("%v - %v", i, as.String()))
	}
	msq := survey.MultiSelect{
		Message: "Select tracks",
		Options: opt,
	}

	valid := survey.ComposeValidators()
	fmt.Println(survey.AskOne(&msq, &ans, valid))

	fmt.Println(ans)
	return ans
}

func SelectAudio(allStreams []AudioData) []AudioData {
	ans := []string{}
	opt := []string{}
	for i, as := range allStreams {
		opt = append(opt, fmt.Sprintf("%v - %v", i, as.String()))
	}
	msq := survey.MultiSelect{
		Message: "Select tracks",
		Options: opt,
	}

	valid := survey.ComposeValidators()
	survey.AskOne(&msq, &ans, valid)
	picked := []AudioData{}
	for _, st := range allStreams {
		for _, an := range ans {
			if strings.Contains(an, st.String()) {
				picked = append(picked, st)
			}
		}
	}
	fmt.Println(ans)
	return picked
}

func (ad *AudioData) String() string {
	str := fmt.Sprintf("audio: %v, %v channels", ad.chanLayout, ad.chanNum)
	if ad.language != "" {
		str += " (" + ad.language + ")"
	}
	return str
}

func (mr *mediaFileReport) Issues() []string {
	targetDimentions := dimentions{}
	str := []string{}
	switch mr.mediaType {
	case MediaTypeFilmHD:
		targetDimentions.width, targetDimentions.height = 1920, 1080
		for i, video := range mr.vData {
			if err := dimentionIssue(video.dimentions, targetDimentions); err != nil {
				str = append(str, fmt.Sprintf("Video %v: %v", i, err.Error()))
			}

		}

	}
	return str
}

func dimentionIssue(actual, target dimentions) error {
	if actual.width == target.width && actual.height == target.height {
		return nil
	}
	if actual.width < target.width && actual.height < target.height {
		return fmt.Errorf("Dimention to small for target")
	}
	if actual.width >= target.width && actual.height <= target.height {
		return fmt.Errorf("Need Downscale")
	}
	if actual.width <= target.width && actual.height >= target.height {
		return fmt.Errorf("Need Downscale")
	}
	if actual.width > target.width && actual.height > target.height {
		return fmt.Errorf("Need Downscale")
	}
	return fmt.Errorf("Issue unknown")
}

func fpsIssue(actual string) error {
	if actual != "25 fps" {
		return fmt.Errorf("%v", actual)
	}
	return nil
}

//GETTERS
func (mr *mediaFileReport) FPS() string {
	return mr.vData[0].fps
}

func (mr *mediaFileReport) Audio() []AudioData {
	return mr.aData
}

func (mr *mediaFileReport) Video() []VideoData {
	return mr.vData
}

func (vid *VideoData) SAR() string {
	return vid.sar
}
func (vid *VideoData) Dimentions() (int, int) {
	return vid.dimentions.width, vid.dimentions.height
}

func InterlaceByIdet(path string) (bool, error) {
	com, err := command.New(
		command.CommandLineArguments("ffmpeg", "-filter:v idet -frames:v 37500 -an -f rawvideo -y NUL -i "+path),
		command.Set(command.BUFFER_ON),
		//command.Set(command.TERMINAL_ON),
	)
	fmt.Printf("\nSearching interlace: %v\n", path)
	if err != nil {
		return false, err
	}
	if err := com.Run(); err != nil {
		return false, err
	}
	buf := com.StdOut() + "\n" + com.StdErr()
	for _, line := range strings.Split(buf, "\n") {
		if strings.Contains(line, "[Parsed_idet_0") {
			fmt.Println(line)
		}
	}

	return false, nil
}
