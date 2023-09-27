# ffmpeg-clipper

## Description
Tool for making video clips. Download a executable from the releases and drop it in a folder with a bunch of video files. Run it from that folder and it will open your default web browser pointing to a new local web server on a random port. This tool takes all the complicated parts of using FFmpeg to edit simple video clips and hides it behind a (hopefully) easy to use web interface. You should have FFmpeg, FFplay, and FFprobe either in the same directory or installed globally on your system.

## TODO
- return new filename - done
- scroll on list box - done
- error/message close button - done
- sane defaults - done
- add support for ffmpeg in path - done
- add brightness - done
- field validation - done
- use other media player - done
- preferences/profiles - done
- add source width/height validation - done
- revisit runSystemCommand stdout/stderr - done
- styling - done-ish
- work indicator - done
- use ffprobe - done
- add default eq filter value text - done
- don't freeze frontend on video play - done
- add delete video button - done
- add help - done
- clean up failed clips - done
- create actual readme - done-ish
- add linux / macos support
- sub second start and stop support
- marshal json responses
- big code cleanup/refactor - done
- support for other encoders
 - hevc/libx265 - done
 - h264_amf
 - h264_nvenc - done
 - h264_qsv
 - hevc_amf
 - hevc_nvenc - done
 - hevc_qsv
 - av1/libaom-av1 - done
 - av1_nvenc
 - av1_qsv
 - av1_amf