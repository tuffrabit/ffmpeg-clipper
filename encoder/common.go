package encoder

import (
	"errors"
	"ffmpeg-clipper/config"
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

	scaleFactor, ok := bodyJson["scaleFactor"].(float64)
	if !ok {
		return nil, errors.New("encoder.getCommonClipParameters: could not determine scale factor")
	}

	saturation, ok := bodyJson["saturation"].(float64)
	if !ok {
		return nil, errors.New("encoder.getCommonClipParameters: could not determine saturation")
	}

	contrast, ok := bodyJson["contrast"].(float64)
	if !ok {
		return nil, errors.New("encoder.getCommonClipParameters: could not determine contrast")
	}

	brightness, ok := bodyJson["brightness"].(float64)
	if !ok {
		return nil, errors.New("encoder.getCommonClipParameters: could not determine brightness")
	}

	gamma, ok := bodyJson["gamma"].(float64)
	if !ok {
		return nil, errors.New("encoder.getCommonClipParameters: could not determine gamma")
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
	}, nil
}
