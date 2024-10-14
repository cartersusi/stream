# HLS and MPEG-DASH Encoding with Go

## About
Commands used for generating hls and dash in [bstore](https://github.com/cartersusi/bstore)

To see how to integrate in a server environment check out.
1. [Upload API](https://github.com/cartersusi/bstore/blob/main/pkg/bstore/upload.go)
2. [File Hosting API](https://github.com/cartersusi/bstore/blob/main/pkg/bstore/serve.go)

## Installation 
```sh
go get github.com/cartersusi/stream
```

## Usage
```go
package main

import (
	"fmt"
	"os"

	"github.com/cartersusi/stream"
)

func dash(input_path, codec string) {
	encoder := stream.VideoEncoder{}
	err := encoder.New(input_path, codec, "dash")
	if err != nil {
		fmt.Println("Error:", err)
	}

	encoder.Print()
	err = encoder.Encode()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func hls(input_path, codec string) {
	encoder := stream.VideoEncoder{
		InputFile:  input_path,
		Codec:      codec,
		StreamType: stream.HLS,
		OutputDir:  "custom_output",
		OutputFile: "custom_output/index.m3u8",
	}

	encoder.CheckAudio() // not necessary, will auto check when encoding
	encoder.SetCommand() // not necessary, will auto set when encoding
	encoder.Print()

	err := encoder.Encode()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func main() {
	input_path := os.Args[1]
	codec := os.Args[2]

	dash(input_path, codec)
	hls(input_path, codec)
}
```