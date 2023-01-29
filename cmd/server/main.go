package main

import (
	"os"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/zibbp/tethys/internal/file"
	transportHTTP "github.com/zibbp/tethys/internal/transport/http"
	"github.com/zibbp/tethys/internal/utils"
)

func Run() error {

	// Create upload folder
	err := os.MkdirAll("/tmp/uploads", 0755)
	if err != nil {
		log.Error().Err(err).Msg("Error creating upload folder")
		return err
	}

	// Old file cleanup
	c := cron.New()
	_, err = c.AddFunc("0 * * * *", utils.RemoveOldFiles)
	if err != nil {
		log.Error().Err(err).Msg("Error adding cron job")
		return err
	}
	c.Start()

	fileService := file.NewService()

	handler := transportHTTP.NewHandler(fileService)

	if err := handler.Serve(); err != nil {
		log.Error().Err(err).Msg("Error creating HTTP handler")
		return err
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal().Err(err).Msg("Error running server")
	}
}
