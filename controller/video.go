package controller

import (
	"ffmpeg-clipper/html/templ"
	"ffmpeg-clipper/video"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func GetVideos(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "text/html")

	component, err := templ.GetAvailableVideosSelect(false)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.GetVideos: could not get video list %v", err))
		return
	}

	component.Render(r.Context(), w)
}

func GetVideoPlayer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryValues := r.URL.Query()
	videoPath := queryValues.Get("path")
	header := w.Header()
	header.Set("Content-Type", "text/html")

	videoPlayerComponent, err := templ.GetVideoPlayer(videoPath)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.GetVideoPlayer: could not get video player %v", err))
		return
	}

	/*cmd := exec.Command(
		ffmpeg.FfprobePath,
		"-v",
		"error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=width,height",
		"-of",
		"csv=s=x:p=0",
		videoPath,
	)
	output, err := common.RunSystemCommand(cmd)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.GetVideoPlayer: ffprobe failed\nstderr: %v\nerr: %v", output, err))
		return
	}

	output = strings.TrimSuffix(output, "\r\n")
	output = strings.TrimSuffix(output, "\n")
	output = strings.TrimSuffix(output, "\r")*/

	resolution, err := video.GetVideoResolution(videoPath)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.GetVideoPlayer: could not get video resolution: %v", err))
		return
	}

	videoDetailsComponent, err := templ.GetVideoDetailsOutOfBand(resolution)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.GetVideoPlayer: could not get video details: %v", err))
		return
	}

	videoPlayerComponent.Render(r.Context(), w)
	videoNameComponent := templ.GetVideoNameOutOfBand(videoPath)
	videoNameComponent.Render(r.Context(), w)
	videoDetailsComponent.Render(r.Context(), w)
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	videoPath := r.FormValue("path")
	header := w.Header()
	header.Set("Content-Type", "text/html")

	_, err := os.Stat(videoPath)
	if err == nil {
		log.Printf("deleteing %v", videoPath)
		err = os.Remove(videoPath)
		if err != nil {
			handleResponseError(w, fmt.Sprintf("controller.DeleteVideo: could not delete video file: %v", err))
			return
		}
	}

	availableVideosSelectComponent, err := templ.GetAvailableVideosSelect(true)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.DeleteVideo: could not get video list %v", err))
		return
	}

	videoPlayerComponent, err := templ.GetVideoPlayerOutOfBand("")
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.DeleteVideo: could not get video player %v", err))
		return
	}

	videoDetailsComponent, err := templ.GetVideoDetailsOutOfBand("")
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.DeleteVideo: could not get video details: %v", err))
		return
	}

	availableVideosSelectComponent.Render(r.Context(), w)
	videoNameComponent := templ.GetVideoNameOutOfBand("")
	videoNameComponent.Render(r.Context(), w)
	videoPlayerComponent.Render(r.Context(), w)
	videoDetailsComponent.Render(r.Context(), w)
}
