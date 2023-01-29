package utils

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func SanitizeFileName(fileName string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r == '/' || r == '\\' || r == ':' || r == '*' ||
			r == '?' || r == '"' || r == '<' || r == '>' || r == ' ' ||
			r == '|' || r == '%':
			return '_'
		default:
			return r
		}
	}, fileName)
}

func RemoveOldFiles() {
	dir := "/tmp/uploads"

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip root folder
		if path == dir {
			return nil
		}

		if time.Since(info.ModTime()) > time.Hour {
			err := os.RemoveAll(path)
			if err != nil {
				return err
			}
		}

		return nil

	})

	if err != nil {
		log.Error().Err(err).Msg("Error removing old files")
	}
}
