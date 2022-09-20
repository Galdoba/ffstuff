package probe

import (
	"fmt"
	"strings"

	"github.com/malashin/ffinfo"
)

//Stream - общий интерфейс стрима для видео, аудио, (ТУДУ) субтитров
type Stream interface {
	PrintStreamData() string
	CodecType() string
}

func (vid *VideoData) PrintStreamData() string {
	return vid.String()
}

func (aud *AudioData) PrintStreamData() string {
	return aud.String()
}

////////////////////////////////////////////////////////////////////////////////

//AudioData - актуальная информация аудиострима для работы с ffmpeg
type AudioData struct {
	chanLayout string
	chanNum    int
	sampleRate int
	language   string
	fcMapKey   string
	position   string
}

func (ad *AudioData) String() string {
	str := ""
	if ad.fcMapKey != "" {
		str += " [" + ad.fcMapKey + "] "
	}
	str += fmt.Sprintf("%v channels", ad.chanNum)
	if ad.chanLayout != "" {
		str += ": " + ad.chanLayout + ""
	}
	return str
}

func (aud *AudioData) CodecType() string {
	return "audio"
}

//getters
func (ad *AudioData) ChanLayout() string {
	return ad.chanLayout
}
func (ad *AudioData) ChanNum() int {
	return ad.chanNum
}
func (ad *AudioData) SampleRate() int {
	return ad.sampleRate
}
func (ad *AudioData) Language() string {
	return ad.language
}
func (ad *AudioData) FCmapKey() string {
	return ad.fcMapKey
}
func (ad *AudioData) Position() string {
	return ad.position
}

//setters
func (ad *AudioData) SetChanLayout(c string) {
	ad.chanLayout = c
}
func (ad *AudioData) SetChanNum(c int) {
	ad.chanNum = c
}
func (ad *AudioData) SetSampleRate(s int) {
	ad.sampleRate = s
}
func (ad *AudioData) SetLanguage(l string) {
	ad.language = l
}
func (ad *AudioData) SetFcMapKey(f string) {
	ad.fcMapKey = f
}
func (ad *AudioData) SetPosition(p string) {
	ad.position = p
}

func fillAudioData(stream ffinfo.Stream) AudioData {
	aud := AudioData{}
	aud.chanNum = stream.Channels
	aud.sampleRate = stream.Channels
	aud.chanLayout = stream.ChannelLayout
	aud.definePosition(stream.ChannelLayout)
	aud.language = stream.Tags.Language
	return aud
}

func (aud *AudioData) definePosition(layout string) {
	switch {
	case containsAnyOf(layout, "(FL)"):
		aud.position = "L"
	case containsAnyOf(layout, "(FR)"):
		aud.position = "R"
	case containsAnyOf(layout, "(LFE)"):
		aud.position = "LFE"
	case containsAnyOf(layout, "(BL)"):
		aud.position = "Ls"
	case containsAnyOf(layout, "(BR)"):
		aud.position = "Rs"
	}
}

func containsAnyOf(str string, el ...string) bool {
	for _, e := range el {
		if strings.Contains(str, strings.ToUpper(e)) {
			return true
		}
	}
	return false
}

func (aud *AudioData) SetMapKey(file, stream int) {
	aud.fcMapKey = fmt.Sprintf("%v:a:%v", file, stream)
}

////////////////////////////////////////////////////////////////////////////////

type VideoData struct {
	fps        string
	dimentions Dimentions
	sar        string
	dar        string
	fcMapKey   string
	issues     []string
}

func (vd *VideoData) String() string {
	str := ""
	if vd.fcMapKey != "" {
		str += " [" + vd.fcMapKey + "]"
	}
	str += " " + vd.dimentions.String()
	fps := ""
	switch vd.fps {
	default:
		fps = vd.fps + " (!)"
	case "2997/125", "24000/1001", "27021/1127":
		fps = "23.976"
	case "25/1":
		fps = "25"
	case "24/1":
		fps = "24"

	}
	str += " FPS:" + fps + " "
	if vd.dar+vd.sar != "" {
		str += "["
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
	return str
}

func (vd *VideoData) CodecType() string {
	return "video"
}

func (vd *VideoData) FPS() string {
	return vd.fps
}

func (vid *VideoData) Dimentions() (int, int) {
	return vid.dimentions.width, vid.dimentions.height
}

func (vid *VideoData) SAR() string {
	return vid.sar
}

func (vid *VideoData) DAR() string {
	return vid.dar
}

func (vid *VideoData) FCmapKey() string {
	return vid.fcMapKey
}

func (vd *VideoData) SetFps(f string) {
	vd.fps = f
}
func (vd *VideoData) SetDimentions(d Dimentions) {
	vd.dimentions = d
}
func (vd *VideoData) SetSar(s string) {
	vd.sar = s
}
func (vd *VideoData) SetDar(d string) {
	vd.dar = d
}
func (vd *VideoData) SetFcMapKey(file, stream int) {
	vd.fcMapKey = fmt.Sprintf("%v:v:%v", file, stream)
}

func fillVideoData(stream ffinfo.Stream) VideoData {
	vid := VideoData{}
	//vid := fillVideoData(stream)
	switch stream.RFrameRate {
	default:
		vid.fps = stream.RFrameRate + " (WARNING)"
	case "2997/125", "24000/1001", "24/1", "25/1", "27021/1127":
		vid.fps = stream.RFrameRate
	}
	vid.dimentions = Dimentions{stream.Width, stream.Height}
	//vid.issues = dimentionIssue(vid.Dimentions, targetDimentions(mr.mediaType))
	vid.sar = stream.SampleAspectRatio
	vid.dar = stream.DisplayAspectRatio
	return vid
}

//////////////

type Dimentions struct {
	width  int
	height int
}

func (d *Dimentions) String() string {
	return fmt.Sprintf("%vx%v", d.width, d.height)
}
