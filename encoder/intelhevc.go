package encoder

import (
	"errors"
	"ffmpeg-clipper/common"
	"ffmpeg-clipper/ffmpeg"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
)

func ClipIntelHevc(bodyJson map[string]interface{}) (string, error) {
	commonclipParams, err := getCommonClipParameters(bodyJson)
	if err != nil {
		return "", fmt.Errorf("encoder.ClipIntelHevc: could not get common clip params: %w", err)
	}

	encodingPreset, ok := bodyJson["encodingPreset"].(string)
	if !ok {
		return "", errors.New("encoder.ClipIntelHevc: could not determine encoding preset")
	}

	qualityTarget, ok := bodyJson["qualityTarget"].(float64)
	if !ok {
		return "", errors.New("encoder.ClipIntelHevc: could not determine quality target")
	}

	videoExtension := filepath.Ext(commonclipParams.Video)
	videoName := commonclipParams.Video[:len(commonclipParams.Video)-len(videoExtension)]
	newVideoName := fmt.Sprintf("%v_clip%v%v", videoName, common.GetRandomString(), videoExtension)

	cmd := exec.Command(ffmpeg.FfmpegPath,
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
		"hevc_qsv",
		"-preset",
		encodingPreset,
		"-global_quality",
		fmt.Sprintf("%v", qualityTarget),
		"-look_ahead",
		"1",
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
		log.Printf("encoder.ClipIntelHevc: error running ffmpeg: %v\n", cmdOutput)
	}

	return newVideoName, nil
}
