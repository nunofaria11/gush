package main

import (
	"gush/entities"
	"gush/shortener"
	"gush/utils"
	"io/ioutil"
	"net/url"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

// PORT - The port the service is running on
const PORT = "8080"

func createShortURL(ctx iris.Context) {

	var hash string
	rawBodyAsBytes, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		ctx.Writef("%v", err)
	}

	urlToShorten := string(rawBodyAsBytes)
	urlInfo := entities.NewURLInfo(urlToShorten)

	ok := false

	for !ok {
		hash = utils.RandomString(8)
		_, ok = shortener.SetShortURL(hash, urlInfo)
	}

	fullRequestURI := ctx.FullRequestURI()
	url, err := url.Parse(fullRequestURI)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	hashedURL := url.Scheme + "://" + url.Host + "/" + hash
	hashedURLBytes := []byte(hashedURL)

	ctx.Write(hashedURLBytes)
}

func redirectShortURL(ctx iris.Context) {
	hash := ctx.Params().GetString("hash")
	urlInfo, ok := shortener.GetShortURL(hash)

	if !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	ctx.Redirect(urlInfo.URL, iris.StatusPermanentRedirect)
}

// GetURLInfo - retrieves a URL info
func getURLInfo(ctx iris.Context) {
	hash := ctx.Params().GetString("hash")
	urlInfo, ok := shortener.GetShortURL(hash)

	if !ok {
		ctx.StatusCode(iris.StatusNotFound)
		return
	}

	ctx.Header("Last-Modified", urlInfo.Time.UTC().String())
	ctx.Writef(urlInfo.URL)
}

func main() {

	app := iris.New()
	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())

	app.Post("/", createShortURL)
	app.Get("/{hash:string}", redirectShortURL)
	app.Get("/info/{hash:string}", getURLInfo)

	app.Run(iris.Addr(":"+PORT), iris.WithoutServerError(iris.ErrServerClosed))
}
