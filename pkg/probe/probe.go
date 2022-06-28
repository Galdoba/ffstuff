package probe

import (
	"fmt"

	"github.com/Galdoba/devtools/cli/command"
	"github.com/malashin/ffinfo"
	"gopkg.in/AlecAivazis/survey.v1"
)

type Media struct {
	data            string
	f               *ffinfo.File
	mediaReportType string
}

type Report interface {
	Report() string
}

func NewMedia(path string) (*Media, error) {
	f, e := ffinfo.Probe(path)
	mr := Media{}
	if e != nil {
		return nil, e
	}
	mr.f = f
	mr.data = f.String()
	com, err := command.New(command.CommandLineArguments(fmt.Sprintf("ffprobe -i %v", path)),
		command.Set(command.BUFFER_ON),
		command.Set(command.TERMINAL_ON),
	)
	if err != nil {
		return &mr, err
	}
	com.Run()
	mr.data += "\n" + com.StdOut() + "\n" + com.StdErr()

	return &mr, e
}

type mediaFileReport struct {
	filename string
	name     string
	vData    []videoData
	aData    []audioData
}

type videoData struct {
	fps        string
	dimentions dimentions
}

type audioData struct {
	chanLayout string
	chanNum    int
	sampleRate int
	language   string
}

type dimentions struct {
	width  int
	height int
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

func (mr Media) MediaFileReport() *mediaFileReport {
	allStreams := mr.f.Format.NbStreams
	inRep := mediaFileReport{}
	inRep.filename = mr.f.Format.Filename
	for st := 0; st < allStreams; st++ {
		stream := mr.f.Streams[st]
		switch stream.CodecType {
		default:
			fmt.Println("DEBUG: unimplemented or inknown stream type: stream", st)
		case "video":
			vid := videoData{}
			switch stream.RFrameRate {
			default:
				vid.fps = stream.RFrameRate + " (unknown)"
			case "2997/125", "24000/1001":
				vid.fps = "23.98 fps"
			case "24/1":
				vid.fps = "24 fps"
			case "25/1":
				vid.fps = "25 fps"
			}
			vid.dimentions = dimentions{stream.Width, stream.Height}

			inRep.vData = append(inRep.vData, vid)
		case "audio":
			aud := audioData{}
			aud.chanNum = stream.Channels
			aud.sampleRate = stream.Channels
			aud.chanLayout = stream.ChannelLayout
			aud.language = stream.Tags.Language
			inRep.aData = append(inRep.aData, aud)
		}

	}
	fmt.Println(mr.f.String())
	fmt.Println(mr.f.Format.Filename)
	fmt.Println("------------")
	fmt.Println(inRep)

	return &inRep
}

func (inR mediaFileReport) String() string {
	str := fmt.Sprintf("File: %v\n", inR.filename)
	for i := 0; i < len(inR.vData); i++ {
		if i == 0 {
			str += fmt.Sprintf("Video:\n")
		}
		str += fmt.Sprintf(" Stream %v: %v, %v", i, inR.vData[i].dimentions.String(), inR.vData[i].fps)
		str += fmt.Sprintf("\n")
	}
	for i := 0; i < len(inR.aData); i++ {
		if i == 0 {
			str += fmt.Sprintf("Audio:\n")
		}
		str += fmt.Sprintf(" Stream %v: %v", i, inR.aData[i].String())

	}
	return str
}

func SelectAudio(mr *mediaFileReport) []string {
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

func (ad *audioData) String() string {
	return fmt.Sprintf("audio: %v, %v channels (%v)", ad.chanLayout, ad.chanNum, ad.language)
}
