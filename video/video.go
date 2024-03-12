package video

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetAvailableVideos() ([]string, error) {
	var availableVideos []string
	dirEntries, err := os.ReadDir("./")
	if err != nil {
		return nil, fmt.Errorf("controller.GetAvailableVideos: could not list directory: %w", err)
	}

	allowedExtensions := map[string]struct{}{
		".mp4":  {},
		".mkv":  {},
		".avi":  {},
		".flv":  {},
		".mov":  {},
		".wmv":  {},
		".ogg":  {},
		".webm": {},
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			fileExtension := filepath.Ext(entry.Name())

			_, ok := allowedExtensions[fileExtension]
			if ok {
				availableVideos = append(availableVideos, entry.Name())
			}
		}
	}

	return availableVideos, nil
}
