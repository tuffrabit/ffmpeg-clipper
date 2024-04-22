package templ

import (
	"errors"
	templates "ffmpeg-clipper/templ"
	video "ffmpeg-clipper/video"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

func GetAvailableVideosSelect(oob bool) (templ.Component, error) {
	videos, err := video.GetAvailableVideos()
	if err != nil {
		return nil, fmt.Errorf("html/templ/video.GetAvailableVideosSelect: could not get video list: %w", err)
	}

	return templates.GetAvailableVideosSelect(videos, oob), nil
}

func GetVideoPlayer(path string) (templ.Component, error) {
	if path == "" {
		return templates.GetVideoPlayer(""), nil
	}

	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("html/templ/video.GetVideoPlayer: %v does not exist: %w", path, err)
	}

	return templates.GetVideoPlayer("streamvideo/" + path), nil
}

func GetVideoPlayerOutOfBand(path string) (templ.Component, error) {
	if path == "" {
		return templates.GetVideoPlayerOutOfBand(""), nil
	}

	_, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("html/templ/video.GetVideoPlayerOutOfBand: %v does not exist: %w", path, err)
	}

	return templates.GetVideoPlayerOutOfBand("streamvideo/" + path), nil
}

func GetVideoDetailsOutOfBand(resolution string) (templ.Component, error) {
	if resolution == "" {
		return templates.GetVideoDetailsOutOfBand("", ""), nil
	}

	if !strings.Contains(resolution, "x") {
		return nil, errors.New("html/templ/video.GetVideoDetailsOutOfBand: resolution string not valid")
	}

	parts := strings.Split(resolution, "x")
	if len(parts) < 2 {
		return nil, errors.New("html/templ/video.GetVideoDetailsOutOfBand: resolution string does not contain enough parts")
	}

	_, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("html/templ/video.GetVideoDetailsOutOfBand: could not parse width: %w", err)
	}

	_, err = strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("html/templ/video.GetVideoDetailsOutOfBand: could not parse height: %w", err)
	}

	return templates.GetVideoDetailsOutOfBand(parts[0], parts[1]), nil
}

func GetVideoResolutionOutOfBand(resolution string) templ.Component {
	return templates.GetVideoResolutionOutOfBand(resolution)
}

func GetVideoNameOutOfBand(name string) templ.Component {
	return templates.GetVideoNameOutOfBand(name)
}
