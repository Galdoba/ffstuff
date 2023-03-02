package main

import (
	"fmt"
	"image"
	"image/color"

	vidio "github.com/Galdoba/ffstuff/imported/Vidio"
)

/*
import Image
from subprocess import Popen, PIPE

fps, duration = 24, 100
p = Popen(['ffmpeg', '-y', '-f', 'image2pipe', '-vcodec', 'mjpeg', '-r', '24', '-i', '-', '-vcodec', 'mpeg4', '-qscale', '5', '-r', '24', 'video.avi'], stdin=PIPE)
for i in range(fps * duration):
    im = Image.new("RGB", (300, 300), (i, 1, 1))
    im.save(p.stdin, 'JPEG')
p.stdin.close()
p.wait()
*/

func main() {
	video, _ := vidio.NewVideo(`d:\\tests\\duplicates\\Eve_cut.mov`)
	options := vidio.Options{
		FPS:     video.FPS(),
		Bitrate: video.Bitrate(),
	}
	if video.HasStreams() {
		options.StreamFile = video.FileName()
	}
	video.Width()
	// copyFound := 0
	// copyCounter := 0
	// frameCounter := 0
	// frameMap := make(map[int]int)
	//var lastFrame []byte
	img := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
	//lastImg := image.NewRGBA(image.Rect(0, 0, video.Width(), video.Height()))
	video.SetFrameBuffer(img.Pix)
	//var colorSet1 []color.Color
	//var colorSet2 []color.Color
	writer, err := vidio.NewVideoWriter("output.mp4", video.Width(), video.Height(), &options)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer writer.Close()
	//writer.Write(video.FrameBuffer())
	for video.Read() {
		// frameCounter++
		// //		frameBts := video.FrameBuffer()
		// colorSet1 = memorizeColors(img, video.Width()/4, video.Height()/4)
		// switch isCopy(colorSet1, colorSet2) {
		// case true:
		// 	copyCounter++
		// 	copyFound++
		// 	fmt.Println(frameCounter, copyCounter)
		// case false:
		// 	copyCounter = 0

		// }
		// frameMap[frameCounter] = copyCounter

		// //fmt.Printf("unique frames: %v/%v    \r", frameCounter-copyFound, frameCounter)
		// colorSet2 = colorSet1

		err := writer.Write(video.FrameBuffer())
		if err != nil {
			panic(err.Error())
		}

	}
	fmt.Println("done")
}

func memorizeColors(img *image.RGBA, width, height int) []color.Color {
	clr := []color.Color{}
	for w := 0; w < width; w++ {
		for h := 0; h < height; h++ {
			clr = append(clr, img.At(w, h))
		}
	}
	return clr
}

func isCopy(bts1, bts2 []color.Color) bool {
	if len(bts1) != len(bts2) {
		//fmt.Println("not copy:", len(bts1), "!=", len(bts2))
		return false
	}
	for i := 0; i < len(bts1); i++ {
		if bts2[i] != bts1[i] {
			//fmt.Println("not copy:", bts2[i], "!=", bts1[i], i)
			return false
		}
	}
	//fmt.Printf("full copy       \n")
	return true
}

/*
ffmpeg -hwaccel cuda -c:v h264_cuvid -i d:\tests\duplicates\Eve_cut.mov -pix_fmt bgr24 -f rawvideo -
ffmpeg -f image2pipe -c:v png -r 30000/1001 -i -

ffmpeg -i d:\\tests\\duplicates\\Eve_cut.mov -vf "select='if(gt(scene,0.01),st(1,t),lte(t-ld(1),1))',setpts=N/FRAME_RATE/TB" trimmed.mp4

ffprobe -show_frames d:\\tests\\duplicates\\Eve_cut.mov

ffprobe -show_frames -select_streams v -print_format -unit json=c=1 0001.wmv

 ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 d:\\tests\\duplicates\\Eve_cut.mov

"ffmpeg",
		"-i", video.filename,
		"-f", "image2pipe",
		"-loglevel", "quiet",
		"-pix_fmt", "rgba",
		"-vcodec", "rawvideo",
		"-map", fmt.Sprintf("0:v:%d", video.stream),
		"-",

ffmpeg -i _HD_51.mp4 -vf select='eq(n\,100)+eq(n\,184)+eq(n\,213)' -vsync 0 frames%d.jpg
ffmpeg -t 225 -i _HD_51.mp4 -vf select='not(eq(n\,100))+not(eq(n\,184))+not(eq(n\,213))' -vsync 0 -vcodec rawvideo out.mp4

pipe = subprocess.Popen([FFMPEG_BIN, "-i", src,
                    "-loglevel", "quiet",
                    "-vf", "select=not(mod(n\,100))",
                    "-vsync", "0",
                    "-an",
                    "-f", "image2pipe",
                    "-pix_fmt", "bgr24",
                    "-vcodec", "rawvideo", "-"],
                    stdin=subprocess.PIPE,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.DEVNULL)

ffmpeg -i _HD_51.mp4 -loglevel quiet -vf select=not(mod(n\,100)) -vsync 0 -an -f image2pipe -pix_fmt yuv422 -vcodec rawvideo - | ffmpeg -f rawvideo -i - -vcodec h264 test.mp4


*/

func copyVid() {
	video, _ := vidio.NewVideo(`d:\\tests\\duplicates\\Eve_cut.mov`)
	options := vidio.Options{
		FPS:     video.FPS(),
		Bitrate: video.Bitrate(),
	}
	if video.HasStreams() {
		options.StreamFile = video.FileName()
	}

	writer, _ := vidio.NewVideoWriter("output.mp4", video.Width(), video.Height(), &options)

	writer.Close()
	fmt.Println(writer)
	i := 0
	for video.Read() {
		err := writer.Write(video.FrameBuffer())
		fmt.Println(i, writer)
		if err != nil {
			fmt.Println(err.Error())
		}

		i++
		if i > 50 {
			break
		}
	}
}

//ffmpeg -t 10 -i Eve_cut.mov -an -vcodec copy -f matroska - | ffmpeg -i - -vcodec copy  10SecC2.mov

/*
readder
cmd := exec.Command(
		"ffmpeg",
		"-i", video.filename,
		"-f", "image2pipe",
		"-loglevel", "quiet",
		"-pix_fmt", "rgba",
		"-vcodec", "rawvideo",
		"-map", fmt.Sprintf("0:v:%d", video.stream),
		"-",
	)
ffmpeg -t 1 -i _HD_51.mp4 -f image2pipe -loglevel quiet -pix_fmt rgba -vcodec rawvideo -map 0:v:0 - | ffmpeg -y -loglevel quiet -f rawvideo -vcodec rawvideo -s 1920x1080 -pix_fmt rgba -r 25 -i - -vcodec libx264 -pix_fmt yuv420p outTest.mp4






writter
command := []string{
		"-y", // overwrite output file if it exists.
		"-loglevel", "quiet",
		"-f", "rawvideo",
		"-vcodec", "rawvideo",
		"-s", fmt.Sprintf("%dx%d", writer.width, writer.height), // frame w x h.
		"-pix_fmt", "rgba",
		"-r", fmt.Sprintf("%.02f", writer.fps), // frames per second.
		"-i", "-", // The input comes from stdin.
	}

*/
