package http

import (
	"fmt"
	"gush/services"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const defaultHTTPPort = "8080"

func buildRedirectURL(r *http.Request, hash string) string {
	var scheme, host string

	if r.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}
	host = r.Host

	return fmt.Sprintf("%v://%v/%v", scheme, host, hash)
}

func postShortURL(w http.ResponseWriter, r *http.Request) {

	// Validate "Content-Type" header
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		log.Printf("An error occurred when Content-Type header: %v", err)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	if mediaType != "text/plain" {
		log.Printf("Unsupported media type: %s", mediaType)
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rawBodyAsBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("An error occurred when parsing body: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	urlToShorten := string(rawBodyAsBytes)

	hash, ok := services.SetShortURL(urlToShorten)

	if !ok {
		log.Printf("An error occurred generating hash...")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedURL := buildRedirectURL(r, hash)
	hashedURLBytes := []byte(hashedURL)

	w.Header().Set("Content-Type", "text/plain")

	w.Write(hashedURLBytes)
}

func getRedirectShortURL(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	hash := vars["hash"]

	urlInfo, ok := services.GetShortURLInfo(hash)

	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	log.Printf("Redirecting to: %v", urlInfo.URL)
	http.Redirect(w, r, urlInfo.URL, http.StatusPermanentRedirect)
}

func getURLInfo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	hash := vars["hash"]

	urlInfo, ok := services.GetShortURLInfo(hash)

	if !ok {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	hashedURLBytes := []byte(urlInfo.URL)

	w.Header().Set("Last-Modified", urlInfo.CreatedAt.Format(http.TimeFormat))
	w.Header().Set("Content-Type", "text/plain")

	w.Write(hashedURLBytes)
}

func getHTTPEnvPort() string {
	envPort, ok := os.LookupEnv("HTTP_PORT")

	if !ok || len(envPort) == 0 {
		envPort = defaultHTTPPort
	}

	return envPort
}

// Run - runs HTTP environment
func Run() {

	r := mux.NewRouter()

	r.HandleFunc("/", postShortURL).Methods(http.MethodPost)
	r.HandleFunc("/{hash}", getRedirectShortURL).Methods(http.MethodGet)
	r.HandleFunc("/info/{hash}", getURLInfo).Methods(http.MethodGet)

	port := getHTTPEnvPort()

	log.Printf("Listening HTTP in port %v ...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
