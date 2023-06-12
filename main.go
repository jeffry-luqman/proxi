package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/mattn/go-isatty"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
)

const version = "23.06.121325"

const (
	Reset uint8 = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
)

const (
	Black uint8 = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
)

var (
	port    = "8080"
	debug   = false
	targets = map[string]string{"/": "http://localhost:3000"}

	configFile = "proxi.yml"
	client     = fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
	}
)

func main() {
	server := &fasthttp.Server{
		Handler: proxyHandler,
	}
	fmt.Println("Starting " + format("proxi "+version, Green))
	fmt.Println("Please open " + format("http://localhost:"+port, Blue, BlinkSlow) + " in browser")
	if err := server.ListenAndServe(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %s", err)
	}
}

func init() {
	b, _ := os.ReadFile(configFile)
	var c struct {
		Port    int
		Debug   bool
		Targets map[string]string
	}
	yaml.Unmarshal(b, &c)
	for k, v := range c.Targets {
		targets[k] = v
	}
	if c.Port > 0 {
		port = fmt.Sprintf("%v", c.Port)
	}
	debug = c.Debug
}

func proxyHandler(c *fasthttp.RequestCtx) {
	req := &c.Request
	res := &c.Response

	originalURL := string(req.Header.RequestURI())
	defer req.SetRequestURI(originalURL)

	addr := getURL(originalURL)
	req.SetRequestURI(addr)

	if scheme := getScheme(addr); len(scheme) > 0 {
		req.URI().SetSchemeBytes(scheme)
	}
	req.Header.Del("Connection")
	reqInfo := ""
	if debug {
		reqInfo = formatRequest(c.Method(), addr)
		log.Println(reqInfo)
	}
	now := time.Now()
	err := client.Do(req, res)
	if debug {
		fmt.Println(formatResponse(res.Header.StatusCode(), now, reqInfo, err))
	}
	res.Header.Del("Connection")
}

func getURL(originalURL string) string {
	for k, v := range targets {
		if k != "/" && strings.HasPrefix(originalURL, k) {
			return v + originalURL
		}
	}
	return targets["/"] + originalURL
}

func getScheme(s string) []byte {
	uri := unsafe.Slice(unsafe.StringData(s), len(s))
	i := bytes.IndexByte(uri, '/')
	if i < 1 || uri[i-1] != ':' || i == len(uri)-1 || uri[i+1] != '/' {
		return nil
	}
	return uri[:i-1]
}

func formatRequest(method []byte, addr string) string {
	return format(string(method)+"\t"+addr, Cyan)
}

func formatResponse(code int, now time.Time, req string, err error) string {
	res := formatStatus(code) + "\t" + formatDuration(now) + "\t" + req
	if err != nil {
		res += "\t" + format(err.Error(), Red)
	}
	return res
}

func formatStatus(code int) string {
	if code >= http.StatusOK {
		if code < http.StatusBadRequest {
			return format(code, Green)
		} else if code < http.StatusInternalServerError {
			return format(code, Yellow)
		}
	}
	return format(code, Red)
}

func formatDuration(now time.Time) string {
	d := time.Now().Sub(now)
	if d < time.Second/2 {
		return format(d, Magenta)
	} else if d < time.Second {
		return format(d, Yellow)
	}
	return format(d, Red)
}

func format(text any, attribute ...uint8) string {
	s := fmt.Sprintf("%v", text)
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return s
	}
	format := make([]string, len(attribute))
	for i, v := range attribute {
		format[i] = fmt.Sprintf("%v", v)
	}
	return "\x1b[" + strings.Join(format, ";") + "m" + s + "\x1b[0m"
}
