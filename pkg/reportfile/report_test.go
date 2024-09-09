package reportfile

import (
	"fmt"
	"testing"
)

func TestReport(t *testing.T) {
	r := New(
		NewField("num_files", "3"),
		NewField("processed___video", "path/video.mp4"),
		NewField("processed__audio1", "path/audio1.mp4"),
		NewField("processed__audio2", "path/audio2.mp4"),
		NewField("bitrate original", "12345 kb/s"),
		NewField("bitrate processed", "23456 kb/s"),
		NewField("rename_error", "rename failed, return original name"),
	)
	bt, _ := r.Marshal()
	fmt.Println(string(bt))
}
