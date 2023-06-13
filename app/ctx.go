package app

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type Ctx struct {
	StartAt        time.Time
	FinishAt       time.Time
	Duration       time.Duration
	Method         string
	EndPoint       string
	CompleteURL    string
	Header         map[string]string
	Body           string
	StatusCode     int
	ResponseHeader map[string]string
	ResponseBody   string
}

func newCtx(req *fasthttp.Request, originalURL string) *Ctx {
	c := &Ctx{
		StartAt:     time.Now(),
		Method:      string(req.Header.Method()),
		CompleteURL: string(req.RequestURI()),
	}
	c.EndPoint, _, _ = strings.Cut(originalURL, "?")
	c.Header = map[string]string{}
	if config.Log.File.IncludeRequestHeaders {
		req.Header.VisitAll(func(key, value []byte) {
			c.Header[string(key)] = string(value)
		})
	}
	if config.Log.File.IncludeRequestBody {
		c.Body = string(req.Body())
	}
	if config.Log.Console.Enable && config.Log.Console.PrintRequestImmediately {
		fmt.Println(c.StartAt.String() + "\n" + c.fmtMethod() + "\t" + c.fmtURL())
	}

	return c
}

func (c *Ctx) logging(res *fasthttp.Response, err error) {
	c.FinishAt = time.Now()
	c.Duration = c.FinishAt.Sub(c.StartAt).Round(time.Millisecond)
	c.StatusCode = res.Header.StatusCode()
	if config.Log.Console.Enable {
		fmt.Println(c.fmtStatus() + "\t" + c.fmtDuration() + "\t" + c.fmtMethod() + "\t" + c.fmtURL())
	}
	c.ResponseHeader = map[string]string{}
	if config.Log.File.IncludeResponseHeaders {
		res.Header.VisitAll(func(key, value []byte) {
			c.ResponseHeader[string(key)] = string(value)
		})
	}
	if config.Log.File.IncludeResponseBody {
		c.ResponseBody = string(res.Body())
	}
	if config.Log.File.Enable {
		fileLogger.Info().Msg(c.Method + "\t" + c.CompleteURL)
	}
}

func (c *Ctx) fmtMethod() string {
	return Fmt(c.Method, Green)
}

func (c *Ctx) fmtURL() string {
	return Fmt(c.CompleteURL, Cyan)
}

func (c *Ctx) fmtStatus() string {
	if c.StatusCode >= http.StatusOK {
		if c.StatusCode < http.StatusBadRequest {
			return Fmt(c.StatusCode, Green)
		} else if c.StatusCode < http.StatusInternalServerError {
			return Fmt(c.StatusCode, Yellow)
		}
	}
	return Fmt(c.StatusCode, Red)
}

func (c *Ctx) fmtDuration() string {
	if c.Duration < time.Second/2 {
		return Fmt(c.Duration, Magenta)
	} else if c.Duration < time.Second {
		return Fmt(c.Duration, Yellow)
	}
	return Fmt(c.Duration, Red)
}
