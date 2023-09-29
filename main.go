package main

import (
	"ffmpeg-clipper/common"
	"ffmpeg-clipper/controller"
	"ffmpeg-clipper/html"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

var frontendUri string
var indexHtmlContent string

func main() {
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
	router.DELETE("/deletevideo", controller.DeleteVideo)
	router.POST("/clipvideo", controller.ClipVideo)
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
		common.RunSystemCommand(cmd)
	}(port)

	log.Fatal(http.Serve(listener, router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, indexHtmlContent)
}

func GetMainJs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "text/html")

	mainJsContent, err := html.GetMainJsContent(frontendUri)
	if err != nil {
		message := fmt.Sprintf("main.GetMainJs: could not get js template content: %v", err)
		log.Println(message)
		header.Set("Content-Type", "application/json")
		fmt.Fprint(w, "{\"error:\""+message+"}")
		return
	}

	fmt.Fprint(w, mainJsContent)
}
