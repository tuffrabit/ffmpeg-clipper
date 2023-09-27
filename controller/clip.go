package controller

import (
	"encoding/json"
	"errors"
	"ffmpeg-clipper/common"
	"ffmpeg-clipper/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
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

func ClipVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	header := w.Header()
	header.Set("Content-Type", "application/json")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: could not read request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	var bodyJson map[string]interface{}

	err = json.Unmarshal(bodyBytes, &bodyJson)
	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: could not json marshal request body: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	encoder, ok := bodyJson["encoder"].(string)
	if !ok {
		message := "controller.ClipVideo: could not determine encoder"
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	var newVideoName string

	switch config.EncoderType(encoder) {
	case config.Libx264EncoderType:
		newVideoName, err = clipLibx264(bodyJson)
	case config.Libx265EncoderType:
		newVideoName, err = clipLibx265(bodyJson)
	case config.NvencH264EncoderType:
		newVideoName, err = clipNvencH264(bodyJson)
	case config.NvencHevcEncoderType:
		newVideoName, err = clipNvencHevc(bodyJson)
	case config.LibaomAv1EncoderType:
		newVideoName, err = clipLibaomAv1(bodyJson)
	default:
		err = fmt.Errorf("%v is not a valid encoder", encoder)
	}

	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: could not create clip: %v", err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	_, err = os.Stat(newVideoName)
	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: new video clip %v was not created: %v", newVideoName, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		return
	}

	ffprobeErr := false

	cmd := exec.Command(
		FfprobePath,
		"-v",
		"error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=width,height",
		"-of",
		"csv=s=x:p=0",
		newVideoName,
	)
	output, err := common.RunSystemCommand(cmd)
	if err != nil {
		message := fmt.Sprintf("controller.ClipVideo: ffprobe failed\nstderr: %v\nerr: %v", output, err)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		ffprobeErr = true
	}

	output = strings.TrimSuffix(output, "\r\n")
	output = strings.TrimSuffix(output, "\n")
	output = strings.TrimSuffix(output, "\r")

	if output == "" {
		message := fmt.Sprintf("controller.ClipVideo: generated clip %v is invalid", newVideoName)
		log.Println(message)
		fmt.Fprint(w, generateErrorJson(message))
		ffprobeErr = true
	}

	if ffprobeErr {
		_, err = os.Stat(newVideoName)
		if err == nil {
			os.Remove(newVideoName)
		}
		return
	}

	fmt.Fprint(w, "{\"newVideoName\": \""+newVideoName+"\"}")
}

func getCommonClipParameters(bodyJson map[string]interface{}) (*CommonClipParameters, error) {
	video, ok := bodyJson["video"].(string)
	if !ok {
		return nil, errors.New("controller.getCommonClipParameters: could not determine video")
	}

	startTime, ok := bodyJson["startTime"].(string)
	if !ok {
		return nil, errors.New("controller.getCommonClipParameters: could not determine start time")
	}

	endTime, ok := bodyJson["endTime"].(string)
	if !ok {
		return nil, errors.New("controller.getCommonClipParameters: could not determine end time")
	}

	scaleFactor, ok := bodyJson["scaleFactor"].(float64)
	if !ok {
		return nil, errors.New("controller.getCommonClipParameters: could not determine scale factor")
	}

	saturation, ok := bodyJson["saturation"].(float64)
	if !ok {
		return nil, errors.New("controller.getCommonClipParameters: could not determine saturation")
	}

	contrast, ok := bodyJson["contrast"].(float64)
	if !ok {
		return nil, errors.New("controller.getCommonClipParameters: could not determine contrast")
	}

	brightness, ok := bodyJson["brightness"].(float64)
	if !ok {
		return nil, errors.New("controller.getCommonClipParameters: could not determine brightness")
	}

	gamma, ok := bodyJson["gamma"].(float64)
	if !ok {
		return nil, errors.New("controller.getCommonClipParameters: could not determine gamma")
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

func clipLibx264(bodyJson map[string]interface{}) (string, error) {
	commonclipParams, err := getCommonClipParameters(bodyJson)
	if err != nil {
		return "", fmt.Errorf("controller.clipLibx264: could not get common clip params: %w", err)
	}

	encodingPreset, ok := bodyJson["encodingPreset"].(string)
	if !ok {
		return "", errors.New("controller.clipLibx264: could not determine encoding preset")
	}

	qualityTarget, ok := bodyJson["qualityTarget"].(float64)
	if !ok {
		return "", errors.New("controller.clipLibx264: could not determine quality target")
	}

	videoExtension := filepath.Ext(commonclipParams.Video)
	videoName := commonclipParams.Video[:len(commonclipParams.Video)-len(videoExtension)]
	newVideoName := fmt.Sprintf("%v_clip%v%v", videoName, common.GetRandomString(), videoExtension)

	cmd := exec.Command(FfmpegPath,
		"-nostats",
		"-hide_banner",
		"-loglevel",
		"error",
		"-ss",
		commonclipParams.StartTime,
		"-to",
		commonclipParams.EndTime,
		"-i",
		commonclipParams.Video,
		"-c:v",
		"libx264",
		"-preset",
		encodingPreset,
		"-crf",
		fmt.Sprintf("%v", qualityTarget),
		"-vf",
		fmt.Sprintf("scale=iw/%v:-1:flags=bicubic,eq=saturation=%v:contrast=%v:brightness=%v:gamma=%v",
			commonclipParams.ScaleFactor,
			commonclipParams.Saturation,
			commonclipParams.Contrast,
			commonclipParams.Brightness,
			commonclipParams.Gamma,
		),
		newVideoName,
	)
	cmdOutput, err := common.RunSystemCommand(cmd)
	if err != nil {
		log.Printf("controller.clipLibx264: error running ffmpeg: %v\n", cmdOutput)
	}

	return newVideoName, nil
}

func clipLibx265(bodyJson map[string]interface{}) (string, error) {
	commonclipParams, err := getCommonClipParameters(bodyJson)
	if err != nil {
		return "", fmt.Errorf("controller.clipLibx265: could not get common clip params: %w", err)
	}

	encodingPreset, ok := bodyJson["encodingPreset"].(string)
	if !ok {
		return "", errors.New("controller.clipLibx265: could not determine encoding preset")
	}

	qualityTarget, ok := bodyJson["qualityTarget"].(float64)
	if !ok {
		return "", errors.New("controller.clipLibx265: could not determine quality target")
	}

	videoExtension := filepath.Ext(commonclipParams.Video)
	videoName := commonclipParams.Video[:len(commonclipParams.Video)-len(videoExtension)]
	newVideoName := fmt.Sprintf("%v_clip%v%v", videoName, common.GetRandomString(), videoExtension)

	cmd := exec.Command(FfmpegPath,
		"-nostats",
		"-hide_banner",
		"-loglevel",
		"error",
		"-ss",
		commonclipParams.StartTime,
		"-to",
		commonclipParams.EndTime,
		"-i",
		commonclipParams.Video,
		"-c:v",
		"libx265",
		"-preset",
		encodingPreset,
		"-crf",
		fmt.Sprintf("%v", qualityTarget),
		"-vf",
		fmt.Sprintf("scale=iw/%v:-1:flags=bicubic,eq=saturation=%v:contrast=%v:brightness=%v:gamma=%v",
			commonclipParams.ScaleFactor,
			commonclipParams.Saturation,
			commonclipParams.Contrast,
			commonclipParams.Brightness,
			commonclipParams.Gamma,
		),
		newVideoName,
	)
	cmdOutput, err := common.RunSystemCommand(cmd)
	if err != nil {
		log.Printf("controller.clipLibx265: error running ffmpeg: %v\n", cmdOutput)
	}

	return newVideoName, nil
}

func clipLibaomAv1(bodyJson map[string]interface{}) (string, error) {
	commonclipParams, err := getCommonClipParameters(bodyJson)
	if err != nil {
		return "", fmt.Errorf("controller.clipLibaomAv1: could not get common clip params: %w", err)
	}

	qualityTarget, ok := bodyJson["qualityTarget"].(float64)
	if !ok {
		return "", errors.New("controller.clipLibaomAv1: could not determine quality target")
	}

	videoExtension := filepath.Ext(commonclipParams.Video)
	videoName := commonclipParams.Video[:len(commonclipParams.Video)-len(videoExtension)]
	newVideoName := fmt.Sprintf("%v_clip%v%v", videoName, common.GetRandomString(), videoExtension)

	cmd := exec.Command(FfmpegPath,
		"-nostats",
		"-hide_banner",
		"-loglevel",
		"error",
		"-ss",
		commonclipParams.StartTime,
		"-to",
		commonclipParams.EndTime,
		"-i",
		commonclipParams.Video,
		"-c:v",
		"libaom-av1",
		"-crf",
		fmt.Sprintf("%v", qualityTarget),
		"-vf",
		fmt.Sprintf("scale=iw/%v:-1:flags=bicubic,eq=saturation=%v:contrast=%v:brightness=%v:gamma=%v",
			commonclipParams.ScaleFactor,
			commonclipParams.Saturation,
			commonclipParams.Contrast,
			commonclipParams.Brightness,
			commonclipParams.Gamma,
		),
		newVideoName,
	)
	cmdOutput, err := common.RunSystemCommand(cmd)
	if err != nil {
		log.Printf("controller.clipLibaomAv1: error running ffmpeg: %v\n", cmdOutput)
	}

	return newVideoName, nil
}

func clipNvencH264(bodyJson map[string]interface{}) (string, error) {
	commonclipParams, err := getCommonClipParameters(bodyJson)
	if err != nil {
		return "", fmt.Errorf("controller.clipNvencH264: could not get common clip params: %w", err)
	}

	encodingPreset, ok := bodyJson["encodingPreset"].(string)
	if !ok {
		return "", errors.New("controller.clipNvencH264: could not determine encoding preset")
	}

	qualityTarget, ok := bodyJson["qualityTarget"].(float64)
	if !ok {
		return "", errors.New("controller.clipNvencH264: could not determine quality target")
	}

	videoExtension := filepath.Ext(commonclipParams.Video)
	videoName := commonclipParams.Video[:len(commonclipParams.Video)-len(videoExtension)]
	newVideoName := fmt.Sprintf("%v_clip%v%v", videoName, common.GetRandomString(), videoExtension)

	cmd := exec.Command(FfmpegPath,
		"-nostats",
		"-hide_banner",
		"-loglevel",
		"error",
		"-ss",
		commonclipParams.StartTime,
		"-to",
		commonclipParams.EndTime,
		"-i",
		commonclipParams.Video,
		"-c:v",
		"h264_nvenc",
		"-rc",
		"constqp",
		"-preset",
		encodingPreset,
		"-qp",
		fmt.Sprintf("%v", qualityTarget),
		"-vf",
		fmt.Sprintf("scale=iw/%v:-1:flags=bicubic,eq=saturation=%v:contrast=%v:brightness=%v:gamma=%v",
			commonclipParams.ScaleFactor,
			commonclipParams.Saturation,
			commonclipParams.Contrast,
			commonclipParams.Brightness,
			commonclipParams.Gamma,
		),
		newVideoName,
	)
	cmdOutput, err := common.RunSystemCommand(cmd)
	if err != nil {
		log.Printf("controller.clipNvencH264: error running ffmpeg: %v\n", cmdOutput)
	}

	return newVideoName, nil
}

func clipNvencHevc(bodyJson map[string]interface{}) (string, error) {
	commonclipParams, err := getCommonClipParameters(bodyJson)
	if err != nil {
		return "", fmt.Errorf("controller.clipNvencHevc: could not get common clip params: %w", err)
	}

	encodingPreset, ok := bodyJson["encodingPreset"].(string)
	if !ok {
		return "", errors.New("controller.clipNvencHevc: could not determine encoding preset")
	}

	qualityTarget, ok := bodyJson["qualityTarget"].(float64)
	if !ok {
		return "", errors.New("controller.clipNvencHevc: could not determine quality target")
	}

	videoExtension := filepath.Ext(commonclipParams.Video)
	videoName := commonclipParams.Video[:len(commonclipParams.Video)-len(videoExtension)]
	newVideoName := fmt.Sprintf("%v_clip%v%v", videoName, common.GetRandomString(), videoExtension)

	cmd := exec.Command(FfmpegPath,
		"-nostats",
		"-hide_banner",
		"-loglevel",
		"error",
		"-ss",
		commonclipParams.StartTime,
		"-to",
		commonclipParams.EndTime,
		"-i",
		commonclipParams.Video,
		"-c:v",
		"hevc_nvenc",
		"-rc",
		"constqp",
		"-preset",
		encodingPreset,
		"-qp",
		fmt.Sprintf("%v", qualityTarget),
		"-vf",
		fmt.Sprintf("scale=iw/%v:-1:flags=bicubic,eq=saturation=%v:contrast=%v:brightness=%v:gamma=%v",
			commonclipParams.ScaleFactor,
			commonclipParams.Saturation,
			commonclipParams.Contrast,
			commonclipParams.Brightness,
			commonclipParams.Gamma,
		),
		newVideoName,
	)
	cmdOutput, err := common.RunSystemCommand(cmd)
	if err != nil {
		log.Printf("controller.clipNvencHevc: error running ffmpeg: %v\n", cmdOutput)
	}

	return newVideoName, nil
}
