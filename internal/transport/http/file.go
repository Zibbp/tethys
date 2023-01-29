package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FileService interface {
	UploadFile(ctx context.Context, r *http.Request) (string, error)
	GetFile(ctx context.Context, id string, fileName string) ([]byte, error)
}

func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {

	upload, err := h.Service.FileService.UploadFile(r.Context(), r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	protocol := "http"
	if r.TLS != nil {
		protocol = "https"
	}

	url := fmt.Sprintf("%s://%s/%s", protocol, r.Host, upload)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(url))
}

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	id := chi.URLParam(r, "id")
	file, err := h.Service.FileService.GetFile(r.Context(), id, fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.WriteHeader(http.StatusOK)
	w.Write(file)

}
