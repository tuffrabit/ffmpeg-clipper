package encoder

import (
	"errors"
	"ffmpeg-clipper/config"
	"fmt"
	"strconv"
)

type CommonClipParameters struct {
	Video           string             `json:"video"`
	StartTime       string             `json:"startTime"`
	EndTime         string             `json:"endTime"`
	ScaleFactor     float64            `json:"scaleFactor"`
	Encoder         config.EncoderType `json:"encoder"`
	Saturation      float64            `json:"saturation"`
	Contrast        float64            `json:"contrast"`
	Brightness      float64            `json:"brightness"`
	Gamma           float64            `json:"gamma"`
	Exposure        float64            `json:"exposure"`
	BlackLevel      float64            `json:"blackLevel"`
	PlayAfter       bool               `json:"playAfter"`
	AlternatePlayer string             `json:"alternatePlayer"`
}

func getCommonClipParameters(bodyJson map[string]interface{}) (*CommonClipParameters, error) {
	video, ok := bodyJson["video"].(string)
	if !ok {
		return nil, errors.New("encoder.getCommonClipParameters: could not determine video")
	}

	startTime, ok := bodyJson["startTime"].(string)
	if !ok {
		return nil, errors.New("encoder.getCommonClipParameters: could not determine start time")
	}

	endTime, ok := bodyJson["endTime"].(string)
	if !ok {
		return nil, errors.New("encoder.getCommonClipParameters: could not determine end time")
	}

	scaleFactor, err := getFloat64FromStringInterface(bodyJson["scaleFactor"])
	if err != nil {
		return nil, fmt.Errorf("encoder.getCommonClipParameters: could not determine scale factor: %w", err)
	}

	saturation, err := getFloat64FromStringInterface(bodyJson["saturation"])
	if err != nil {
		return nil, fmt.Errorf("encoder.getCommonClipParameters: could not determine saturation: %w", err)
	}

	contrast, err := getFloat64FromStringInterface(bodyJson["contrast"])
	if err != nil {
		return nil, fmt.Errorf("encoder.getCommonClipParameters: could not determine contrast: %w", err)
	}

	brightness, err := getFloat64FromStringInterface(bodyJson["brightness"])
	if err != nil {
		return nil, fmt.Errorf("encoder.getCommonClipParameters: could not determine brightness: %w", err)
	}

	gamma, err := getFloat64FromStringInterface(bodyJson["gamma"])
	if err != nil {
		return nil, fmt.Errorf("encoder.getCommonClipParameters: could not determine gamma: %w", err)
	}

	exposure, err := getFloat64FromStringInterface(bodyJson["exposure"])
	if err != nil {
		return nil, fmt.Errorf("encoder.getCommonClipParameters: could not determine exposure: %w", err)
	}

	blackLevel, err := getFloat64FromStringInterface(bodyJson["blackLevel"])
	if err != nil {
		return nil, fmt.Errorf("encoder.getCommonClipParameters: could not determine blackLevel: %w", err)
	}

	return &CommonClipParameters{
		Video:       video,
		StartTime:   startTime,
		EndTime:     endTime,
		ScaleFactor: scaleFactor,
		Saturation:  saturation,
		Contrast:    contrast,
		Brightness:  brightness,
		Gamma:       gamma,
		Exposure:    exposure,
		BlackLevel:  blackLevel,
	}, nil
}

func getFloat64FromStringInterface(value interface{}) (float64, error) {
	stringValue, ok := value.(string)
	if !ok {
		return 0, errors.New("encoder.getFloat64FromStringInterface: string could not be asserted")
	}

	floatValue, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return 0, fmt.Errorf("encoder.getFloat64FromStringInterface: could not parse float: %w", err)
	}

	return floatValue, nil
}
