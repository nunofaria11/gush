package http

import (
	"fmt"
	"gush/services"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
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

	hash, err := services.SetShortURL(urlToShorten)

	if err != nil {
		log.Printf("An error occurred generating hash: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedURL := buildRedirectURL(r, hash)
	hashedURLBytes := []byte(hashedURL)

	w.Header().Set("Content-Type", "text/plain")

	w.Write(hashedURLBytes)
}

func getRedirectShortURL(w http.ResponseWriter, r *http.Request) {
	var redirectURL string

	vars := mux.Vars(r)
	hash := vars["hash"]

	urlInfo, err := services.GetShortURLInfo(hash)

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	redirectURL = urlInfo.URL
	parsedURL, err := url.Parse(redirectURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(parsedURL.Scheme) == 0 {
		redirectURL = "http://" + redirectURL
	}

	log.Printf("Redirecting to: %v", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
}

func getURLInfo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	hash := vars["hash"]

	urlInfo, err := services.GetShortURLInfo(hash)

	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	hashedURLBytes := []byte(urlInfo.URL)

	w.Header().Set("Last-Modified", urlInfo.CreatedAt.Format(http.TimeFormat))
	w.Header().Set("Content-Type", "text/plain")

	w.Write(hashedURLBytes)
}

func getHTTPEnvPort() string {
	envPort, ok := os.LookupEnv("PORT")

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
