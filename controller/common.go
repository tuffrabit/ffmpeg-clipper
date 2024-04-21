package controller

import (
	"encoding/json"
	"ffmpeg-clipper/common"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
)

func Active(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		ticker := time.NewTicker(2 * time.Second)

		for range ticker.C {
			// Write
			err := websocket.Message.Send(ws, "ping")
			if err != nil {
				log.Printf("controller.active: websocket send error: %v", err)
				return
			}

			t := time.Now()
			//log.Printf("client ping at %v", t)

			// Read
			msg := ""
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Printf("controller.active: websocket receive error: %v", err)
				return
			}

			if msg == "pong" {
				common.SetClientLastUpdateTime(t)
			}
		}
	}).ServeHTTP(w, r)
}

func handleResponseError(w http.ResponseWriter, message string) {
	log.Println(message)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, generateErrorJson(message))
}

func generateErrorJson(errorMessage string) string {
	errorJson := struct {
		Error string
	}{
		Error: errorMessage,
	}

	errorJsonBytes, err := json.Marshal(errorJson)
	if err != nil {
		log.Printf("controller.generateErrorJson: Failed to marshal error message to json: %v", err)
		log.Printf("original error message: %v", errorMessage)
		return "{\"error\": \"Could not generate full error message. See console for details.\"}"
	}

	return string(errorJsonBytes)
}
