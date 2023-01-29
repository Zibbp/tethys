package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Router  chi.Router
	Server  *http.Server
	Service Service
}

type Service struct {
	FileService FileService
}

type Response struct {
	Message string `json:"message"`
}

func NewHandler(fileService FileService) *Handler {
	log.Info().Msg("Creating HTTP handler")
	h := &Handler{
		Service: Service{
			FileService: fileService,
		},
	}

	h.Router = chi.NewRouter()

	h.mapRoutes()

	h.Server = &http.Server{
		Addr:    ":4000",
		Handler: h.Router,
	}

	return h

}

func (h *Handler) mapRoutes() {
	h.Router.Get("/", h.GetRoot)
	h.Router.Put("/{fileName}", h.UploadFile)
	h.Router.Get("/{id}/{fileName}", h.GetFile)
}

func (h *Handler) GetRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	text :=
		`Welcome to Tethys.
Upload a file with:
curl --upload-file file https://example.com
Download a file with:
curl https://example.com/file -o file
Files are deleted after one hour.
`
	w.Write([]byte(text))
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("Error starting HTTP server")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Info().Msg("Shutting down")
	return nil
}
