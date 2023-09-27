package controller

import (
	"encoding/json"
	"ffmpeg-clipper/config"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var configJson *config.ConfigJson

func InitConfig() error {
	var err error
	configJson, err = config.GetConfig()
	if err != nil {
		return fmt.Errorf("controller.InitConfig: could not get config: %v", err)
	}

	return nil
}

func GetConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	resultBytes, err := json.Marshal(configJson)
	if err != nil {
		message := fmt.Sprintf("controller.GetConfig: could not marshal struct to json: %v", err)
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
		message := fmt.Sprintf("controller.SaveProfile: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	payloadJson := config.ClipProfileJson{}
	err = json.Unmarshal(bodyBytes, &payloadJson)
	if err != nil {
		message := fmt.Sprintf("controller.SaveProfile: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	err = config.SaveProfile(&payloadJson)
	if err != nil {
		message := fmt.Sprintf("controller.SaveProfile: could not save profile: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	configJson, err = config.GetConfig()
	if err != nil {
		message := fmt.Sprintf("controller.SaveProfile: could not load updated config: %v", err)
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
		message := fmt.Sprintf("controller.DeleteProfile: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	payloadJson := struct {
		ProfileName string `json:"profileName"`
	}{}
	err = json.Unmarshal(bodyBytes, &payloadJson)
	if err != nil {
		message := fmt.Sprintf("controller.DeleteProfile: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	err = config.DeleteProfile(payloadJson.ProfileName)
	if err != nil {
		message := fmt.Sprintf("controller.DeleteProfile: could not delete profile: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	configJson, err = config.GetConfig()
	if err != nil {
		message := fmt.Sprintf("controller.DeleteProfile: could not load updated config: %v", err)
		fmt.Fprint(w, generateErrorJson(message))
		log.Fatalln(message)
	}

	fmt.Fprint(w, "{}")
}
