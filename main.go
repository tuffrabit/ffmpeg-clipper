package main

import (
	"ffmpeg-clipper/common"
	"ffmpeg-clipper/controller"
	"ffmpeg-clipper/html"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

var frontendUri string
var indexHtmlContent string

func main() {
	logFile := setupLog()
	if logFile != nil {
		defer logFile.Close()
	}

	err := controller.InitConfig()
	if err != nil {
		log.Fatalf("main.main: could not init config: %v", err)
	}

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/mainjs", GetMainJs)
	router.GET("/checkffmpeg", controller.CheckFFmpeg)
	router.GET("/getavailablevideos", controller.GetAvailableVideos)
	router.GET("/getvideodetails/:name", controller.GetVideoDetails)
	router.GET("/getconfig", controller.GetConfig)
	router.POST("/saveprofile", controller.SaveProfile)
	router.DELETE("/deleteprofile", controller.DeleteProfile)
	router.POST("/playvideo", controller.PlayVideo)
	router.DELETE("/deletevideo", controller.DeleteVideo2)
	router.POST("/clipvideo", controller.ClipVideo)
	router.ServeFiles("/streamvideo/*filepath", http.Dir("."))

	router.GET("/active", controller.Active)
	router.GET("/static/pico.min.css", GetPicoCss)
	router.GET("/static/htmx.min.js", GetHtmxJs)
	router.GET("/static/modal.js", GetModalJs)
	router.GET("/static/main.js", GetMainJs)
	router.GET("/videos.html", controller.GetVideos)
	router.GET("/getvideoplayer.html", controller.GetVideoPlayer)
	router.POST("/deletevideo.html", controller.DeleteVideo)
	router.GET("/profiles.html", controller.GetProfiles)
	router.GET("/getprofile.html", controller.GetProfile)
	router.GET("/getencodersettings.html", controller.GetEncoderSettings)
	router.POST("/saveprofile.html", controller.SaveProfile)
	router.POST("/deleteprofile.html", controller.DeleteProfile)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("main.main: could not start http listener: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	portStr := strconv.Itoa(port)
	frontendUri = "http://localhost:" + portStr
	wsUri := "ws://localhost:" + portStr
	indexHtmlContent, err = html.GetIndexHtmlContent(frontendUri, wsUri)
	if err != nil {
		log.Fatalf("main.main: could not load index html: %v", err)
	}

	log.Println("Using port:", port)

	go func(listener net.Listener, router *httprouter.Router) {
		log.Fatal(http.Serve(listener, router))
	}(listener, router)

	cmd := exec.Command("explorer", frontendUri)
	common.RunSystemCommand(cmd)

	common.SetClientLastUpdateTime(time.Now())
	common.MonitorClient()
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, indexHtmlContent)
}

func GetMainJs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "text/html")

	content, err := html.GetMainJsContent(frontendUri)
	if err != nil {
		message := fmt.Sprintf("main.GetMainJs: could not get js template content: %v", err)
		log.Println(message)
		header.Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"error:\""+message+"}")
		return
	}

	fmt.Fprint(w, content)
}

func GetPicoCss(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "text/css")

	content, err := html.GetPicoCssContent()
	if err != nil {
		message := fmt.Sprintf("main.GetPicoCss: could not get css template content: %v", err)
		log.Println(message)
		header.Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"error:\""+message+"}")
		return
	}

	fmt.Fprint(w, content)
}

func GetHtmxJs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "text/html")

	content, err := html.GetHtmxJsContent()
	if err != nil {
		message := fmt.Sprintf("main.GetHtmxJs: could not get js template content: %v", err)
		log.Println(message)
		header.Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"error:\""+message+"}")
		return
	}

	fmt.Fprint(w, content)
}

func GetModalJs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "text/html")

	content, err := html.GetModalJsContent()
	if err != nil {
		message := fmt.Sprintf("main. GetModalJs: could not get js template content: %v", err)
		log.Println(message)
		header.Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"error:\""+message+"}")
		return
	}

	fmt.Fprint(w, content)
}

func setupLog() *os.File {
	log.SetOutput(os.Stdout)

	useFile := true
	logFile, err := os.OpenFile("ffmpeg-clipper.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		useFile = false
		log.Println("could not open ffmpeg-clipper.log file for write")
	}

	if useFile {
		multi := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multi)
		return logFile
	}

	return nil
}
