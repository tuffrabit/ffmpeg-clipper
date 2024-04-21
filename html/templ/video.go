package templ

import (
	templates "ffmpeg-clipper/templ"
	video "ffmpeg-clipper/video"
	"fmt"
	"os"

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

func GetVideoResolutionOutOfBand(resolution string) templ.Component {
	return templates.GetVideoResolutionOutOfBand(resolution)
}

func GetVideoNameOutOfBand(name string) templ.Component {
	return templates.GetVideoNameOutOfBand(name)
}
