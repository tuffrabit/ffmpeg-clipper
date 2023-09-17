package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"ffmpeg-clipper/config"
	"ffmpeg-clipper/html"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

var configJson *config.ConfigJson
var frontendUri string
var indexHtmlContent string
var ffmpegPath string
var ffplayPath string
var ffprobePath string

func main() {
	var err error

	configJson, err = config.GetConfig()
	if err != nil {
		log.Fatalf("main.main: could not get config: %v", err)
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/checkffmpeg", CheckFFmpeg)
	router.GET("/getavailablevideos", GetAvailableVideos)
	router.GET("/getvideodetails/:name", GetVideoDetails)
	router.GET("/getconfig", GetConfig)
	router.POST("/saveprofile", SaveProfile)
	router.DELETE("/deleteprofile", DeleteProfile)
	router.POST("/playvideo", PlayVideo)
	router.DELETE("/deletevideo", DeleteVideo)
	router.POST("/clipvideo", ClipVideo)
	router.ServeFiles("/streamvideo/*filepath", http.Dir("."))

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("main.main: could not start http listener: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	frontendUri = "http://localhost:" + strconv.Itoa(port)
	indexHtmlContent, err = html.GetIndex2HtmlContent(frontendUri)
	if err != nil {
		log.Fatalf("main.main: could not load index html: %v", err)
	}

	fmt.Println("Using port:", port)

	go func(port int) {
		time.Sleep(2 * time.Second)

		cmd := exec.Command("explorer", frontendUri)
		runSystemCommand(cmd)
	}(port)

	log.Fatal(http.Serve(listener, router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, indexHtmlContent)
}

func CheckFFmpeg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ffmpegPath = "./ffmpeg.exe"
	ffplayPath = "./ffplay.exe"
	ffprobePath = "./ffprobe.exe"
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
			ffmpegPath = fmt.Sprintf("%v%vffmpeg.exe", pathEntry, pathSeparator)
			ffplayPath = fmt.Sprintf("%v%vffplay.exe", pathEntry, pathSeparator)
			ffprobePath = fmt.Sprintf("%v%vffprobe.exe", pathEntry, pathSeparator)
			break
		}
	}

	_, err := os.Stat(ffmpegPath)
	if err == nil {
		result.FFmpegExists = true
	}

	_, err = os.Stat(ffplayPath)
	if err == nil {
		result.FFplayExists = true
	}

	_, err = os.Stat(ffprobePath)
	if err == nil {
		result.FFprobeExists = true
	}

	header := w.Header()
	header.Set("Content-Type", "application/json")

	resultBytes, err := json.Marshal(result)
	if err != nil {
		message := fmt.Sprintf("main.CheckFFmpeg: could not marshal struct to json: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
	} else {
		fmt.Fprint(w, string(resultBytes))
	}
}

func GetAvailableVideos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := struct {
		AvailableVideos []string
	}{}

	header := w.Header()
	header.Set("Content-Type", "application/json")

	var availableVideos []string
	dirEntries, err := os.ReadDir("./")
	if err != nil {
		message := fmt.Sprintf("main.GetAvailableVideos: could not list directory: %v", err)
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
		message := fmt.Sprintf("main.GetAvailableVideos: could not marshal struct to json: %v", err)
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
		message := fmt.Sprintf("main.GetVideoDetails: %v does not exist: %v", video, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	cmd := exec.Command(
		"./ffprobe.exe",
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
	output, err := runSystemCommand(cmd)
	if err != nil {
		message := fmt.Sprintf("main.GetVideoDetails: ffprobe failed\nstderr: %v\nerr: %v", output, err)
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
		message := fmt.Sprintf("main.GetVideoDetails: could not marshal struct to json: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
	} else {
		fmt.Fprint(w, string(resultBytes))
	}
}

func GetConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	resultBytes, err := json.Marshal(configJson)
	if err != nil {
		message := fmt.Sprintf("main.GetConfig: could not marshal struct to json: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
	} else {
		fmt.Fprint(w, string(resultBytes))
	}
}

func SaveProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf("main.SaveProfile: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	payloadJson := config.ClipProfileJson{}
	err = json.Unmarshal(bodyBytes, &payloadJson)
	if err != nil {
		message := fmt.Sprintf("main.SaveProfile: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	err = config.SaveProfile(&payloadJson)
	if err != nil {
		message := fmt.Sprintf("main.SaveProfile: could not save profile: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	configJson, err = config.GetConfig()
	if err != nil {
		message := fmt.Sprintf("main.SaveProfile: could not load updated config: %v", err)
		fmt.Fprint(w, generateErrorJson(message))
		log.Fatalln(message)
	}

	fmt.Fprint(w, "{}")
}

func DeleteProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf("main.DeleteProfile: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	payloadJson := struct {
		ProfileName string `json:"profileName"`
	}{}
	err = json.Unmarshal(bodyBytes, &payloadJson)
	if err != nil {
		message := fmt.Sprintf("main.DeleteProfile: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	err = config.DeleteProfile(payloadJson.ProfileName)
	if err != nil {
		message := fmt.Sprintf("main.DeleteProfile: could not delete profile: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	configJson, err = config.GetConfig()
	if err != nil {
		message := fmt.Sprintf("main.DeleteProfile: could not load updated config: %v", err)
		fmt.Fprint(w, generateErrorJson(message))
		log.Fatalln(message)
	}

	fmt.Fprint(w, "{}")
}

func PlayVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf("main.PlayVideo: could not read request body: %v", err)
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
		message := fmt.Sprintf("main.PlayVideo: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	_, err = os.Stat(payloadJson.Video)
	if err != nil {
		message := fmt.Sprintf("main.PlayVideo: %v does not exist: %v", payloadJson.Video, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	err = playVideoFile(payloadJson.Video, payloadJson.AlternatePlayer)
	if err != nil {
		message := fmt.Sprintf("main.PlayVideo: could not play video: %v", err)
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
		message := fmt.Sprintf("main.DeleteVideo: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	payloadJson := struct {
		Video string `json:"video"`
	}{}
	err = json.Unmarshal(bodyBytes, &payloadJson)
	if err != nil {
		message := fmt.Sprintf("main.DeleteVideo: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	_, err = os.Stat(payloadJson.Video)
	if err != nil {
		message := fmt.Sprintf("main.DeleteVideo: %v does not exist: %v", payloadJson.Video, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	err = os.Remove(payloadJson.Video)
	if err != nil {
		message := fmt.Sprintf("main.DeleteVideo: could not delete video file: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	fmt.Fprint(w, "{}")
}

func ClipVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf("main.ClipVideo: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	payloadJson := struct {
		Video           string  `json:"video"`
		StartTime       string  `json:"startTime"`
		EndTime         string  `json:"endTime"`
		ScaleFactor     float32 `json:"scaleFactor"`
		EncodingPreset  string  `json:"encodingPreset"`
		QualityTarget   int     `json:"qualityTarget"`
		Saturation      float32 `json:"saturation"`
		Contrast        float32 `json:"contrast"`
		Brightness      float32 `json:"brightness"`
		Gamma           float32 `json:"gamma"`
		PlayAfter       bool    `json:"playAfter"`
		AlternatePlayer string  `json:"alternatePlayer"`
	}{}

	err = json.Unmarshal(bodyBytes, &payloadJson)
	if err != nil {
		message := fmt.Sprintf("main.ClipVideo: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	videoExtension := filepath.Ext(payloadJson.Video)
	videoName := payloadJson.Video[:len(payloadJson.Video)-len(videoExtension)]
	newVideoName := fmt.Sprintf("%v_clip%v%v", videoName, getRandomString(), videoExtension)

	cmd := exec.Command(ffmpegPath,
		"-nostats",
		"-hide_banner",
		"-loglevel",
		"error",
		"-ss",
		payloadJson.StartTime,
		"-to",
		payloadJson.EndTime,
		"-i",
		payloadJson.Video,
		"-c:v",
		"libx264",
		"-preset",
		payloadJson.EncodingPreset,
		"-crf",
		strconv.Itoa(payloadJson.QualityTarget),
		"-vf",
		fmt.Sprintf("scale=iw/%v:-1:flags=bicubic,eq=saturation=%v:contrast=%v:brightness=%v:gamma=%v",
			payloadJson.ScaleFactor,
			payloadJson.Saturation,
			payloadJson.Contrast,
			payloadJson.Brightness,
			payloadJson.Gamma,
		),
		newVideoName,
	)
	cmdOutput, err := runSystemCommand(cmd)
	if err != nil {
		_, err = os.Stat(newVideoName)
		if err == nil {
			os.Remove(newVideoName)
		}

		fmt.Fprint(w, generateErrorJson(cmdOutput))
		return
	}

	_, err = os.Stat(newVideoName)
	if err != nil {
		message := fmt.Sprintf("main.ClipVideo: new video clip %v was not created: %v", newVideoName, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	fmt.Fprint(w, "{\"newVideoName\": \""+newVideoName+"\"}")
}

func playVideoFile(videoFileName string, alternatePlayerPath string) error {
	var cmd *exec.Cmd

	if alternatePlayerPath != "" {
		currentDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("main.playVideoFile: could not get current dir: %w", err)
		}

		video := filepath.Join(currentDir, videoFileName)
		cmd = exec.Command(alternatePlayerPath, video)
	} else {
		cmd = exec.Command(ffplayPath,
			"-nostats",
			"-hide_banner",
			"-loglevel",
			"error",
			videoFileName,
		)
	}

	cmdOutput, err := runSystemCommand(cmd)
	if err != nil {
		if err.Error() != cmdOutput {
			return fmt.Errorf("main.playVideoFile: could not play video\nstderr: %verr: %w", cmdOutput, err)
		} else {
			return fmt.Errorf("main.playVideoFile: could not play video\nerr: %w", err)
		}
	}

	return nil
}

func generateErrorJson(errorMessage string) string {
	errorJson := struct {
		Error string
	}{
		Error: errorMessage,
	}

	errorJsonBytes, err := json.Marshal(errorJson)
	if err != nil {
		log.Printf("main.generateErrorJson: Failed to marshal error message to json: %v", err)
		log.Printf("original error message: %v", errorMessage)
		return "{\"error\": \"Could not generate full error message. See console for details.\"}"
	}

	return string(errorJsonBytes)
}

func runSystemCommand(cmd *exec.Cmd) (string, error) {
	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	cmdString := cmd.String()

	log.Printf("Running %v\n", cmdString)

	err := cmd.Run()
	outString := cmdOut.String()
	errString := cmdErr.String()
	if err != nil || errString != "" {
		if outString != "" {
			log.Printf("%v stdout: %v", cmdString, outString)
		}

		if errString != "" {
			log.Printf("%v stderr: %v", cmdString, errString)

			if err == nil {
				err = errors.New(errString)
			}
		}

		log.Println(err)

		return errString, err
	} else {
		if outString != "" {
			log.Printf("%v stdout: %v", cmdString, outString)
		}

		return outString, nil
	}
}

func getRandomString() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X", b)
}
