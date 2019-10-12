package http

import (
	"gush/services"
	"io/ioutil"
	"log"
	"mime"
	"net/url"
	"os"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

// defaultHTTPPort The default port the service is running on
const defaultHTTPPort = "8080"

// CreateShortURL Creates a short URL
func createShortURL(ctx iris.Context) {

	// Validate "Content-Type" header
	mediaType, _, err := mime.ParseMediaType(ctx.GetHeader("Content-Type"))
	if err != nil {
		log.Printf("An error occurred when Content-Type header: %v", err)
		ctx.StatusCode(iris.StatusUnsupportedMediaType)
		return
	}
	if mediaType != "text/plain" {
		log.Printf("Unsupported media type: %s", mediaType)
		ctx.StatusCode(iris.StatusUnsupportedMediaType)
		return
	}

	rawBodyAsBytes, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		log.Printf("An error occurred when parsing body: %v", err)
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	urlToShorten := string(rawBodyAsBytes)

	hash, ok := services.SetShortURL(urlToShorten)

	if !ok {
		log.Printf("An error occurred generating hash...")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	fullRequestURI := ctx.FullRequestURI()
	url, err := url.Parse(fullRequestURI)
	if err != nil {
		log.Printf("An error occurred extracting request URI...")
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	hashedURL := url.Scheme + "://" + url.Host + "/" + hash
	hashedURLBytes := []byte(hashedURL)

	ctx.Write(hashedURLBytes)
}

func redirectShortURL(ctx iris.Context) {
	hash := ctx.Params().GetString("hash")
	urlInfo, ok := services.GetShortURLInfo(hash)

	if !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	log.Printf("Redirecting to: %v", urlInfo.URL)
	ctx.Redirect(urlInfo.URL, iris.StatusPermanentRedirect)
}

func getURLInfo(ctx iris.Context) {

	hash := ctx.Params().GetString("hash")
	urlInfo, ok := services.GetShortURLInfo(hash)

	if !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	ctx.Header("Last-Modified", urlInfo.Time.UTC().String())
	ctx.Writef(urlInfo.URL)
}

// Run - runs HTTP environment
func Run() {

	app := iris.New()
	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())

	app.Post("/", createShortURL)
	app.Get("/{hash:string}", redirectShortURL)
	app.Get("/info/{hash:string}", getURLInfo)

	envPort, ok := os.LookupEnv("HTTP_PORT")
	if !ok || len(envPort) == 0 {
		envPort = defaultHTTPPort
	}

	port := ":" + envPort

	app.Run(iris.Addr(port), iris.WithoutServerError(iris.ErrServerClosed))
}
