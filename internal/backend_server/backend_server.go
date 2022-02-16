package backend_server

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/ronmount/ozon_go/internal/database"
	"github.com/ronmount/ozon_go/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var log = logrus.New()

type backendServer struct {
	storage models.Database
	router  *mux.Router
}

func NewServer() (*backendServer, error) {
	mode := os.Getenv("STORAGE_TYPE")
	if len(mode) == 0 {
		return nil, models.MissingStorageType{}
	}

	server := backendServer{}
	server.router = mux.NewRouter()
	server.router.StrictSlash(true)
	server.router.HandleFunc("/shortener/", server.createShortLinkHandler).Methods("POST")
	server.router.HandleFunc("/shortener/{url:[a-zA-Z0-9_]+}", server.getFullLinkHandler).Methods("GET")

	if mode == "memory" {
		server.storage, _ = database.NewMemoryStorage()
	} else if mode == "postgresql" {
		var err error
		server.storage, err = database.NewPSQLStorage()
		if err != nil {
			return nil, models.PSQLStorage{}
		}
	} else {
		return nil, models.WrongStorageType{}
	}
	return &server, nil
}

func SendJSON(w http.ResponseWriter, v interface{}) {
	log.Info(v)
	if answer, err := json.Marshal(v); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(answer)
	} else {
		handleErrors(w, models.HTTP500{})
	}
}

func handleErrors(w http.ResponseWriter, err error) {
	var status int

	if errors.As(err, &models.HTTP404{}) {
		status = http.StatusNotFound
	} else if errors.As(err, &models.HTTP500{}) {
		status = http.StatusInternalServerError
	}
	log.Error(status)
	http.Error(w, "", status)
}

func (b *backendServer) createShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(r.URL.Path)

	if fullLink := r.FormValue("url"); len(fullLink) > 0 {
		if links, err := b.storage.AddURL(fullLink); err != nil {
			handleErrors(w, err)
		} else {
			SendJSON(w, links)
		}
	} else {
		handleErrors(w, models.HTTP500{})
	}
}

func (b *backendServer) getFullLinkHandler(w http.ResponseWriter, r *http.Request) {
	log.Info(r.URL.Path)

	if shortLink, keyExists := mux.Vars(r)["url"]; keyExists {
		if links, err := b.storage.GetURL(shortLink); err != nil {
			handleErrors(w, err)
		} else {
			SendJSON(w, links)
		}
	} else {
		handleErrors(w, models.HTTP500{})
	}
}

func (b *backendServer) Run() error {
	if err := http.ListenAndServe("0.0.0.0:8080", b.router); err != nil {
		return models.RunError{}
	}
	return nil
}
