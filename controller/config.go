package controller

import (
	"encoding/json"
	"ffmpeg-clipper/config"
	"ffmpeg-clipper/html/templ"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func InitConfig() error {
	err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("controller.InitConfig: could not load config: %v", err)
	}

	return nil
}

func GetProfiles(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "text/html")

	component, err := templ.GetProfileNames(config.GetConfig())
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.GetProfiles: could not get profile list %v", err))
		return
	}

	component.Render(r.Context(), w)
}

func GetProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryValues := r.URL.Query()
	profileName := queryValues.Get("name")
	header := w.Header()
	header.Set("Content-Type", "text/html")

	var profile config.ClipProfileJson

	for _, p := range config.GetConfig().ClipProfiles {
		if profileName == p.ProfileName {
			profile = p
			break
		}
	}

	component := templ.GetProfile(profile)
	component.Render(r.Context(), w)
}

func GetEncoderSettings(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	queryValues := r.URL.Query()
	profileName := queryValues.Get("name")
	encoderName := queryValues.Get("encoder")
	header := w.Header()
	header.Set("Content-Type", "text/html")

	if !config.ValidateEncoderType(encoderName) {
		handleResponseError(w, fmt.Sprintf("controller.GetEncoderSettings: encoder %v is not valid", encoderName))
		return
	}

	encoderType := config.EncoderType(encoderName)
	encoderSettings, err := config.GetEncoderSettingsFromProfile(profileName, encoderType)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.GetEncoderSettings: could not get encoder %v from profile %v", encoderName, profileName))
		return
	}

	component := templ.GetEncoderSettings(encoderType, encoderSettings)
	component.Render(r.Context(), w)
}

func SaveProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	profileName := r.FormValue("name")

	if profileName == "" {
		profileName = r.Header.Get("HX-Prompt")
	}

	header := w.Header()
	header.Set("Content-Type", "text/html")

	log.Printf("Profile to save: %v\n", profileName)

	component, err := templ.GetProfileNames(config.GetConfig())
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.SaveProfile: could not get profile list %v", err))
		return
	}

	component.Render(r.Context(), w)

	/*err = config.SaveProfile(&payloadJson)
	if err != nil {
		message := fmt.Sprintf("controller.SaveProfile: could not save profile: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}
	*/
}

func DeleteProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	profileName := r.FormValue("name")
	header := w.Header()
	header.Set("Content-Type", "text/html")

	err := config.DeleteProfile(profileName)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.DeleteProfile: could not delete profile: %v", err))
		return
	}

	err = config.LoadConfig()
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.DeleteProfile: could not load updated config: %v", err))
		return
	}

	component, err := templ.GetProfileNames(config.GetConfig())
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.DeleteProfile: could not get profile list %v", err))
		return
	}

	component.Render(r.Context(), w)
}

func GetConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	resultBytes, err := json.Marshal(config.GetConfig())
	if err != nil {
		message := fmt.Sprintf("controller.GetConfig: could not marshal struct to json: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
	} else {
		fmt.Fprint(w, string(resultBytes))
	}
}

func SaveProfile2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	fmt.Fprint(w, "{}")
}

func DeleteProfile2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	fmt.Fprint(w, "{}")
}
