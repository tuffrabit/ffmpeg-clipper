package templ

import (
	"ffmpeg-clipper/config"
	templates "ffmpeg-clipper/templ"

	"github.com/a-h/templ"
)

func GetProfileNames(configJson *config.ConfigJson) (templ.Component, error) {
	var profiles []string

	for _, profile := range configJson.ClipProfiles {
		profiles = append(profiles, profile.ProfileName)
	}

	return templates.GetProfileNames(profiles), nil
}

func GetProfile(profile config.ClipProfileJson) templ.Component {
	return templates.GetProfile(profile)
}
