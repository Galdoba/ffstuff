package inputinfo

import (
	"fmt"
	"strings"
)

type inputdata struct {
	data []string
}

func (i inputdata) String() string {
	str := ""
	for _, line := range i.data {
		str += line + "\n"
	}
	return str
}

type MediaFile struct {
	filename     string
	duration     float64
	start        float64
	bitrate      int //kb/s
	videoStreams []videostream
	audioStreams []videostream
}

func CollectData(input inputdata) *MediaFile {
	mf := MediaFile{}

	return &mf
}

func (inp *inputdata) parseName() string {
	name := ""
	for _, line := range inp.data {
		if !(strings.Contains(line, "ffmpeg ") && (strings.Contains(line, " -i "))) && !strings.Contains(line, " from ") {
			continue
		}
		flds := strings.Fields(line)
		for _, f := range flds {
			fmt.Println("-----")
			fmt.Println(f)
			fmt.Println("-----")
			if !strings.Contains(f, ".") {
				continue
			}
			fs := strings.TrimPrefix(f, "'")
			fs = strings.TrimSuffix(f, "':")
			name = fs
		}
	}
	if name == "" {
		panic(inp.String())
	}
	return name
}

/*
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from 'SeaHorses_FTR_1080p_RU-XX_20_24fps_PR422HQ.mov'

> ffmpeg -i S_LUBOVYU_I_YAROSTYU_235_24FPS_RUS51.mov
  INPUT 0: S_LUBOVYU_I_YAROSTYU_235_24FPS_RUS51.mov
  Duration: 01:56:46.00, start: 0.000000, bitrate: 173412 kb/s
    0:0 (eng) Video: prores (HQ) (apch / 0x68637061), yuv422p10le(tv, bt709/unknown/unknown, progressive), 1920x1080, 168410 kb/s, SAR 1:1 DAR 16:9, 24 fps, 24 tbr, 24 tbn, 24 tbc (default)
    0:1 (eng) Data: none (tmcd / 0x64636D74), 0 kb/s (default)
      handler_name    : Apple Time Code Media Handler
    0:2 (eng) Audio: pcm_s16le (sowt / 0x74776F73), 48000 Hz, 5.1(side), s16, 4608 kb/s (default)
    S_lyubovyu_i_yarostyu
$ ffmpeg -i Spletnica_s02e07__PRT230106115407_SER_00087_18.mp4
ffmpeg version 4.4.1-full_build-www.gyan.dev Copyright (c) 2000-2021 the FFmpeg developers
  built with gcc 11.2.0 (Rev1, Built by MSYS2 project)
  configuration: --enable-gpl --enable-version3 --enable-static --disable-w32threads --disable-autodetect --enable-fontconfig --enable-iconv --enable-gnutls --enable-libxml2 --enable-gmp --enable-lzma --enable-libsnappy --enable-zlib --enable-librist --enable-libsrt --enable-libssh --enable-libzmq --enable-avisynth --enable-libbluray --enable-libcaca --enable-sdl2 --enable-libdav1d --enable-libzvbi --enable-librav1e --enable-libsvtav1 --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxvid --enable-libaom --enable-libopenjpeg --enable-libvpx --enable-libass --enable-frei0r --enable-libfreetype --enable-libfribidi --enable-libvidstab --enable-libvmaf --enable-libzimg --enable-amf --enable-cuda-llvm --enable-cuvid --enable-ffnvcodec --enable-nvdec --enable-nvenc --enable-d3d11va --enable-dxva2 --enable-libmfx --enable-libglslang --enable-vulkan --enable-opencl --enable-libcdio --enable-libgme --enable-libmodplug --enable-libopenmpt --enable-libopencore-amrwb --enable-libmp3lame --enable-libshine --enable-libtheora --enable-libtwolame --enable-libvo-amrwbenc --enable-libilbc --enable-libgsm --enable-libopencore-amrnb --enable-libopus --enable-libspeex --enable-libvorbis --enable-ladspa --enable-libbs2b --enable-libflite --enable-libmysofa --enable-librubberband --enable-libsoxr --enable-chromaprint
  libavutil      56. 70.100 / 56. 70.100
  libavcodec     58.134.100 / 58.134.100
  libavformat    58. 76.100 / 58. 76.100
  libavdevice    58. 13.100 / 58. 13.100
  libavfilter     7.110.100 /  7.110.100
  libswscale      5.  9.100 /  5.  9.100
  libswresample   3.  9.100 /  3.  9.100
  libpostproc    55.  9.100 / 55.  9.100
Input #0, mov,mp4,m4a,3gp,3g2,mj2, from 'Spletnica_s02e07__PRT230106115407_SER_00087_18.mp4':
  Metadata:
    major_brand     : isom
    minor_version   : 512
    compatible_brands: isomiso2avc1mp41
    encoder         : Lavf58.45.100
  Duration: 00:54:07.02, start: 0.000000, bitrate: 8376 kb/s
  Stream #0:0(und): Video: h264 (High 4:2:2) (avc1 / 0x31637661), yuv422p, 1920x1080 [SAR 1:1 DAR 16:9], 7897 kb/s, 25 fps, 25 tbr, 12800 tbn, 50 tbc (default)
    Metadata:
      handler_name    : VideoHandler
      vendor_id       : [0][0][0][0]
      timecode        : 01:00:00:00
  Stream #0:1(rus): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, 5.1, fltp, 341 kb/s (default)
    Metadata:
      handler_name    : SoundHandler
      vendor_id       : [0][0][0][0]
  Stream #0:2(eng): Audio: aac (LC) (mp4a / 0x6134706D), 48000 Hz, stereo, fltp, 128 kb/s
    Metadata:
      handler_name    : SoundHandler
      vendor_id       : [0][0][0][0]
  Stream #0:3(eng): Data: none (tmcd / 0x64636D74)
    Metadata:
      handler_name    : TimeCodeHandler
      timecode        : 01:00:00:00


*/

type videostream struct {
	pix_fmt       string
	width, height int
	bitrate       int
	sar           string
	dar           string
	fps           string
	tbr           string
	tbn           string
	tbc           string
}

type audiostream struct {
	hertz          int
	bitrate        string
	channels       int
	channel_layout string
}
