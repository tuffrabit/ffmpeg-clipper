package templ

import (
	templates "ffmpeg-clipper/templ"
	video "ffmpeg-clipper/video"
	"fmt"
	"os"

	"github.com/a-h/templ"
)

func GetAvailableVideos() (templ.Component, error) {
	videos, err := video.GetAvailableVideos()
	if err != nil {
		return nil, fmt.Errorf("html/templ/video.GetAvailableVideos: could not get video list: %w", err)
	}

	return templates.GetAvailableVideos(videos), nil
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

func GetVideoPlayerOutOfBand() templ.Component {
	return templates.GetVideoPlayerOutOfBand()
}

func GetVideoResolutionOutOfBand(resolution string) templ.Component {
	return templates.GetVideoResolutionOutOfBand(resolution)
}
