package stream

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func (v *VideoEncoder) Encode() error {
	err := v.ValidateEncode()
	if err != nil {
		return err
	}

	cmd := exec.Command("bash", "-c", v.Command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error running command:", err)
	}

	return nil
}

func (v *VideoEncoder) ValidateEncode() error {
	if v.InputFile == "" {
		log.Println("Error: Input file not set")
		return errors.New("Input file not set")
	}
	if v.Codec == "" {
		log.Println("WARNING: Codec not set, using default codec libx264")
		v.Codec = "libx264" //default codec, should work on most ffmpeg installations
	}
	if v.StreamType != DASH && v.StreamType != HLS {
		log.Println("Error: Invalid stream type")
		return errors.New("Invalid stream type, must be DASH(0) or HLS(1)")
	}
	if v.Command == "" {
		log.Println("WARNING: Command not set, setting ffmpeg command")
		v.SetCommand()
	}
	if !v.Audio {
		v.CheckAudio()
	}
	if _, err := os.Stat(v.OutputDir); os.IsNotExist(err) {
		log.Println("WARNING: Output directory does not exist, creating directory")
		os.MkdirAll(v.OutputDir, os.ModePerm)
	}

	file_dir := filepath.Dir(v.OutputFile)
	if file_dir != v.OutputDir {
		log.Println("WARNING: Output file directory does not match output directory, creating directory")
	}
	if _, err := os.Stat(file_dir); os.IsNotExist(err) {
		log.Println("WARNING: Output file directory does not exist, creating directory")
		os.MkdirAll(file_dir, os.ModePerm)
	}

	return nil
}

func (v *VideoEncoder) DASHcmd() {
	audio_cmd := "-c:a libopus -b:a 128k"
	segment_cmd := `-dash_segment_type mp4 -adaptation_sets "id=0,streams=v id=1,streams=a"`
	if !v.Audio {
		audio_cmd = ""
		segment_cmd = `-dash_segment_type mp4 -adaptation_sets "id=0,streams=v"`
	}
	v.Command = fmt.Sprintf(`ffmpeg -i %s \
  -map 0 -c:v %s -b:v 1000k -keyint_min 150 -g 150 -sc_threshold 0 %s \
  -f dash -seg_duration 4 -use_template 1 -use_timeline 1 -init_seg_name 'init-$RepresentationID$.m4s' \
  -media_seg_name 'chunk-$RepresentationID$-$Number$.m4s' \
  %s \
  %s`,
		v.InputFile, v.Codec, audio_cmd, segment_cmd, v.OutputFile)
}

func (v *VideoEncoder) HLScmd() {
	audio_cmd := "-c:a aac -b:a 128k"
	segment_cmd := `-var_stream_map "v:0,a:0"`
	if !v.Audio {
		audio_cmd = ""
		segment_cmd = `-var_stream_map "v:0"`
	}
	v.Command = fmt.Sprintf(`ffmpeg -i %s \
  -map 0 -c:v %s -b:v 1000k -keyint_min 150 -g 150 -sc_threshold 0 %s \
  -f hls \
  -hls_time 4 \
  -hls_playlist_type vod \
  -hls_segment_filename %s/segment_%%03d.ts \
  -master_pl_name /master.m3u8 \
  %s \
  %s`,
		v.InputFile, v.Codec, audio_cmd, v.OutputDir, segment_cmd, v.OutputFile)
}

func GetCMD(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func Contains(source, target string) bool {
	length := len(target)
	if length > len(source) {
		return false
	}

	for i := 0; i <= len(source)-length; i++ {
		if source[i:i+length] == target {
			return true
		}
	}

	return false
}
