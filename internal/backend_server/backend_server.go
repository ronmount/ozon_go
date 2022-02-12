package backend_server

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/ronmount/ozon_go/internal/config_parser"
	"github.com/ronmount/ozon_go/internal/database"
	"github.com/ronmount/ozon_go/internal/models"
	"github.com/ronmount/ozon_go/internal/my_errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var log = logrus.New()

type backendServer struct {
	storage models.Database
	router  *mux.Router
}

func handleErrors(w http.ResponseWriter, err error) {
	var status int

	if errors.As(err, &my_errors.HTTP404{}) {
		status = http.StatusNotFound
	} else if errors.As(err, &my_errors.HTTP500{}) {
		status = http.StatusInternalServerError
	}
	log.Error(status)
	http.Error(w, "", status)
}

func NewServer() (*backendServer, error) {
	mode := os.Getenv("STORAGE_TYPE")
	if len(mode) == 0 {
		return nil, my_errors.MissingStorageType{}
	}

	server := backendServer{}
	server.router = mux.NewRouter()
	server.router.StrictSlash(true)
	server.router.HandleFunc("/shortener/", server.createShortLinkHandler).Methods("POST")
	server.router.HandleFunc("/shortener/{url:[a-zA-Z0-9_]+}", server.getFullLinkHandler).Methods("GET")

	var err error
	if mode == "redis" {
		server.storage, err = database.NewRedisStorage()
		if err != nil {
			return nil, my_errors.RedisStorage{}
		}
	} else if mode == "postgresql" {
		config, err := config_parser.LoadPSQLConfigs()
		if err != nil {
			return nil, my_errors.PSQLStorage{}
		}
		server.storage, err = database.NewPSQLStorage(*config)
		if err != nil {
			return nil, my_errors.PSQLStorage{}
		}
	} else {
		return nil, my_errors.WrongStorageType{}
	}
	return &server, nil
}

func sendJSON(w http.ResponseWriter, v interface{}) {
	log.Info(v)
	if answer, err := json.Marshal(v); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(answer)
	} else {
		handleErrors(w, my_errors.HTTP500{})
	}
}

func (b *backendServer) createShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	if fullLink := r.FormValue("url"); len(fullLink) > 0 {
		if links, err := b.storage.AddURL(fullLink); err != nil {
			handleErrors(w, err)
		} else {
			sendJSON(w, links)
		}
	} else {
		handleErrors(w, my_errors.HTTP500{})
	}
}

func (b *backendServer) getFullLinkHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	if shortLink, keyExists := mux.Vars(r)["url"]; keyExists {
		if links, err := b.storage.GetURL(shortLink); err != nil {
			handleErrors(w, err)
		} else {
			sendJSON(w, links)
		}
	} else {
		handleErrors(w, my_errors.HTTP500{})
	}
}

func (b *backendServer) Run() error {
	if err := http.ListenAndServe("0.0.0.0:8080", b.router); err != nil {
		return my_errors.RunError{}
	}
	return nil
}