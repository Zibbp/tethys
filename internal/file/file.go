package file

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/tethys/internal/utils"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) UploadFile(ctx context.Context, r *http.Request) (string, error) {
	id, err := gonanoid.New(6)
	if err != nil {
		return "", fmt.Errorf("error generating id: %w", err)
	}

	fileName := chi.URLParam(r, "fileName")
	fileName = utils.SanitizeFileName(fileName)

	// Create folder
	err = os.MkdirAll(fmt.Sprintf("/tmp/uploads/%s", id), 0755)
	if err != nil {
		return "", fmt.Errorf("error creating folder: %w", err)
	}

	// Create temp file
	file, err := os.Create(fmt.Sprintf("/tmp/uploads/%s/%s", id, fileName))
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error copying file to temp file")
		return "", fmt.Errorf("error copying file to temp file: %w", err)
	}

	log.Info().Msgf("Uploaded file %s", fileName)

	return fmt.Sprintf("%s/%s", id, fileName), nil
}

func (s *Service) GetFile(ctx context.Context, id string, fileName string) ([]byte, error) {
	file, err := os.ReadFile(fmt.Sprintf("/tmp/uploads/%s/%s", id, fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return file, nil
}
