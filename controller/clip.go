package controller

import (
	"ffmpeg-clipper/config"
	ffmpegEncoder "ffmpeg-clipper/encoder"
	"ffmpeg-clipper/html/templ"
	"ffmpeg-clipper/video"
	"fmt"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func ClipVideo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	encoder := r.FormValue("encoder")
	header := w.Header()
	header.Set("Content-Type", "text/html")

	if encoder == "" {
		handleResponseError(w, "controller.ClipVideo: encoder is a required field")
		return
	}

	fields := make(map[string]interface{})
	fields["video"] = r.FormValue("video")
	fields["startTime"] = r.FormValue("start-time")
	fields["endTime"] = r.FormValue("end-time")
	fields["scaleFactor"] = r.FormValue("scale-factor")
	fields["saturation"] = r.FormValue("saturation")
	fields["contrast"] = r.FormValue("contrast")
	fields["brightness"] = r.FormValue("brightness")
	fields["gamma"] = r.FormValue("gamma")
	fields["exposure"] = r.FormValue("exposure")
	fields["blackLevel"] = r.FormValue("black-level")
	fields["encoder"] = r.FormValue("encoder")
	fields["encodingPreset"] = r.FormValue("encoding-preset")
	fields["qualityTarget"] = r.FormValue("quality-target")

	var newVideoName string
	var err error

	switch config.EncoderType(encoder) {
	case config.Libx264EncoderType:
		newVideoName, err = ffmpegEncoder.ClipLibx264(fields)
	case config.Libx265EncoderType:
		newVideoName, err = ffmpegEncoder.ClipLibx265(fields)
	case config.LibaomAv1EncoderType:
		newVideoName, err = ffmpegEncoder.ClipLibaomAv1(fields)
	case config.NvencH264EncoderType:
		newVideoName, err = ffmpegEncoder.ClipNvencH264(fields)
	case config.NvencHevcEncoderType:
		newVideoName, err = ffmpegEncoder.ClipNvencHevc(fields)
	case config.IntelH264EncoderType:
		newVideoName, err = ffmpegEncoder.ClipIntelH264(fields)
	case config.IntelHevcEncoderType:
		newVideoName, err = ffmpegEncoder.ClipIntelHevc(fields)
	case config.IntelAv1EncoderType:
		newVideoName, err = ffmpegEncoder.ClipIntelAv1(fields)
	default:
		err = fmt.Errorf("%v is not a valid encoder", encoder)
	}

	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.ClipVideo: could not create clip: %v", err))
		return
	}

	_, err = os.Stat(newVideoName)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.ClipVideo: new video clip %v was not created: %v", newVideoName, err))
		return
	}

	ffprobeErr := false

	/*cmd := exec.Command(
		ffmpeg.FfprobePath,
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
		handleResponseError(w, fmt.Sprintf("controller.ClipVideo: ffprobe failed\nstderr: %v\nerr: %v", output, err))
		return
	}

	output = strings.TrimSuffix(output, "\r\n")
	output = strings.TrimSuffix(output, "\n")
	output = strings.TrimSuffix(output, "\r")

	if output == "" {
		handleResponseError(w, fmt.Sprintf("controller.ClipVideo: generated clip %v is invalid", newVideoName))
		ffprobeErr = true
	}*/

	resolution, err := video.GetVideoResolution(newVideoName)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.ClipVideo: generated clip %v is invalid", newVideoName))
		ffprobeErr = true
	}

	if ffprobeErr {
		_, err = os.Stat(newVideoName)
		if err == nil {
			os.Remove(newVideoName)
		}
		return
	}

	availableVideosSelectComponent, err := templ.GetAvailableVideosSelect(true)
	if err != nil {
		handleResponseError(w, fmt.Sprintf("controller.ClipVideo: could not get video list %v", err))
		return
	}

	availableVideosSelectComponent.Render(r.Context(), w)

	playAfter := r.FormValue("play-after")
	if playAfter != "" {
		videoPlayerComponent, err := templ.GetVideoPlayerOutOfBand(newVideoName)
		if err != nil {
			handleResponseError(w, fmt.Sprintf("controller.ClipVideo: could not get video player %v", err))
			return
		}

		videoDetailsComponent, err := templ.GetVideoDetailsOutOfBand(resolution)
		if err != nil {
			handleResponseError(w, fmt.Sprintf("controller.GetVideoPlayer: could not get video details: %v", err))
			return
		}

		videoPlayerComponent.Render(r.Context(), w)
		videoDetailsComponent.Render(r.Context(), w)
		videoNameComponent := templ.GetVideoNameOutOfBand(newVideoName)
		videoNameComponent.Render(r.Context(), w)
	}
}
