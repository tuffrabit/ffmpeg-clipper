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

func ClipNvencH264(bodyJson map[string]interface{}) (string, error) {
	commonclipParams, err := getCommonClipParameters(bodyJson)
	if err != nil {
		return "", fmt.Errorf("encoder.ClipNvencH264: could not get common clip params: %w", err)
	}

	encodingPreset, ok := bodyJson["encodingPreset"].(string)
	if !ok {
		return "", errors.New("encoder.ClipNvencH264: could not determine encoding preset")
	}

	qualityTarget, err := getFloat64FromStringInterface(bodyJson["qualityTarget"])
	if err != nil {
		return "", fmt.Errorf("encoder.ClipNvencH264: could not determine qualityTarget: %w", err)
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
		"h264_nvenc",
		"-rc",
		"constqp",
		"-preset",
		encodingPreset,
		"-qp",
		fmt.Sprintf("%v", qualityTarget),
		"-vf",
		fmt.Sprintf("scale=iw/%v:-1:flags=bicubic,exposure=%v:black=%v,eq=saturation=%v:contrast=%v:brightness=%v:gamma=%v",
			commonclipParams.ScaleFactor,
			commonclipParams.Exposure,
			commonclipParams.BlackLevel,
			commonclipParams.Saturation,
			commonclipParams.Contrast,
			commonclipParams.Brightness,
			commonclipParams.Gamma,
		),
		newVideoName,
	)
	cmdOutput, err := common.RunSystemCommand(cmd)
	if err != nil {
		log.Printf("encoder.ClipNvencH264: error running ffmpeg: %v\n", cmdOutput)
	}

	return newVideoName, nil
}
