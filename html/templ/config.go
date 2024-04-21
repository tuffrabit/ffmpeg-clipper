package templ

import (
	"ffmpeg-clipper/config"
	templates "ffmpeg-clipper/templ"

	"github.com/a-h/templ"
)

func GetProfileAndList(selectedProfile string, configJson *config.ConfigJson, profile config.ClipProfileJson) templ.Component {
	var profiles []string

	for _, profile := range configJson.ClipProfiles {
		profiles = append(profiles, profile.ProfileName)
	}

	return templates.GetProfileAndList(selectedProfile, profiles, profile)
}

func GetProfileNames(selectedProfile string, configJson *config.ConfigJson) templ.Component {
	var profiles []string

	for _, profile := range configJson.ClipProfiles {
		profiles = append(profiles, profile.ProfileName)
	}

	return templates.GetProfileNames(selectedProfile, profiles)
}

func GetProfile(profile config.ClipProfileJson) templ.Component {
	return templates.GetProfile(profile)
}
