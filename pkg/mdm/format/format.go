package format

import "fmt"

const (
	Trailer4K   = "Trailer 4K"
	TrailerHD   = "Trailer HD"
	TrailerSD   = "Trailer SD"
	Film4K      = "Film 4K"
	FilmHD      = "Film HD"
	FilmSD      = "Film SD"
	PureSound   = "Sound"
	DimentionSD = 10
	DimentionHD = 11
	Dimention4K = 12
)

var dimentionsSD VideoDimention
var dimentionsHD VideoDimention
var dimentions4K VideoDimention

func init() {
	dimentionsSD = VideoDimention{720, 576, "SD", ""}
	dimentionsHD = VideoDimention{1920, 1080, "HD", ""}
	dimentions4K = VideoDimention{3840, 2160, "4K", ""}
}

type VideoDimention struct {
	width       int
	height      int
	maximumSize string
	issue       string
}

func (vd *VideoDimention) MaxSize() string {
	return vd.maximumSize
}

func (vd *VideoDimention) Issue() string {
	return vd.issue
}

func NewDimention(w, h int) *VideoDimention {
	vd := VideoDimention{}
	vd.width = w
	vd.height = h
	switch {
	case w == 720 && h == 576:
		return &dimentionsSD
	case w == 1920 && h == 1080:
		return &dimentionsSD
	case w == 3840 && h == 2160:
		return &dimentionsSD
	}
	if w < 720 {
		vd.issue = "Less than SD"
	}
	if w <= 1919 {
		vd.maximumSize = "SD"
	}
	if w <= 3839 {
		vd.maximumSize = "HD"
	}
	if w >= 3840 {
		vd.maximumSize = "4K"
	}
	return &vd
}

func DimentionPreset(dp int) VideoDimention {
	switch dp {
	default:
		return VideoDimention{0, 0, "", "dimention unknown"}
	case DimentionSD:
		return dimentionsSD
	case DimentionHD:
		return dimentionsHD
	case Dimention4K:
		return dimentions4K
	}
}

type TargetFormat struct {
	Dimention   VideoDimention
	SAR         string
	VidCodec    string
	CRF         string
	Pix_fmt     string
	G           string
	PresetVideo string
	AudCodec    string
	Compression string
	MapMetadata string
	MapChapters string
}

/*
# 4K(Ultra HD 4K, 3840 × 2160, 16:9):
clear && \
mv  ~/IN/FILE ~/IN/_IN_PROGRESS/ && \
fflite -r 25 -i ~/IN/_IN_PROGRESS/FILE \
-filter_complex "[0:a:0]aresample=48000,atempo=25/24[audio];
[0:v:0]setsar=1/1[video_4k];
[0:v:0]scale=1920:1080,setsar=1/1[video_hd];
[0:v:0]scale=720:576,setsar=64/45,unsharp=3:3:0.3:3:3:0[video_sd]" \
-map [audio]    @alac0 NAME_AUDIORUS20.m4a \
-map [video_4k] @crf18 NAME_4K.mp4 \
-map [video_hd] @crf16 NAME_HD.mp4 \
-map [video_sd] @crf13 NAME_SD.mp4 \
&& touch  NAME.ready \
&& at now + 10 hours <<< "mv ~/IN/_IN_PROGRESS/FILE OUTPATH"
*/

func SetAs(format string) (*TargetFormat, error) {
	tf := TargetFormat{}
	switch format {
	default:
		return nil, fmt.Errorf("unknown format '%v'", format)
	case Film4K:
		tf.Dimention = DimentionPreset(Dimention4K)
		tf.SAR = "1/1"
		tf.CRF = "18"
	case FilmHD:
		tf.Dimention = DimentionPreset(DimentionHD)
		tf.SAR = "1/1"
		tf.CRF = "16"
	case FilmSD:
		tf.Dimention = DimentionPreset(DimentionSD)
		tf.SAR = "64/45"
		tf.CRF = "13"
	}
	tf.PresetVideo = "medium"
	tf.VidCodec = "libx264"
	tf.Pix_fmt = "yuv420p"
	tf.G = "0"
	tf.AudCodec = "alac"
	tf.Compression = "0"
	tf.MapMetadata = "-1"
	tf.MapChapters = "-1"
	return &tf, nil
}
