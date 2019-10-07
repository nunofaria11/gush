package main

import (
	"io/ioutil"
	"time"

	"github.com/kataras/iris"

	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

// URLInfo - type used to hold URL info in memory
type URLInfo struct {
	url       string
	timestamp time.Time
}

// PORT - The port the service is running on
const PORT = "8080"

// shortURLMap - the map used to store URLs in memory
var shortURLMap map[string]string

func createShortURL(ctx iris.Context) {

	rawBodyAsBytes, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		ctx.Writef("%v", err)
	}

	url := string(rawBodyAsBytes)
}

func redirectShortURL(ctx iris.Context) {
	// TODO: 301 Moved Permanently
}

// GetURLInfo - retrieves a URL info
func getURLInfo(ctx iris.Context) {
	url := ctx.Params().GetString("url")

	// TODO
	//  1. set "Last-Modified" date header
	//  2. set "Content-Type" as "text/plain"
}

func main() {

	app := iris.New()
	app.Logger().SetLevel("debug")

	app.Use(recover.New())
	app.Use(logger.New())

	shortURLMap = make(map[string]string)

	app.Post("/", createShortURL)
	app.Get("/{url:string}", redirectShortURL)
	app.Get("/info/{url:string}", getURLInfo)

	app.Run(iris.Addr(":"+PORT), iris.WithoutServerError(iris.ErrServerClosed))
}
