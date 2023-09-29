package controller

import (
	"encoding/json"
	"ffmpeg-clipper/ffmpeg"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CheckFFmpeg(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	result := ffmpeg.CheckFFmpeg()
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
