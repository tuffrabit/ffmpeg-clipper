package controller

import (
	"encoding/json"
	"ffmpeg-clipper/common"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func GetAvailableVideos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := struct {
		AvailableVideos []string
	}{}

	header := w.Header()
	header.Set("Content-Type", "application/json")

	var availableVideos []string
	dirEntries, err := os.ReadDir("./")
	if err != nil {
		message := fmt.Sprintf("controller.GetAvailableVideos: could not list directory: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	allowedExtensions := map[string]struct{}{
		".mp4": {},
		".mkv": {},
		".avi": {},
		".flv": {},
		".mov": {},
		".wmv": {},
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			fileExtension := filepath.Ext(entry.Name())

			_, ok := allowedExtensions[fileExtension]
			if ok {
				availableVideos = append(availableVideos, entry.Name())
			}
		}
	}

	result.AvailableVideos = availableVideos

	resultBytes, err := json.Marshal(result)
	if err != nil {
		message := fmt.Sprintf("controller.GetAvailableVideos: could not marshal struct to json: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
	} else {
		fmt.Fprint(w, string(resultBytes))
	}
}

func GetVideoDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	video := ps.ByName("name")

	header := w.Header()
	header.Set("Content-Type", "application/json")

	_, err := os.Stat(video)
	if err != nil {
		message := fmt.Sprintf("controller.GetVideoDetails: %v does not exist: %v", video, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	cmd := exec.Command(
		FfprobePath,
		"-v",
		"error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=width,height",
		"-of",
		"csv=s=x:p=0",
		video,
	)
	output, err := common.RunSystemCommand(cmd)
	if err != nil {
		message := fmt.Sprintf("controller.GetVideoDetails: ffprobe failed\nstderr: %v\nerr: %v", output, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	output = strings.TrimSuffix(output, "\r\n")
	output = strings.TrimSuffix(output, "\n")
	output = strings.TrimSuffix(output, "\r")

	result := struct {
		Resolution string
	}{
		Resolution: output,
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		message := fmt.Sprintf("controller.GetVideoDetails: could not marshal struct to json: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
	} else {
		fmt.Fprint(w, string(resultBytes))
	}
}

func PlayVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf("controller.PlayVideo: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	payloadJson := struct {
		Video           string `json:"video"`
		AlternatePlayer string `json:"alternatePlayer"`
	}{}
	err = json.Unmarshal(bodyBytes, &payloadJson)
	if err != nil {
		message := fmt.Sprintf("controller.PlayVideo: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	_, err = os.Stat(payloadJson.Video)
	if err != nil {
		message := fmt.Sprintf("controller.PlayVideo: %v does not exist: %v", payloadJson.Video, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	err = playVideoFile(payloadJson.Video, payloadJson.AlternatePlayer)
	if err != nil {
		message := fmt.Sprintf("controller.PlayVideo: could not play video: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	fmt.Fprint(w, "{}")
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf("controller.DeleteVideo: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	payloadJson := struct {
		Video string `json:"video"`
	}{}
	err = json.Unmarshal(bodyBytes, &payloadJson)
	if err != nil {
		message := fmt.Sprintf("controller.DeleteVideo: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	_, err = os.Stat(payloadJson.Video)
	if err != nil {
		message := fmt.Sprintf("controller.DeleteVideo: %v does not exist: %v", payloadJson.Video, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	err = os.Remove(payloadJson.Video)
	if err != nil {
		message := fmt.Sprintf("controller.DeleteVideo: could not delete video file: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	fmt.Fprint(w, "{}")
}

func playVideoFile(videoFileName string, alternatePlayerPath string) error {
	var cmd *exec.Cmd

	if alternatePlayerPath != "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("controller.playVideoFile: could not get current dir: %w", err)
		}

		video := filepath.Join(currentDir, videoFileName)
		cmd = exec.Command(alternatePlayerPath, video)
	} else {
		cmd = exec.Command(FfplayPath,
			"-nostats",
			"-hide_banner",
			"-loglevel",
			"error",
			videoFileName,
		)
	}

	cmdOutput, err := common.RunSystemCommand(cmd)
	if err != nil {
		if err.Error() != cmdOutput {
			return fmt.Errorf("controller.playVideoFile: could not play video\nstderr: %verr: %w", cmdOutput, err)
		} else {
			return fmt.Errorf("controller.playVideoFile: could not play video\nerr: %w", err)
		}
	}

	return nil
}
