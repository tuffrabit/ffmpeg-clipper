package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/a-h/templ"
	"golang.org/x/exp/slices"
)

const CONFIG_FILENAME = "ffmpeg-clipper-config.json"

var configJson ConfigJson

type ConfigJson struct {
	ClipProfiles []ClipProfileJson `json:"profiles"`
}

type ClipProfileJson struct {
	ProfileName     string                         `json:"profileName"`
	ScaleFactor     float32                        `json:"scaleFactor"`
	Encoder         EncoderType                    `json:"encoder"`
	EncoderSettings ClipProfileJsonEncoderSettings `json:"encoderSettings"`
	Saturation      float32                        `json:"saturation"`
	Contrast        float32                        `json:"contrast"`
	Brightness      float32                        `json:"brightness"`
	Gamma           float32                        `json:"gamma"`
	Exposure        float32                        `json:"exposure"`
	BlackLevel      float32                        `json:"black_level"`
	PlayAfter       bool                           `json:"playAfter"`
}

type ClipProfileJsonEncoderSettings struct {
	Libx264   Libx264EncoderSettings   `json:"libx264"`
	Libx265   Libx265EncoderSettings   `json:"libx265"`
	LibaomAv1 LibaomAv1EncoderSettings `json:"libaom-av1"`
	NvencH264 NvencH264EncoderSettings `json:"h264_nvenc"`
	NvencHevc NvencHevcEncoderSettings `json:"hevc_nvenc"`
	IntelH264 IntelH264EncoderSettings `json:"h264_qsv"`
	IntelHevc IntelHevcEncoderSettings `json:"hevc_qsv"`
	IntelAv1  IntelAv1EncoderSettings  `json:"av1_qsv"`
}

func (s *ClipProfileJsonEncoderSettings) SetEncoderSettings(t EncoderType, i EncoderSettingsInterface) {
	switch t {
	case Libx264EncoderType:
		val, ok := i.(*Libx264EncoderSettings)
		if ok {
			s.Libx264.EncodingPreset = val.EncodingPreset
			s.Libx264.QualityTarget = val.QualityTarget
		}
	case Libx265EncoderType:
		val, ok := i.(*Libx265EncoderSettings)
		if ok {
			s.Libx265.EncodingPreset = val.EncodingPreset
			s.Libx265.QualityTarget = val.QualityTarget
		}
	case LibaomAv1EncoderType:
		val, ok := i.(*LibaomAv1EncoderSettings)
		if ok {
			s.LibaomAv1.QualityTarget = val.QualityTarget
		}
	case NvencH264EncoderType:
		val, ok := i.(*NvencH264EncoderSettings)
		if ok {
			s.NvencH264.EncodingPreset = val.EncodingPreset
			s.NvencH264.QualityTarget = val.QualityTarget
		}
	case NvencHevcEncoderType:
		val, ok := i.(*NvencHevcEncoderSettings)
		if ok {
			s.NvencHevc.EncodingPreset = val.EncodingPreset
			s.NvencHevc.QualityTarget = val.QualityTarget
		}
	case IntelH264EncoderType:
		val, ok := i.(*IntelH264EncoderSettings)
		if ok {
			s.IntelH264.EncodingPreset = val.EncodingPreset
			s.IntelH264.QualityTarget = val.QualityTarget
		}
	case IntelHevcEncoderType:
		val, ok := i.(*IntelHevcEncoderSettings)
		if ok {
			s.IntelHevc.EncodingPreset = val.EncodingPreset
			s.IntelHevc.QualityTarget = val.QualityTarget
		}
	case IntelAv1EncoderType:
		val, ok := i.(*IntelAv1EncoderSettings)
		if ok {
			s.IntelAv1.EncodingPreset = val.EncodingPreset
			s.IntelAv1.QualityTarget = val.QualityTarget
		}
	}
}

type GetEncoderSettingsFunc func() templ.Component

type EncoderSettingsInterface interface {
	Validate() bool
	GetEncodingPreset() string
	GetQualityTarget() int
}

type Libx264EncoderSettings struct {
	EncodingPreset string `json:"encodingPreset"`
	QualityTarget  int    `json:"qualityTarget"`
}

func (s Libx264EncoderSettings) Validate() bool {
	if s.QualityTarget < 0 || s.QualityTarget > 51 {
		return false
	}

	return true
}

func (s Libx264EncoderSettings) GetEncodingPreset() string {
	return s.EncodingPreset
}

func (s Libx264EncoderSettings) GetQualityTarget() int {
	return s.QualityTarget
}

type Libx265EncoderSettings struct {
	EncodingPreset string `json:"encodingPreset"`
	QualityTarget  int    `json:"qualityTarget"`
}

func (s Libx265EncoderSettings) Validate() bool {
	if s.QualityTarget < 0 || s.QualityTarget > 51 {
		return false
	}

	return true
}

func (s Libx265EncoderSettings) GetEncodingPreset() string {
	return s.EncodingPreset
}

func (s Libx265EncoderSettings) GetQualityTarget() int {
	return s.QualityTarget
}

type LibaomAv1EncoderSettings struct {
	QualityTarget int `json:"qualityTarget"`
}

func (s LibaomAv1EncoderSettings) Validate() bool {
	if s.QualityTarget < 0 || s.QualityTarget > 63 {
		return false
	}

	return true
}

func (s LibaomAv1EncoderSettings) GetEncodingPreset() string {
	return ""
}

func (s LibaomAv1EncoderSettings) GetQualityTarget() int {
	return s.QualityTarget
}

type NvencH264EncoderSettings struct {
	EncodingPreset string `json:"encodingPreset"`
	QualityTarget  int    `json:"qualityTarget"`
}

func (s NvencH264EncoderSettings) Validate() bool {
	if s.QualityTarget < 0 || s.QualityTarget > 51 {
		return false
	}

	return true
}

func (s NvencH264EncoderSettings) GetEncodingPreset() string {
	return s.EncodingPreset
}

func (s NvencH264EncoderSettings) GetQualityTarget() int {
	return s.QualityTarget
}

type NvencHevcEncoderSettings struct {
	EncodingPreset string `json:"encodingPreset"`
	QualityTarget  int    `json:"qualityTarget"`
}

func (s NvencHevcEncoderSettings) Validate() bool {
	if s.QualityTarget < 0 || s.QualityTarget > 51 {
		return false
	}

	return true
}

func (s NvencHevcEncoderSettings) GetEncodingPreset() string {
	return s.EncodingPreset
}

func (s NvencHevcEncoderSettings) GetQualityTarget() int {
	return s.QualityTarget
}

type IntelH264EncoderSettings struct {
	EncodingPreset string `json:"encodingPreset"`
	QualityTarget  int    `json:"qualityTarget"`
}

func (s IntelH264EncoderSettings) Validate() bool {
	if s.QualityTarget < 1 || s.QualityTarget > 51 {
		return false
	}

	return true
}

func (s IntelH264EncoderSettings) GetEncodingPreset() string {
	return s.EncodingPreset
}

func (s IntelH264EncoderSettings) GetQualityTarget() int {
	return s.QualityTarget
}

type IntelHevcEncoderSettings struct {
	EncodingPreset string `json:"encodingPreset"`
	QualityTarget  int    `json:"qualityTarget"`
}

func (s IntelHevcEncoderSettings) Validate() bool {
	if s.QualityTarget < 1 || s.QualityTarget > 51 {
		return false
	}

	return true
}

func (s IntelHevcEncoderSettings) GetEncodingPreset() string {
	return s.EncodingPreset
}

func (s IntelHevcEncoderSettings) GetQualityTarget() int {
	return s.QualityTarget
}

type IntelAv1EncoderSettings struct {
	EncodingPreset string `json:"encodingPreset"`
	QualityTarget  int    `json:"qualityTarget"`
}

func (s IntelAv1EncoderSettings) Validate() bool {
	if s.QualityTarget < 1 || s.QualityTarget > 51 {
		return false
	}

	return true
}

func (s IntelAv1EncoderSettings) GetEncodingPreset() string {
	return s.EncodingPreset
}

func (s IntelAv1EncoderSettings) GetQualityTarget() int {
	return s.QualityTarget
}

type EncoderType string

const (
	Libx264EncoderType   EncoderType = "libx264"
	Libx265EncoderType   EncoderType = "libx265"
	LibaomAv1EncoderType EncoderType = "libaom-av1"
	NvencH264EncoderType EncoderType = "h264_nvenc"
	NvencHevcEncoderType EncoderType = "hevc_nvenc"
	IntelH264EncoderType EncoderType = "h264_qsv"
	IntelHevcEncoderType EncoderType = "hevc_qsv"
	IntelAv1EncoderType  EncoderType = "av1_qsv"
)

func NewProfile(name string) *ClipProfileJson {
	libx264EncoderSettings := Libx264EncoderSettings{
		EncodingPreset: "slow",
		QualityTarget:  24,
	}

	libx265EncoderSettings := Libx265EncoderSettings{
		EncodingPreset: "slow",
		QualityTarget:  29,
	}

	libaomAv1EncoderSettings := LibaomAv1EncoderSettings{
		QualityTarget: 29,
	}

	nvencH264EncoderSettings := NvencH264EncoderSettings{
		EncodingPreset: "p4",
		QualityTarget:  26,
	}

	nvencHevcEncoderSettings := NvencHevcEncoderSettings{
		EncodingPreset: "p4",
		QualityTarget:  31,
	}

	intelH264EncoderSettings := IntelH264EncoderSettings{
		EncodingPreset: "medium",
		QualityTarget:  29,
	}

	intelHevcEncoderSettings := IntelHevcEncoderSettings{
		EncodingPreset: "medium",
		QualityTarget:  31,
	}

	intelAv1EncoderSettings := IntelAv1EncoderSettings{
		EncodingPreset: "medium",
		QualityTarget:  33,
	}

	encoderSettings := ClipProfileJsonEncoderSettings{
		Libx264:   libx264EncoderSettings,
		Libx265:   libx265EncoderSettings,
		LibaomAv1: libaomAv1EncoderSettings,
		NvencH264: nvencH264EncoderSettings,
		NvencHevc: nvencHevcEncoderSettings,
		IntelH264: intelH264EncoderSettings,
		IntelHevc: intelHevcEncoderSettings,
		IntelAv1:  intelAv1EncoderSettings,
	}

	return &ClipProfileJson{
		ProfileName:     name,
		ScaleFactor:     2.666,
		Encoder:         Libx264EncoderType,
		EncoderSettings: encoderSettings,
		Saturation:      1,
		Contrast:        1,
		Brightness:      0,
		Gamma:           1,
		Exposure:        0,
		BlackLevel:      0,
		PlayAfter:       true,
	}
}

func GetEncoderTypes() map[EncoderType]string {
	m := make(map[EncoderType]string)

	m[Libx264EncoderType] = "CPU H.264"
	m[Libx265EncoderType] = "CPU H.265"
	m[LibaomAv1EncoderType] = "CPU AV1"
	m[NvencH264EncoderType] = "Nvidia H.264"
	m[NvencHevcEncoderType] = "Nvidia H.265"
	m[IntelH264EncoderType] = "Intel H.264"
	m[IntelHevcEncoderType] = "Intel H.265"
	m[IntelAv1EncoderType] = "Intel AV1"

	return m
}

func ValidateEncoderType(encoderName string) bool {
	for t := range GetEncoderTypes() {
		if encoderName == string(t) {
			return true
		}
	}

	return false
}

func LoadConfig() error {
	_, err := os.Stat(CONFIG_FILENAME)
	if err != nil {
		err = createDefaultConfigFile()
		if err != nil {
			return fmt.Errorf("config.LoadConfig: could not create default config: %w", err)
		}
	}

	_, err = os.Stat(CONFIG_FILENAME)
	if err != nil {
		return errors.New("config.LoadConfig: config file does not exist")
	}

	file, err := os.Open(CONFIG_FILENAME)
	if err != nil {
		return fmt.Errorf("config.LoadConfig: could not open config file: %w", err)
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("config.LoadConfig: could not read config file: %w", err)
	}

	var tempConfigJson ConfigJson
	err = json.Unmarshal(fileBytes, &tempConfigJson)
	if err != nil {
		return fmt.Errorf("config.LoadConfig: could not marshal json: %w", err)
	}

	configJson = tempConfigJson

	return nil
}

func GetConfig() *ConfigJson {
	return &configJson
}

func GetConfigWithLoad() (*ConfigJson, error) {
	err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("config.GetConfigWithLoad: could not load config: %w", err)
	}

	return &configJson, nil
}

func GetProfile(name string) *ClipProfileJson {
	for _, p := range configJson.ClipProfiles {
		if name == p.ProfileName {
			return &p
		}
	}

	return nil
}

func GetEncoderSettingsFromProfile(name string, encoder EncoderType) (ClipProfileJsonEncoderSettings, error) {
	profile := GetProfile(name)
	if profile == nil {
		return ClipProfileJsonEncoderSettings{}, fmt.Errorf("config.GetEncoderSettingsFromProfile: profile %v does not exist", name)
	}

	return profile.EncoderSettings, nil
}

func SaveProfile(profileJson *ClipProfileJson) error {
	configJson, err := GetConfigWithLoad()
	if err != nil {
		return fmt.Errorf("config.SaveProfile: could not get config: %w", err)
	}

	profileIndex := -1
	for index, profile := range configJson.ClipProfiles {
		if profile.ProfileName == profileJson.ProfileName {
			profileIndex = index
			break
		}
	}

	if profileIndex != -1 {
		configJson.ClipProfiles[profileIndex] = *profileJson
	} else {
		configJson.ClipProfiles = append(configJson.ClipProfiles, *profileJson)
	}

	err = writeConfigFile(configJson)
	if err != nil {
		return fmt.Errorf("config.SaveProfile: could not write config: %w", err)
	}

	err = LoadConfig()
	if err != nil {
		return fmt.Errorf("config.SaveProfile: could not load updated config: %w", err)
	}

	return nil
}

func DeleteProfile(profileName string) error {
	configJson, err := GetConfigWithLoad()
	if err != nil {
		return fmt.Errorf("config.DeleteProfile: could not get config: %w", err)
	}

	profileIndex := -1
	for index, profile := range configJson.ClipProfiles {
		if profile.ProfileName == profileName {
			profileIndex = index
			break
		}
	}

	configJson.ClipProfiles = slices.Delete(configJson.ClipProfiles, profileIndex, profileIndex+1)
	err = writeConfigFile(configJson)
	if err != nil {
		return fmt.Errorf("config.DeleteProfile: could not write config: %w", err)
	}

	err = LoadConfig()
	if err != nil {
		return fmt.Errorf("config.DeleteProfile: could not load updated config: %w", err)
	}

	return nil
}

func writeConfigFile(configJson *ConfigJson) error {
	jsonBytes, err := json.MarshalIndent(configJson, "", " ")
	if err != nil {
		return fmt.Errorf("config.writeConfigFile: could not marshal json: %w", err)
	}

	err = os.WriteFile(CONFIG_FILENAME, jsonBytes, 0644)
	if err != nil {
		return fmt.Errorf("config.writeConfigFile: could not write config: %w", err)
	}

	return nil
}

func createDefaultConfigFile() error {
	defaultConfigJson := generateDefaultConfigJson()
	err := writeConfigFile(&defaultConfigJson)
	if err != nil {
		return fmt.Errorf("config.createDefaultConfigFile: could not write file: %w", err)
	}

	return nil
}

func generateDefaultConfigJson() ConfigJson {
	libx264EncoderSettings := Libx264EncoderSettings{
		EncodingPreset: "slow",
		QualityTarget:  24,
	}

	libx265EncoderSettings := Libx265EncoderSettings{
		EncodingPreset: "slow",
		QualityTarget:  29,
	}

	libaomAv1EncoderSettings := LibaomAv1EncoderSettings{
		QualityTarget: 29,
	}

	nvencH264EncoderSettings := NvencH264EncoderSettings{
		EncodingPreset: "p4",
		QualityTarget:  26,
	}

	nvencHevcEncoderSettings := NvencHevcEncoderSettings{
		EncodingPreset: "p4",
		QualityTarget:  31,
	}

	intelH264EncoderSettings := IntelH264EncoderSettings{
		EncodingPreset: "medium",
		QualityTarget:  29,
	}

	intelHevcEncoderSettings := IntelHevcEncoderSettings{
		EncodingPreset: "medium",
		QualityTarget:  31,
	}

	intelAv1EncoderSettings := IntelAv1EncoderSettings{
		EncodingPreset: "medium",
		QualityTarget:  33,
	}

	encoderSettings := ClipProfileJsonEncoderSettings{
		Libx264:   libx264EncoderSettings,
		Libx265:   libx265EncoderSettings,
		LibaomAv1: libaomAv1EncoderSettings,
		NvencH264: nvencH264EncoderSettings,
		NvencHevc: nvencHevcEncoderSettings,
		IntelH264: intelH264EncoderSettings,
		IntelHevc: intelHevcEncoderSettings,
		IntelAv1:  intelAv1EncoderSettings,
	}

	huntDayClipProfile := ClipProfileJson{
		ProfileName:     "Hunt - Day",
		ScaleFactor:     2.666,
		Encoder:         Libx264EncoderType,
		EncoderSettings: encoderSettings,
		Saturation:      2,
		Contrast:        1.1,
		Brightness:      0,
		Gamma:           1,
		Exposure:        0,
		BlackLevel:      0,
		PlayAfter:       true,
	}

	huntNightClipProfile := ClipProfileJson{
		ProfileName:     "Hunt - Night",
		ScaleFactor:     2.666,
		Encoder:         Libx264EncoderType,
		EncoderSettings: encoderSettings,
		Saturation:      2,
		Contrast:        1.1,
		Brightness:      0.1,
		Gamma:           1,
		Exposure:        0,
		BlackLevel:      0,
		PlayAfter:       true,
	}

	destinyClipProfile := ClipProfileJson{
		ProfileName:     "Destiny",
		ScaleFactor:     2.666,
		Encoder:         Libx264EncoderType,
		EncoderSettings: encoderSettings,
		Saturation:      1,
		Contrast:        1,
		Brightness:      0,
		Gamma:           1,
		Exposure:        0,
		BlackLevel:      0,
		PlayAfter:       true,
	}

	return ConfigJson{
		ClipProfiles: []ClipProfileJson{
			huntDayClipProfile,
			huntNightClipProfile,
			destinyClipProfile,
		},
	}
}
