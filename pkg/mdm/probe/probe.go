package probe

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/Galdoba/ffstuff/pkg/info"
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
	DimentionSD        = 10
	DimentionHD        = 11
	Dimention4K        = 12
)

type FileReport struct {
	filename  string
	name      string
	f         *ffinfo.File
	data      string
	mediaType string
	vData     []VideoData
	aData     []AudioData
	issues    []issue
}

type issue struct {
	level string
	text  string
}

var dimentionsSD Dimentions
var dimentionsHD Dimentions
var dimentions4K Dimentions

func init() {
	dimentionsSD = Dimentions{720, 576}
	dimentionsHD = Dimentions{1920, 1080}
	dimentions4K = Dimentions{3840, 2160}
}

func Dimention(d int) Dimentions {
	switch d {
	default:
		return Dimentions{}
	case DimentionSD:
		return dimentionsSD
	case DimentionHD:
		return dimentionsHD
	case Dimention4K:
		return dimentions4K
	}
}

/*
для входящего трейлера важно:
колво видеостримов+
расширение+
фпс+
колво аудиостримов+
тип звука

*/

func NewReport(path string) (*FileReport, error) {
	report := FileReport{}
	f, e := ffinfo.Probe(path)
	if e != nil {
		return &report, e
	}

	report.filename = f.Format.Filename
	report.data = f.String()
	// com, err := command.New(command.CommandLineArguments(fmt.Sprintf("ffprobe -i %v", path)),
	// 	command.Set(command.BUFFER_ON),
	// 	command.Set(command.TERMINAL_OFF),
	// )
	// if err != nil {
	// 	return &report, err
	// }
	// com.Run()
	// report.data += "\n" + com.StdOut() + "\n" + com.StdErr()
	allStreams := f.Format.NbStreams
	videoStream := 0
	audStream := 0
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
			case "2997/125", "24000/1001", "24/1", "25/1", "27021/1127":
				vid.fps = stream.RFrameRate
			}
			vid.dimentions = Dimentions{stream.Width, stream.Height}
			//vid.issues = dimentionIssue(vid.Dimentions, targetDimentions(mr.mediaType))
			vid.sar = stream.SampleAspectRatio
			vid.dar = stream.DisplayAspectRatio
			vid.fcMapKey = fmt.Sprintf(":v:%v", videoStream)
			report.vData = append(report.vData, vid)
			videoStream++
		case "audio":
			aud := AudioData{}
			aud.chanNum = stream.Channels
			aud.sampleRate = stream.Channels
			aud.chanLayout = stream.ChannelLayout
			aud.language = stream.Tags.Language
			aud.fcMapKey = fmt.Sprintf(":a:%v", audStream)
			report.aData = append(report.aData, aud)
			audStream++
		}

	}
	//fmt.Println(f.String())
	//fmt.Println(f.Format.Filename)
	//fmt.Println("------------")
	//fmt.Println(report)

	return &report, nil
}

func targetDimentions(mType string) Dimentions {
	switch mType {
	default:
		return Dimentions{1, 1}
	case MediaTypeFilmHD, MediaTypeTrailerHD:
		return Dimentions{1920, 1080}
	case MediaTypeFilm4K, MediaTypeTrailer4K:
		return Dimentions{3840, 2160}

	}
}

func (inR FileReport) String() string {
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

func SelectAudio0(mr *FileReport) []string {
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
		Message: "Аудио стримы: ",
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

func (mr *FileReport) Issues() []string {
	targetDimentions := Dimentions{}
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

func dimentionIssue(actual, target Dimentions) error {
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
func (mr *FileReport) FPS() string {
	return mr.vData[0].fps
}

func (mr *FileReport) Audio() []AudioData {
	return mr.aData
}

func (mr *FileReport) Video() []VideoData {
	return mr.vData
}

func InterlaceByIdet(path string) (bool, error) {
	length, err := info.Duration(path)
	if err != nil {
		return false, err
	}
	leng := (length / 2)
	if leng > 180 {
		leng = 180
	}
	com, err := command.New(
		command.CommandLineArguments("ffmpeg", "-filter:v idet -frames:v 999 -an -f rawvideo -y NUL -i "+path), //-ss "+fmt.Sprintf("%v", leng)+"
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
	dataRAW := ""
	buf := com.StdOut() + "\n" + com.StdErr()
	for _, line := range strings.Split(buf, "\n") {
		if strings.Contains(line, "[Parsed_idet_0") {
			fmt.Println(line)
			if strings.Contains(line, "Multi frame detection") {

			}
			dataRAW = line
		}
	}
	data := strings.Fields(dataRAW)
	frames := make(map[string]int)
	for o, d := range data {
		switch o {
		case 7, 9, 11, 13:
			n, err := strconv.Atoi(d)
			if err != nil {
				return false, err
			}
			frames[data[o-1]] = n
		}
	}
	if frames["TFF:"]+frames["BFF:"] >= 40 {
		return true, nil
	}

	return false, nil
}

func StreamsData(paths ...string) ([]Stream, error) {
	streams := []Stream{}
	for i, path := range paths {
		f, e := ffinfo.Probe(path)
		if e != nil {
			return nil, e
		}
		//fmt.Println(f.String())
		allStreams := f.Format.NbStreams
		videoStream := 0
		audStream := 0
		for st := 0; st < allStreams; st++ {
			stream := f.Streams[st]
			switch stream.CodecType {
			default:
				//fmt.Println("DEBUG: unimplemented or unknown stream type: stream", st)
				//TODO: обработка стримов Data (d), Subtitles (s), Attachments (t)
			case "video":
				vid := fillVideoData(stream)
				vid.SetFcMapKey(i, videoStream)
				streams = append(streams, &vid)
				videoStream++
			case "audio":
				aud := fillAudioData(stream)
				aud.SetMapKey(i, audStream)
				streams = append(streams, &aud)
				audStream++
			}
		}
	}
	return streams, nil
}

func NewReport0(path string) (*FileReport, error) {
	report := FileReport{}
	f, e := ffinfo.Probe(path)
	if e != nil {
		return &report, e
	}

	report.filename = f.Format.Filename
	report.data = f.String()
	// com, err := command.New(command.CommandLineArguments(fmt.Sprintf("ffprobe -i %v", path)),
	// 	command.Set(command.BUFFER_ON),
	// 	command.Set(command.TERMINAL_OFF),
	// )
	// if err != nil {
	// 	return &report, err
	// }
	// com.Run()
	// report.data += "\n" + com.StdOut() + "\n" + com.StdErr()
	allStreams := f.Format.NbStreams
	videoStream := 0
	audStream := 0
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
			case "2997/125", "24000/1001", "24/1", "25/1", "27021/1127":
				vid.fps = stream.RFrameRate
			}
			vid.dimentions = Dimentions{stream.Width, stream.Height}
			//vid.issues = dimentionIssue(vid.Dimentions, targetDimentions(mr.mediaType))
			vid.sar = stream.SampleAspectRatio
			vid.dar = stream.DisplayAspectRatio
			vid.fcMapKey = fmt.Sprintf(":v:%v", videoStream)
			report.vData = append(report.vData, vid)
			videoStream++
		case "audio":
			aud := AudioData{}
			aud.chanNum = stream.Channels
			aud.sampleRate = stream.Channels
			aud.chanLayout = stream.ChannelLayout
			aud.language = stream.Tags.Language
			aud.fcMapKey = fmt.Sprintf(":a:%v", audStream)
			report.aData = append(report.aData, aud)
			audStream++
		}

	}
	//fmt.Println(f.String())
	//fmt.Println(f.Format.Filename)
	//fmt.Println("------------")
	//fmt.Println(report)

	return &report, nil
}

func SeparateByTypes(strs []Stream) ([]VideoData, []AudioData) {
	video := []VideoData{}
	audio := []AudioData{}
	for _, str := range strs {
		switch t := str.(type) {
		case *VideoData:
			video = append(video, *t)
		case *AudioData:
			audio = append(audio, *t)
		default:
		}
	}
	return video, audio
}
