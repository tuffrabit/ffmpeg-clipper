package controller

import (
	"errors"
	"ffmpeg-clipper/config"
	"ffmpeg-clipper/html/templ"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

	component := templ.GetProfileNames("", config.GetConfig())
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
	r.ParseForm()
	profileName := r.FormValue("name")

	if profileName == "" {
		profileName = r.Header.Get("HX-Prompt")
	}

	header := w.Header()
	header.Set("Content-Type", "text/html")

	profile, err := getClipProfileJsonFromRequest(profileName, r.Form)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.SaveProfile: could not convert request payload to valid profile %v", err))
		return
	}

	err = config.SaveProfile(profile)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.SaveProfile: could not save profile %v", err))
		return
	}

	component := templ.GetProfileAndList(profileName, config.GetConfig(), *profile)
	component.Render(r.Context(), w)
}

func getClipProfileJsonFromRequest(name string, values url.Values) (*config.ClipProfileJson, error) {
	profileJson := config.GetProfile(name)
	if profileJson == nil {
		profileJson = config.NewProfile(name)
	}

	scaleFactor, err := getFloatValueFromRequest("scale-factor", values)
	if err != nil {
		return nil, fmt.Errorf("controller.getConfigProfileJsonFromRequest: scale-factor is invalid: %w", err)
	}

	if !values.Has("encoder") {
		return nil, errors.New("controller.getConfigProfileJsonFromRequest: encoder value is required")
	}

	encoderSettings, err := getEncoderSettingsFromRequest(values)
	if err != nil {
		return nil, fmt.Errorf("controller.getConfigProfileJsonFromRequest: could not resolve valid encoder settings %w", err)
	}

	saturation, err := getFloatValueFromRequest("saturation", values)
	if err != nil {
		return nil, fmt.Errorf("controller.getConfigProfileJsonFromRequest: saturation is invalid: %w", err)
	}

	contrast, err := getFloatValueFromRequest("contrast", values)
	if err != nil {
		return nil, fmt.Errorf("controller.getConfigProfileJsonFromRequest: contrast is invalid: %w", err)
	}

	brightness, err := getFloatValueFromRequest("brightness", values)
	if err != nil {
		return nil, fmt.Errorf("controller.getConfigProfileJsonFromRequest: brightness is invalid: %w", err)
	}

	gamma, err := getFloatValueFromRequest("gamma", values)
	if err != nil {
		return nil, fmt.Errorf("controller.getConfigProfileJsonFromRequest: gamma is invalid: %w", err)
	}

	playAfter := false
	if values.Has("play-after") {
		playAfter = true
	}

	profileJson.ScaleFactor = scaleFactor
	profileJson.Encoder = config.EncoderType(values.Get("encoder"))
	profileJson.EncoderSettings.SetEncoderSettings(profileJson.Encoder, *encoderSettings)
	profileJson.Saturation = saturation
	profileJson.Contrast = contrast
	profileJson.Brightness = brightness
	profileJson.Gamma = gamma
	profileJson.PlayAfter = playAfter

	return profileJson, nil
}

func getEncoderSettingsFromRequest(values url.Values) (*config.EncoderSettingsInterface, error) {
	if !values.Has("quality-target") {
		return nil, errors.New("controller.getEncoderSettingsFromRequest: quality-target value is required")
	}

	qualityTargetString := values.Get("quality-target")
	qualityTarget, err := strconv.ParseInt(qualityTargetString, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("controller.getEncoderSettingsFromRequest: %v is an invalid value for quality-target: %w", qualityTargetString, err)
	}

	encoderPreset := ""
	if values.Has("encoding-preset") {
		encoderPreset = values.Get("encoding-preset")
	}

	encoderType := config.EncoderType(values.Get("encoder"))

	switch encoderType {
	case config.Libx264EncoderType,
		config.Libx265EncoderType,
		config.NvencH264EncoderType,
		config.NvencHevcEncoderType,
		config.IntelH264EncoderType,
		config.IntelHevcEncoderType,
		config.IntelAv1EncoderType:
		if encoderPreset == "" {
			return nil, errors.New("controller.getEncoderSettingsFromRequest: encoding-preset value is required")
		}
	}

	var encoderSettings config.EncoderSettingsInterface

	switch encoderType {
	case config.Libx264EncoderType:
		encoderSettings = &config.Libx264EncoderSettings{
			EncodingPreset: encoderPreset,
			QualityTarget:  int(qualityTarget),
		}
	case config.Libx265EncoderType:
		encoderSettings = &config.Libx265EncoderSettings{
			EncodingPreset: encoderPreset,
			QualityTarget:  int(qualityTarget),
		}
	case config.LibaomAv1EncoderType:
		encoderSettings = &config.LibaomAv1EncoderSettings{
			QualityTarget: int(qualityTarget),
		}
	case config.NvencH264EncoderType:
		encoderSettings = &config.NvencH264EncoderSettings{
			EncodingPreset: encoderPreset,
			QualityTarget:  int(qualityTarget),
		}
	case config.NvencHevcEncoderType:
		encoderSettings = &config.NvencHevcEncoderSettings{
			EncodingPreset: encoderPreset,
			QualityTarget:  int(qualityTarget),
		}
	case config.IntelH264EncoderType:
		encoderSettings = &config.IntelH264EncoderSettings{
			EncodingPreset: encoderPreset,
			QualityTarget:  int(qualityTarget),
		}
	case config.IntelHevcEncoderType:
		encoderSettings = &config.IntelHevcEncoderSettings{
			EncodingPreset: encoderPreset,
			QualityTarget:  int(qualityTarget),
		}
	case config.IntelAv1EncoderType:
		encoderSettings = &config.IntelAv1EncoderSettings{
			EncodingPreset: encoderPreset,
			QualityTarget:  int(qualityTarget),
		}
	}

	if !encoderSettings.Validate() {
		return nil, errors.New("controller.getEncoderSettingsFromRequest: values are not valid")
	}

	return &encoderSettings, nil
}

func getFloatValueFromRequest(name string, values url.Values) (float32, error) {
	if !values.Has(name) {
		return 0, fmt.Errorf("controller.getFloatValueFromRequest: %v value is required", name)
	}

	stringValue := values.Get(name)
	value, err := strconv.ParseFloat(stringValue, 32)
	if err != nil {
		return 0, fmt.Errorf("controller.getFloatValueFromRequest: %v is an invalid value for %v", stringValue, name)
	}

	return float32(value), nil
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

	component := templ.GetProfileNames("", config.GetConfig())
	component.Render(r.Context(), w)
}
