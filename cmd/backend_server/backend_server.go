package backend_server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ronmount/ozon_go/internal/config_parser"
	"github.com/ronmount/ozon_go/internal/database"
	"github.com/ronmount/ozon_go/internal/models"
	"log"
	"net/http"
	"os"
)

const description500 = "500 Internal Server Error"

type backendServer struct {
	storage models.Database
	router  *mux.Router
}

func NewServer() *backendServer {
	mode := os.Getenv("STORAGE_TYPE")
	if len(mode) == 0 {
		log.Fatal("STORAGE_TYPE env variable not found.")
	}

	server := backendServer{}
	server.router = mux.NewRouter()
	server.router.StrictSlash(true)
	server.router.HandleFunc("/shortener/", server.createShortLinkHandler).Methods("POST")
	server.router.HandleFunc("/shortener/{url:[a-zA-Z0-9_]+}", server.getFullLinkHandler).Methods("GET")

	if mode == "redis" {
		server.storage = database.NewRedisStorage()
	} else if mode == "postgresql" {
		var err error = nil
		config := config_parser.LoadPSQLConfigs()
		server.storage, err = database.NewPSQLStorage(*config)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Unknown database type.")
	}
	return &server
}

func sendJSON(w http.ResponseWriter, v interface{}) {
	answer, err := json.Marshal(v)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		_, wErr := w.Write(answer)
		if wErr != nil {
			log.Println(wErr)
		}
	} else {
		log.Println(err)
		http.Error(w, description500, http.StatusInternalServerError)
	}
}

func (b *backendServer) createShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	if fullLink := r.FormValue("url"); len(fullLink) > 0 {
		links, status := b.storage.AddURL(fullLink)
		log.Println(links, status)
		w.WriteHeader(status)
		if status != http.StatusInternalServerError {
			sendJSON(w, links)
		}
	} else {
		log.Println("Some error occurred while getting url from request.")
		http.Error(w, description500, http.StatusInternalServerError)
	}
}

func (b *backendServer) getFullLinkHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	if shortLink, keyExists := mux.Vars(r)["url"]; keyExists {
		links, status := b.storage.GetURL(shortLink)
		log.Println(links, status)
		w.WriteHeader(status)
		if status == http.StatusOK {
			sendJSON(w, links)
		}
	} else {
		http.Error(w, description500, http.StatusInternalServerError)
	}
}

func (b *backendServer) Run() {
	http.ListenAndServe("0.0.0.0:8080", b.router)
}
