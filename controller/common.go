package controller

import (
	"encoding/json"
	"log"
)

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
