package main

import (
	"log"
	"net/http"
	"github.com/yosssi/ace"
	"github.com/yosssi/ace-proxy"
	"github.com/gin-gonic/gin"
	"flag"
	"github.com/gin-gonic/contrib/gzip"
)

var p = proxy.New(&ace.Options{BaseDir: "views"})

func main() {
	var httpPort string
	var httpsPort string
	var certFile string
	var keyFile string
	var prod bool

	flag.StringVar(&httpPort, "port", "8080", "http port")
	flag.StringVar(&httpsPort, "https", "", "https port")
	flag.StringVar(&certFile, "cert", "", "cert file for https")
	flag.StringVar(&keyFile, "key", "", "key file for https")
	flag.BoolVar(&prod, "prod", false, "is in production")

	flag.Parse()

	if prod {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Handle(http.MethodGet, "/", HomeHandler)
	r.Handle(http.MethodGet, "/stickers", StickersHandler)
	r.Handle(http.MethodPost, "/stickers/sticker", UploadStickerHandler)
	r.StaticFS("/public/", http.Dir("./static/"))
	r.StaticFS("/.well-known/", http.Dir("./.well-known/"))

	r.NoRoute(StatusNotFoundHandler)


	if certFile != "" && keyFile != "" && httpsPort != "" {
		go r.Run(":" + httpPort)
		log.Fatal(r.RunTLS(":" + httpPort, certFile, keyFile))
	} else {
		log.Fatal(r.Run(":" + httpPort))
	}

}

func HomeHandler(c *gin.Context) {
	runTemplate(c, "index", nil)
}

func StatusNotFoundHandler(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNotFound)

	runTemplate(c, "error", map[string]string{
		"Status": "404",
		"Message": "Page not found.",
	})
}


func runTemplate(c *gin.Context, innerPath string, data interface{}) {
	tpl, err := p.Load("base", innerPath, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Header("Content-Type", "text/HTML")

	if err := tpl.Execute(c.Writer, data); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}