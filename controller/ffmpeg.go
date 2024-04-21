package controller

import (
	"ffmpeg-clipper/ffmpeg"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CheckFFmpeg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := ffmpeg.CheckFFmpeg()

	if !result.FFmpegExists || !result.FFprobeExists {
		handleResponseError(w, "controller.CheckFFmpeg: ffmpeg or ffprobe is not present locally or on the path")
	}
}
