package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
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

var frontendUri string
var indexHtmlContent string
var ffmpegPath string
var ffplayPath string

func main() {
	var err error
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/checkffmpeg", CheckFFmpeg)
	router.GET("/getavailablevideos", GetAvailableVideos)
	router.GET("/playvideo/:name", PlayVideo)
	router.POST("/clipvideo", ClipVideo)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	frontendUri = "http://localhost:" + strconv.Itoa(port)
	indexHtmlContent, err = html.GetIndexHtmlContent(frontendUri)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Using port:", port)

	go func(port int) {
		time.Sleep(2 * time.Second)

		cmd := exec.Command("explorer", frontendUri)
		runSystemCommand(cmd, false)
	}(port)

	log.Fatal(http.Serve(listener, router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, indexHtmlContent)
}

func CheckFFmpeg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ffmpegPath = "./ffmpeg.exe"
	ffplayPath = "./ffplay.exe"
	result := struct {
		FFmpegExists bool
		FFplayExists bool
	}{
		FFmpegExists: false,
		FFplayExists: false,
	}

	pathVar := os.Getenv("PATH")
	paths := strings.Split(pathVar, ";")
	pathSeparator := string(os.PathSeparator)

	for _, pathEntry := range paths {
		if strings.Contains(pathEntry, "ffmpeg") {
			ffmpegPath = fmt.Sprintf("%v%vffmpeg.exe", pathEntry, pathSeparator)
			ffplayPath = fmt.Sprintf("%v%vffplay.exe", pathEntry, pathSeparator)
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

	header := w.Header()
	header.Set("Content-Type", "application/json")

	resultBytes, err := json.Marshal(result)
	if err != nil {
		log.Printf("main.CheckFFmpeg: could not marshal struct to json: %v", err)
		fmt.Printf("{\"error\": \"main.CheckFFmpeg: could not marshal struct to json\"}")
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
		log.Printf("main.GetAvailableVideos: could not list directory: %v", err)
		fmt.Fprint(w, "{\"error\": \"main.GetAvailableVideos: could not list directory\"}")
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
		log.Printf("main.GetAvailableVideos: could not marshal struct to json: %v", err)
		fmt.Fprint(w, "{\"error\": \"main.GetAvailableVideos: could not marshal struct to json\"}")
	} else {
		fmt.Fprint(w, string(resultBytes))
	}
}

func PlayVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	video := ps.ByName("name")

	header := w.Header()
	header.Set("Content-Type", "application/json")

	_, err := os.Stat(video)
	if err != nil {
		log.Printf("main.PlayVideo: %v does not exist: %v", video, err)
		fmt.Fprint(w, "{\"error\": \"main.PlayVideo: requested video file does not exist\"}")
		return
	}

	cmd := exec.Command(ffplayPath, video)
	runSystemCommand(cmd, false)
	fmt.Fprint(w, "{}")
}

func ClipVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("main.ClipVideo: could not read request body: %v", err)
		fmt.Fprint(w, "{\"error\": \"main.ClipVideo: could not read request body\"}")
		return
	}

	payloadJson := struct {
		Video          string `json:"video"`
		StartTime      string `json:"startTime"`
		EndTime        string `json:"endTime"`
		ScaleFactor    string `json:"scaleFactor"`
		EncodingPreset string `json:"encodingPreset"`
		QualityTarget  string `json:"qualityTarget"`
		Saturation     string `json:"saturation"`
		Contrast       string `json:"contrast"`
		Brightness     string `json:"brightness"`
		Gamma          string `json:"gamma"`
		PlayAfter      bool   `json:"playAfter"`
	}{}

	err = json.Unmarshal(bodyBytes, &payloadJson)
	if err != nil {
		log.Printf("main.ClipVideo: could not json marshal request body: %v", err)
		fmt.Fprint(w, "{\"error\": \"main.ClipVideo: could not json marshal request body\"}")
		return
	}

	videoExtension := filepath.Ext(payloadJson.Video)
	videoName := payloadJson.Video[:len(payloadJson.Video)-len(videoExtension)]
	newVideoName := fmt.Sprintf("%v_clip%v%v", videoName, getRandomString(), videoExtension)

	cmd := exec.Command(ffmpegPath,
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
		payloadJson.QualityTarget,
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
	runSystemCommand(cmd, false)

	_, err = os.Stat(newVideoName)
	if err != nil {
		log.Printf("main.ClipVideo: new video clip %v was not created: %v", newVideoName, err)
		fmt.Fprint(w, "{\"error\": \"main.ClipVideo: new video clip was not created\"}")
		return
	}

	if payloadJson.PlayAfter {
		cmd := exec.Command(ffplayPath, newVideoName)
		runSystemCommand(cmd, false)
	}

	fmt.Fprint(w, "{\"newVideoName\": \""+newVideoName+"\"}")
}

func runSystemCommand(cmd *exec.Cmd, dieOnError bool) {
	var cmdOut bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr
	cmdString := cmd.String()

	log.Printf("Running %v\n", cmdString)

	err := cmd.Run()
	if err != nil {
		log.Printf("%v stdout: %v", cmdString, cmdOut.String())
		log.Printf("%v stderr: %v", cmdString, cmdErr.String())

		if dieOnError {
			log.Fatal(err)
		} else {
			log.Println(err)
		}
	} else {
		log.Printf("%v stdout: %v", cmdString, cmdOut.String())
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
