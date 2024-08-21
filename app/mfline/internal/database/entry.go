package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Galdoba/ffstuff/pkg/ump"
)

type Entry struct {
	File    string            `json:"file"`
	Profile *ump.MediaProfile `json:"Data,omitempty"`
	// Format         *Format   `json:"format"`
	// Streams        []*Stream `json:"streams,omitempty"`
	// ScansCompleted []string  `json:"scans completed,omitempty"`
}

func NewEntry() *Entry {
	return &Entry{}
}

/*
C
R
U
D
*/

func (db *DBjson) Create(file string) error {
	file = sourceKey(file)
	_, err := os.Stat(db.Dir() + file + ".json")
	if err == nil {
		return fmt.Errorf("can't create entry '%v': file exist", file)
	}
	e := &Entry{}
	e.File = file
	err = db.Save(e)
	if err != nil {
		return fmt.Errorf("creation error: %v", err)
	}
	return nil
}

func (db *DBjson) Read(file string) (*Entry, error) {
	key := sourceKey(file)
	_, err := os.Stat(db.Dir() + key + ".json")
	if err != nil {
		return nil, fmt.Errorf("can't read entry '%v': file not exist", db.Dir()+key+".json")
	}
	e, err := db.Load(key)
	if err != nil {
		return nil, fmt.Errorf("load error: %v", err)
	}
	return e, nil
}

func (db *DBjson) Update(e *Entry) error {
	return db.Save(e)
}

func (db *DBjson) Delete(e *Entry) error {
	file := e.File
	path := db.Dir() + file + ".json"
	return os.Remove(path)
}

type Format struct {
	Bit_rate         string            `json:"bit_rate,omitempty"`
	Duration         string            `json:"duration,omitempty"`
	Filename         string            `json:"filename,omitempty"`
	Format_long_name string            `json:"format_long_name,omitempty"`
	Format_name      string            `json:"format_name,omitempty"`
	Nb_programs      int               `json:"nb_programs,omitempty"`
	Nb_streams       int               `json:"nb_streams,omitempty"`
	Probe_score      int               `json:"probe_score,omitempty"`
	Size             string            `json:"size,omitempty"`
	Start_time       string            `json:"start_time,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
}

type Stream struct {
	Avg_frame_rate         string                  `json:"avg_frame_rate,omitempty"`
	Bit_rate               string                  `json:"bit_rate,omitempty"`
	Bits_per_raw_sample    string                  `json:"bits_per_raw_sample,omitempty"`
	Bits_per_sample        int                     `json:"bits_per_sample,omitempty"`
	Channel_layout         string                  `json:"channel_layout,omitempty"`
	Channels               int                     `json:"channels,omitempty"`
	Chroma_location        string                  `json:"chroma_location,omitempty"`
	Closed_captions        int                     `json:"closed_captions,omitempty"`
	Codec_long_name        string                  `json:"codec_long_name,omitempty"`
	Codec_name             string                  `json:"codec_name,omitempty"`
	Codec_tag              string                  `json:"codec_tag,omitempty"`
	Codec_tag_string       string                  `json:"codec_tag_string,omitempty"`
	Codec_time_base        string                  `json:"codec_time_base,omitempty"`
	Codec_type             string                  `json:"codec_type,omitempty"`
	Coded_height           int                     `json:"coded_height,omitempty"`
	Coded_width            int                     `json:"coded_width,omitempty"`
	Color_primaries        string                  `json:"color_primaries,omitempty"`
	Color_range            string                  `json:"color_range,omitempty"`
	Color_space            string                  `json:"color_space,omitempty"`
	Color_transfer         string                  `json:"color_transfer,omitempty"`
	Display_aspect_ratio   string                  `json:"display_aspect_ratio,omitempty"`
	Divx_packed            string                  `json:"divx_packed,omitempty"`
	Dmix_mode              string                  `json:"dmix_mode,omitempty"`
	Duration               string                  `json:"duration,omitempty"`
	Duration_ts            int                     `json:"duration_ts,omitempty"`
	Field_order            string                  `json:"field_order,omitempty"`
	Has_b_frames           int                     `json:"has_b_frames,omitempty"`
	Height                 int                     `json:"height,omitempty"`
	Id                     string                  `json:"id,omitempty"`
	Index                  int                     `json:"index,omitempty"`
	Is_avc                 string                  `json:"is_avc,omitempty"`
	Level                  int                     `json:"level,omitempty"`
	Loro_cmixlev           string                  `json:"loro_cmixlev,omitempty"`
	Loro_surmixlev         string                  `json:"loro_surmixlev,omitempty"`
	Ltrt_cmixlev           string                  `json:"ltrt_cmixlev,omitempty"`
	Ltrt_surmixlev         string                  `json:"ltrt_surmixlev,omitempty"`
	Max_bit_rate           string                  `json:"max_bit_rate,omitempty"`
	Nal_length_size        string                  `json:"nal_length_size,omitempty"`
	Nb_frames              string                  `json:"nb_frames,omitempty"`
	Pix_fmt                string                  `json:"pix_fmt,omitempty"`
	Profile                string                  `json:"profile,omitempty"`
	Progressive_frames_pct float64                 `json:"progressive_frames_pct,omitempty"`
	Quarter_sample         string                  `json:"quarter_sample,omitempty"`
	R_frame_rate           string                  `json:"r_frame_rate,omitempty"`
	Refs                   int                     `json:"refs,omitempty"`
	Sample_aspect_ratio    string                  `json:"sample_aspect_ratio,omitempty"`
	Sample_fmt             string                  `json:"sample_fmt,omitempty"`
	Sample_rate            string                  `json:"sample_rate,omitempty"`
	SilenceData            []SilenceSegment        `json:"silence_segments,omitempty"`
	Start_pts              int                     `json:"start_pts,omitempty"`
	Start_time             string                  `json:"start_time,omitempty"`
	Time_base              string                  `json:"time_base,omitempty"`
	Width                  int                     `json:"width,omitempty"`
	Side_data_list         []Side_data_list_struct `json:"side_data_list,omitempty"`
	Tags                   map[string]string       `json:"tags,omitempty"`
	Disposition            map[string]int          `json:"disposition,omitempty"`
}

type Side_data_list_struct struct {
	Side_data map[string]string
}

type SilenceSegment struct {
	SilenceStart    float64 `json:"start,omitempty"`
	SilenceEnd      float64 `json:"end,omitempty"`
	SilenceDuration float64 `json:"duration,omitempty"`
	LoudnessBorder  float64 `json:"loudness_border,omitempty"`
}

func (db *DBjson) Save(e *Entry) error {
	dir := db.Dir()
	path := dir + e.File + ".json"
	//fmt.Println("save to", path)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf("save failed: %v", err)
	}
	defer f.Close()
	bt, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return fmt.Errorf("save failed: marshaling: %v", err)
	}
	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("save failed: truncate file: %v", err)
	}
	_, err = f.Write(bt)
	if err != nil {
		return fmt.Errorf("save failed: write to file: %v", err)
	}
	//fmt.Printf("bytes writen: %v (%v)\n", wr, path)
	return nil
}

func (db *DBjson) Load(file string) (*Entry, error) {
	key := sourceKey(file)
	dir := db.Dir()
	path := dir + key + ".json"
	e := &Entry{}
	bt, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load failed: %v", err)
	}
	err = json.Unmarshal(bt, e)
	if err != nil {
		return nil, fmt.Errorf("load failed: %v", err)
	}
	return e, nil
}

func sourceKey(file string) string {
	base := filepath.Base(file)
	parts := strings.Split(base, "--")
	return parts[len(parts)-1]
}
