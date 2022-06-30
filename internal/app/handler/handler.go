package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/romik1505/fileserver/internal/app/service"
)

type Handler struct {
	fileService *service.FileService
}

func NewHandler(fs *service.FileService) *Handler {
	return &Handler{
		fileService: fs,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/files/").Handler(http.StripPrefix("/files/", h.getFile())).Methods(http.MethodGet)
	router.HandleFunc("/upload", ErrorWrapper(h.saveFile)).Methods(http.MethodPost)
	router.HandleFunc("/", h.uploadForm).Methods(http.MethodGet)
	return router
}

type HandleFunc func(w http.ResponseWriter, r *http.Request) error

func ErrorWrapper(h HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
