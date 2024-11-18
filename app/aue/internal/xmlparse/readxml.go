package xmlparse

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Galdoba/ffstuff/app/aue/internal/xmlparse/filmdata"
	"github.com/Galdoba/ffstuff/app/aue/internal/xmlparse/seriesdata"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func ParseFilmVideoData(path string) (*filmdata.VideoData, error) {
	bt, err := ReadAsUTF8(path)
	if err != nil {
		return nil, err
	}
	return filmdata.Unmarshal(bt)
}

func ParseSeriesVideoData(path string) (*seriesdata.VideoData, error) {
	bt, err := ReadAsUTF8(path)
	if err != nil {
		return nil, err
	}
	return seriesdata.Unmarshal(bt)
}

func ReadAsUTF8(path string) ([]byte, error) {
	return readWindows1251AsUTF8(path)
}

func readWindows1251AsUTF8(filename string) ([]byte, error) {
	// Read UTF-8 from a GBK encoded file.
	bt := []byte{}
	f, err := os.Open(filename)
	if err != nil {
		return bt, fmt.Errorf("failed to open file: %v", err)
	}
	r := transform.NewReader(f, charmap.Windows1251.NewDecoder())
	//Read converted UTF-8 from `r` as needed.
	lines := []string{}
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		text := sc.Text()

		text = strings.ReplaceAll(text, "windows-1251", "UTF-8")
		lines = append(lines, text)
	}
	if err = sc.Err(); err != nil {

		return bt, fmt.Errorf("failed to scan file: %v", err)
	}
	if err = f.Close(); err != nil {
		return bt, fmt.Errorf("failed to close file: %v", err)
	}
	text := strings.Join(lines, "\n")

	return []byte(text), nil
}
