package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/romik1505/fileserver/internal/app/config"
)

const (
	MaxUploadFile = 10 << 20 // 10 MB
)

func (h *Handler) getFile() http.Handler {
	return http.FileServer(h.fileService)
}

func (h *Handler) saveFile(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	r.ParseMultipartForm(MaxUploadFile)

	file, handler, err := r.FormFile(config.GetValue(config.MultipartFileKey))
	if err != nil {
		return err
	}
	defer file.Close()

	err = h.fileService.SaveFile(ctx, file, handler.Filename)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"url": fmt.Sprintf("%s/files/%s", config.GetValue(config.Domain), handler.Filename),
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) uploadForm(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./index.html")
}
