package controller

import (
	"ffmpeg-clipper/ffmpeg"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CheckFFmpeg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := ffmpeg.CheckFFmpeg()
	//header := w.Header()
	//header.Set("Content-Type", "application/json")

	/*resultBytes, err := json.Marshal(result)
	if err != nil {
		message := fmt.Sprintf("controller.CheckFFmpeg: could not marshal struct to json: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
	} else {
		fmt.Fprint(w, string(resultBytes))
	}*/

	if !result.FFmpegExists || !result.FFprobeExists {
		handleResponseError(w, "controller.CheckFFmpeg: ffmpeg or ffprobe is not present locally or on the path")
	}
}
