package stream

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DASH = iota
	HLS
)

var NameFMap = map[int]string{
	DASH: "DASH",
	HLS:  "HLS",
}

var MethodFMap = map[int]string{
	DASH: "index.mpd",
	HLS:  "index.m3u8",
}
var VidEXT = []string{".mp4", ".webm", ".ogg", ".wmv", ".mov", ".avchd", ".av1"}

type VideoEncoder struct {
	InputFile  string
	Codec      string
	StreamType int
	OutputDir  string
	OutputFile string
	Audio      bool
	Command    string
}

func (v *VideoEncoder) New(input_file, codec string, stream_type interface{}) error {
	v.InputFile = input_file
	v.Codec = codec
	if stream_type == DASH || stream_type == HLS {
		v.StreamType = stream_type.(int)
	} else if strType, ok := stream_type.(string); ok && strings.ToLower(strType) == "dash" {
		v.StreamType = DASH
	} else if strType, ok := stream_type.(string); ok && strings.ToLower(strType) == "hls" {
		v.StreamType = HLS
	} else {
		return errors.New("Invalid stream type")
	}

	err := v.CheckAll()
	if err != nil {
		return err
	}

	v.SetOutput()
	v.SetCommand()

	return nil
}

func (v *VideoEncoder) Print() {
	fmt.Println("InputFile:", v.InputFile)
	fmt.Println("Codec:", v.Codec)
	fmt.Println("StreamType:", v.StreamType, "|", NameFMap[v.StreamType])
	fmt.Println("OutputDir:", v.OutputDir)
	fmt.Println("OutputFile:", v.OutputFile)
	fmt.Println("Audio:", v.Audio)
	fmt.Println("Command:", v.Command)
}

func (v *VideoEncoder) SetOutput() {
	v.SetOutputDir()
	v.SetOutputFile()
	v.CheckAudio()
}

func (v *VideoEncoder) SetOutputDir() {
	dname, fname := filepath.Split(v.InputFile)
	ext := filepath.Ext(fname)
	fname = strings.TrimSuffix(fname, ext)
	v.OutputDir = filepath.Join(dname, fname)
	os.MkdirAll(v.OutputDir, os.ModePerm)
}

func (v *VideoEncoder) SetOutputFile() {
	v.OutputFile = fmt.Sprintf("%s/%s", v.OutputDir, MethodFMap[v.StreamType])
}

func (v *VideoEncoder) SetCommand() {
	switch v.StreamType {
	case DASH:
		v.DASHcmd()
	case HLS:
		v.HLScmd()
	}
}

func (v *VideoEncoder) CheckAll() error {
	_, fname := filepath.Split(v.InputFile)
	if !CheckEXT(fname) {
		return errors.New("Invalid file extension")
	}

	if !v.CheckCodec() {
		return errors.New("Invalid codec")
	}

	return nil
}

func (v *VideoEncoder) CheckCodec() bool {
	output, err := GetCMD("ffmpeg", "-h", "encoder="+v.Codec)
	if err != nil || output == "" {
		return false
	}

	if Contains(output, "is not recognized by FFmpeg") {
		return false
	}

	return true
}

func (v *VideoEncoder) CheckAudio() {
	has_audio, err := GetCMD("ffprobe", "-i", v.InputFile, "-show_streams", "-select_streams", "a", "-loglevel", "error")
	if err != nil || has_audio == "" {
		v.Audio = false
		return
	}

	v.Audio = true
}

func CheckEXT(fname string) bool {
	ext := filepath.Ext(fname)
	for _, e := range VidEXT {
		if e == ext {
			return true
		}
	}

	return false
}
