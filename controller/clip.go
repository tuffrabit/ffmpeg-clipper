package controller

import (
	"encoding/json"
	"ffmpeg-clipper/common"
	"ffmpeg-clipper/config"
	ffmpegEncoder "ffmpeg-clipper/encoder"
	"ffmpeg-clipper/ffmpeg"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func ClipVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	var bodyJson map[string]interface{}

	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	encoder, ok := bodyJson["encoder"].(string)
	if !ok {
		message := "controller.ClipVideo: could not determine encoder"
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	var newVideoName string

	switch config.EncoderType(encoder) {
	case config.Libx264EncoderType:
		newVideoName, err = ffmpegEncoder.ClipLibx264(bodyJson)
	case config.Libx265EncoderType:
		newVideoName, err = ffmpegEncoder.ClipLibx265(bodyJson)
	case config.LibaomAv1EncoderType:
		newVideoName, err = ffmpegEncoder.ClipLibaomAv1(bodyJson)
	case config.NvencH264EncoderType:
		newVideoName, err = ffmpegEncoder.ClipNvencH264(bodyJson)
	case config.NvencHevcEncoderType:
		newVideoName, err = ffmpegEncoder.ClipNvencHevc(bodyJson)
	case config.IntelH264EncoderType:
		newVideoName, err = ffmpegEncoder.ClipIntelH264(bodyJson)
	case config.IntelHevcEncoderType:
		newVideoName, err = ffmpegEncoder.ClipIntelHevc(bodyJson)
	default:
		err = fmt.Errorf("%v is not a valid encoder", encoder)
	}

	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: could not create clip: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	_, err = os.Stat(newVideoName)
	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: new video clip %v was not created: %v", newVideoName, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	ffprobeErr := false

	cmd := exec.Command(
		ffmpeg.FfprobePath,
		"-v",
		"error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=width,height",
		"-of",
		"csv=s=x:p=0",
		newVideoName,
	)
	output, err := common.RunSystemCommand(cmd)
	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: ffprobe failed\nstderr: %v\nerr: %v", output, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		ffprobeErr = true
	}

	output = strings.TrimSuffix(output, "\r\n")
	output = strings.TrimSuffix(output, "\n")
	output = strings.TrimSuffix(output, "\r")

	if output == "" {
		message := fmt.Sprintf("controller.ClipVideo: generated clip %v is invalid", newVideoName)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		ffprobeErr = true
	}

	if ffprobeErr {
		_, err = os.Stat(newVideoName)
		if err == nil {
			os.Remove(newVideoName)
		}
		return
	}

	fmt.Fprint(w, "{\"newVideoName\": \""+newVideoName+"\"}")
}
