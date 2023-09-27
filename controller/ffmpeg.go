package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

var FfmpegPath string
var FfplayPath string
var FfprobePath string

func CheckFFmpeg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	FfmpegPath = "./ffmpeg.exe"
	FfplayPath = "./ffplay.exe"
	FfprobePath = "./ffprobe.exe"
	result := struct {
		FFmpegExists  bool
		FFplayExists  bool
		FFprobeExists bool
	}{
		FFmpegExists:  false,
		FFplayExists:  false,
		FFprobeExists: false,
	}

	pathVar := os.Getenv("PATH")
	paths := strings.Split(pathVar, ";")
	pathSeparator := string(os.PathSeparator)

	for _, pathEntry := range paths {
		if strings.Contains(pathEntry, "ffmpeg") {
			FfmpegPath = fmt.Sprintf("%v%vffmpeg.exe", pathEntry, pathSeparator)
			FfplayPath = fmt.Sprintf("%v%vffplay.exe", pathEntry, pathSeparator)
			FfprobePath = fmt.Sprintf("%v%vffprobe.exe", pathEntry, pathSeparator)
			break
		}
	}

	_, err := os.Stat(FfmpegPath)
	if err == nil {
		result.FFmpegExists = true
	}

	_, err = os.Stat(FfplayPath)
	if err == nil {
		result.FFplayExists = true
	}

	_, err = os.Stat(FfprobePath)
	if err == nil {
		result.FFprobeExists = true
	}

	header := w.Header()
	header.Set("Content-Type", "application/json")

	resultBytes, err := json.Marshal(result)
	if err != nil {
		message := fmt.Sprintf("controller.CheckFFmpeg: could not marshal struct to json: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
	} else {
		fmt.Fprint(w, string(resultBytes))
	}
}
