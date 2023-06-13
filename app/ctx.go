package app

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type Ctx struct {
	StartAt     time.Time
	FinishAt    time.Time
	Duration    time.Duration
	Method      string
	EndPoint    string
	CompleteURL string
	Header      map[string]string
	Body        []byte
	StatusCode  int
}

func newCtx(req *fasthttp.Request, originalURL, completeURL string) *Ctx {
	c := &Ctx{
		StartAt:     time.Now(),
		Method:      string(req.Header.Method()),
		CompleteURL: completeURL,
	}
	c.EndPoint, _, _ = strings.Cut(originalURL, "?")
	if config.Log.File.IncludeRequestHeaders {
		c.Header = map[string]string{}
		req.Header.VisitAll(func(key, value []byte) {
			c.Header[string(key)] = string(value)
		})
	}
	if config.Log.File.IncludeRequestBody {
		c.Body = req.Body()
	}
	if config.Log.Console.Enable && config.Log.Console.PrintRequestImmediately {
		fmt.Println(c.StartAt.Format(time.TimeOnly) + "\t" + c.fmtMethod() + "\t" + c.fmtURL())
	}

	return c
}

func (c *Ctx) logging(res *fasthttp.Response, err error) {
	c.FinishAt = time.Now()
	c.Duration = c.FinishAt.Sub(c.StartAt).Round(time.Millisecond)
	c.StatusCode = res.Header.StatusCode()
	if config.Log.Console.Enable {
		consoleLog := c.FinishAt.Format(time.TimeOnly) + " " + c.fmtDuration() + "\t" + c.fmtStatus() + "\t" + c.Method + " " + c.CompleteURL
		if err != nil {
			consoleLog += "\t" + Fmt(err.Error(), Red)
		}
		fmt.Println(consoleLog)
	}
	if config.Log.File.Enable {
		logger := fileLogger.Log().
			Time("start", c.StartAt).
			Time("finish", c.FinishAt).
			Dur("duration", c.Duration).
			Str("method", c.Method).
			Str("endpoint", c.EndPoint).
			Str("url", c.CompleteURL).
			Int("status", c.StatusCode).
			Str("level", c.level(err))
		if config.Log.File.IncludeRequestHeaders {
			logger = logger.Any("headers", c.Header)
		}
		if config.Log.File.IncludeRequestBody {
			logger = logger.Bytes("body", c.Body)
		}
		if config.Log.File.IncludeResponseHeaders {
			responseHeader := map[string]string{}
			res.Header.VisitAll(func(key, value []byte) {
				responseHeader[string(key)] = string(value)
			})
			logger = logger.Any("res_headers", responseHeader)
		}
		if config.Log.File.IncludeResponseBody {
			logger = logger.Bytes("res_body", res.Body())
		}
		logger.Send()
	}
}

func (c *Ctx) level(err error) string {
	if err != nil {
		return "error"
	}
	if c.StatusCode >= http.StatusOK {
		if c.StatusCode < http.StatusBadRequest {
			return "info"
		} else if c.StatusCode < http.StatusInternalServerError {
			return "warning"
		}
	}
	return "error"
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
