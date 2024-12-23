package app

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"

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
	fileLogger = zerolog.Logger{}
)

func Run() {
	fmt.Println("Starting " + Fmt("proxi "+Version, Green, BlinkSlow))
	b, _ := os.ReadFile(ConfigFile)
	err := yaml.Unmarshal(b, &Conf)
	if err != nil {
		log.Fatal(Fmt(err.Error(), Red))
	}
	if Conf.TargetStr != "" {
		for _, t := range strings.Split(Conf.TargetStr, ";") {
			pathPrefix, targetURL, _ := strings.Cut(strings.Trim(t, " "), " ")
			if pathPrefix != "" {
				Conf.Targets[pathPrefix] = targetURL
			}
		}
	}
	if Conf.Log.Console.Disable {
		Conf.Log.Console.Enable = false
	}
	if Conf.Log.File.Filename != "" {
		Conf.Log.File.Enable = true
	}
	if Conf.Log.File.Enable {
		fileLogger = zerolog.New(&lumberjack.Logger{
			Filename:   Conf.Log.File.Filename,
			MaxSize:    Conf.Log.File.MaxSize,
			MaxAge:     Conf.Log.File.MaxAge,
			MaxBackups: Conf.Log.File.MaxBackups,
		}).With().Timestamp().Logger()
	}
	if Conf.Metric.Port > 0 {
		Conf.Metric.Enable = true
	}
	if Conf.Metric.Enable {
		metric.Init()
	}

	port := fmt.Sprintf("%v", Conf.Port)
	fmt.Println()
	fmt.Println("Proxi available at " + Fmt("http://localhost:"+port, Blue))
	server := &fasthttp.Server{
		ReadBufferSize: 16384,
		Handler:        handler,
	}
	if err := server.ListenAndServe(":" + port); err != nil {
		log.Fatal("Server", Fmt(err.Error(), Red))
	}
}

func handler(c *fasthttp.RequestCtx) {
	originalURL := string(c.Request.Header.RequestURI())
	targetPrefix, targetURL, targetDir := getTarget(originalURL)
	if targetDir != "" {
		rootDir, stripSlashes, _ := strings.Cut(targetDir, " ")
		stripSlashesInt, _ := strconv.Atoi(stripSlashes)
		staticFileHandler := fasthttp.FSHandler(rootDir, stripSlashesInt)
		staticFileHandler(c)
	} else {
		defer c.Request.SetRequestURI(originalURL)
		c.Request.SetRequestURI(targetURL)
		if scheme := getScheme(targetURL); len(scheme) > 0 {
			c.Request.URI().SetSchemeBytes(scheme)
		}
		c.Request.Header.Del("Connection")
		ctx := newCtx(&c.Request, originalURL, targetPrefix, targetURL)
		err := client.Do(&c.Request, &c.Response)
		ctx.logging(&c.Response, err)
		c.Response.Header.Del("Connection")
	}
}

func getTarget(originalURL string) (string, string, string) {
	for k, v := range Conf.Targets {
		if k != "/" && strings.HasPrefix(originalURL, k) {
			if strings.HasPrefix(v, "http") {
				return k, v + originalURL, ""
			}
			return k, v + originalURL, v
		}
	}
	if strings.HasPrefix(Conf.Targets["/"], "http") {
		return "/", Conf.Targets["/"] + originalURL, ""
	}
	return "/", Conf.Targets["/"] + originalURL, Conf.Targets["/"]
}

func getScheme(s string) []byte {
	uri := unsafe.Slice(unsafe.StringData(s), len(s))
	i := bytes.IndexByte(uri, '/')
	if i < 1 || uri[i-1] != ':' || i == len(uri)-1 || uri[i+1] != '/' {
		return nil
	}
	return uri[:i-1]
}
