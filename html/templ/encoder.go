package templ

import (
	"ffmpeg-clipper/config"
	templates "ffmpeg-clipper/templ"

	"github.com/a-h/templ"
)

func GetEncoderSettings(encoderType config.EncoderType, settings config.ClipProfileJsonEncoderSettings) templ.Component {
	return templates.GetEncoderSettings(encoderType, settings)
}
