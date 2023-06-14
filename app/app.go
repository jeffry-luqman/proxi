package app

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/fasthttp/websocket"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

const Version = "v0.0.1"

var (
	ConfigFile = "proxi.yml"
	client     = fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
	}
	upgrader   = websocket.FastHTTPUpgrader{}
	fileLogger = zerolog.Logger{}
)

func Run() {
	fmt.Println("Starting " + Fmt("proxi "+Version, Green, BlinkSlow))
	b, _ := os.ReadFile(ConfigFile)
	err := yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(Fmt(err.Error(), Red))
	}

	if config.Log.File.Enable {
		fileLogger = zerolog.New(&lumberjack.Logger{
			Filename:   config.Log.File.Filename,
			MaxSize:    config.Log.File.MaxSize,
			MaxAge:     config.Log.File.MaxAge,
			MaxBackups: config.Log.File.MaxBackups,
		}).With().Timestamp().Logger()
	}

	port := fmt.Sprintf("%v", config.Port)
	fmt.Println()
	fmt.Println("Proxi available at " + Fmt("http://localhost:"+port, Blue))
	if config.Metric.Enable {
		metric.Init()
	}
	server := &fasthttp.Server{
		Handler: handler,
	}
	if err := server.ListenAndServe(":" + port); err != nil {
		log.Fatal("Server", Fmt(err.Error(), Red))
	}
}

func handler(c *fasthttp.RequestCtx) {
	originalURL := string(c.Request.Header.RequestURI())
	defer c.Request.SetRequestURI(originalURL)

	pathPrefix, addr := getURL(originalURL)
	c.Request.SetRequestURI(addr)

	if scheme := getScheme(addr); len(scheme) > 0 {
		c.Request.URI().SetSchemeBytes(scheme)
	}
	c.Request.Header.Del("Connection")

	ctx := newCtx(&c.Request, originalURL, pathPrefix, addr)
	err := client.Do(&c.Request, &c.Response)
	ctx.logging(&c.Response, err)

	c.Response.Header.Del("Connection")
}

func getURL(originalURL string) (string, string) {
	for k, v := range config.Targets {
		if k != "/" && strings.HasPrefix(originalURL, k) {
			return k, v + originalURL
		}
	}
	return "/", config.Targets["/"] + originalURL
}

func getScheme(s string) []byte {
	uri := unsafe.Slice(unsafe.StringData(s), len(s))
	i := bytes.IndexByte(uri, '/')
	if i < 1 || uri[i-1] != ':' || i == len(uri)-1 || uri[i+1] != '/' {
		return nil
	}
	return uri[:i-1]
}
