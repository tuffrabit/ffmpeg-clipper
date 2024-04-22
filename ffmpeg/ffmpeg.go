package ffmpeg

import (
	"fmt"
	"os"
	"strings"
)

var FfmpegPath string
var FfprobePath string

type Exists struct {
	FFmpegExists  bool
	FFprobeExists bool
}

func CheckFFmpeg() Exists {
	FfmpegPath = "./ffmpeg.exe"
	FfprobePath = "./ffprobe.exe"
	result := Exists{
		FFmpegExists:  false,
		FFprobeExists: false,
	}

	pathVar := os.Getenv("PATH")
	paths := strings.Split(pathVar, ";")
	pathSeparator := string(os.PathSeparator)

	for _, pathEntry := range paths {
		if strings.Contains(pathEntry, "ffmpeg") {
			FfmpegPath = fmt.Sprintf("%v%vffmpeg.exe", pathEntry, pathSeparator)
			FfprobePath = fmt.Sprintf("%v%vffprobe.exe", pathEntry, pathSeparator)
			break
		}
	}

	_, err := os.Stat(FfmpegPath)
	if err == nil {
		result.FFmpegExists = true
	}

	_, err = os.Stat(FfprobePath)
	if err == nil {
		result.FFprobeExists = true
	}

	return result
}
