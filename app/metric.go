package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var metric = &Metric{}

var (
	upgrader     = websocket.Upgrader{}
	activeConn   = make(map[*websocket.Conn]struct{})
	activeConnMu sync.Mutex
)

type Metric struct {
	StartAt time.Time     `json:"start_at"`
	Uptime  time.Duration `json:"uptime"`

	LevelCount         map[string]uint64        `json:"level.count"`
	LevelDurationAvg   map[string]time.Duration `json:"level.duration.avg"`
	LevelDurationMax   map[string]time.Duration `json:"level.duration.max"`
	LevelDurationMin   map[string]time.Duration `json:"level.duration.min"`
	LevelDurationTotal map[string]time.Duration `json:"level.duration.total"`
	LevelLastCall      map[string]time.Time     `json:"level.last_call"`

	StatusCodeCount         map[int]uint64        `json:"status_code.count"`
	StatusCodeDurationAvg   map[int]time.Duration `json:"status_code.duration.avg"`
	StatusCodeDurationMax   map[int]time.Duration `json:"status_code.duration.max"`
	StatusCodeDurationMin   map[int]time.Duration `json:"status_code.duration.min"`
	StatusCodeDurationTotal map[int]time.Duration `json:"status_code.duration.total"`
	StatusCodeLastCall      map[int]time.Time     `json:"status_code.last_call"`

	MethodCount         map[string]uint64        `json:"method.count"`
	MethodDurationAvg   map[string]time.Duration `json:"method.duration.avg"`
	MethodDurationMax   map[string]time.Duration `json:"method.duration.max"`
	MethodDurationMin   map[string]time.Duration `json:"method.duration.min"`
	MethodDurationTotal map[string]time.Duration `json:"method.duration.total"`
	MethodLastCall      map[string]time.Time     `json:"method.last_call"`

	PathPrefixSuccessCount         map[string]uint64        `json:"path_prefix.success.count"`
	PathPrefixSuccessDurationAvg   map[string]time.Duration `json:"path_prefix.success.duration.avg"`
	PathPrefixSuccessDurationMax   map[string]time.Duration `json:"path_prefix.success.duration.max"`
	PathPrefixSuccessDurationMin   map[string]time.Duration `json:"path_prefix.success.duration.min"`
	PathPrefixSuccessDurationTotal map[string]time.Duration `json:"path_prefix.success.duration.total"`
	PathPrefixSuccessLastCall      map[string]time.Time     `json:"path_prefix.success.last_call"`
	PathPrefixWarningCount         map[string]uint64        `json:"path_prefix.warning.count"`
	PathPrefixWarningDurationAvg   map[string]time.Duration `json:"path_prefix.warning.duration.avg"`
	PathPrefixWarningDurationMax   map[string]time.Duration `json:"path_prefix.warning.duration.max"`
	PathPrefixWarningDurationMin   map[string]time.Duration `json:"path_prefix.warning.duration.min"`
	PathPrefixWarningDurationTotal map[string]time.Duration `json:"path_prefix.warning.duration.total"`
	PathPrefixWarningLastCall      map[string]time.Time     `json:"path_prefix.warning.last_call"`
	PathPrefixErrorCount           map[string]uint64        `json:"path_prefix.error.count"`
	PathPrefixErrorDurationAvg     map[string]time.Duration `json:"path_prefix.error.duration.avg"`
	PathPrefixErrorDurationMax     map[string]time.Duration `json:"path_prefix.error.duration.max"`
	PathPrefixErrorDurationMin     map[string]time.Duration `json:"path_prefix.error.duration.min"`
	PathPrefixErrorDurationTotal   map[string]time.Duration `json:"path_prefix.error.duration.total"`
	PathPrefixErrorLastCall        map[string]time.Time     `json:"path_prefix.error.last_call"`
	PathPrefixTotalCount           map[string]uint64        `json:"path_prefix.total.count"`
	PathPrefixTotalDurationAvg     map[string]time.Duration `json:"path_prefix.total.duration.avg"`
	PathPrefixTotalDurationMax     map[string]time.Duration `json:"path_prefix.total.duration.max"`
	PathPrefixTotalDurationMin     map[string]time.Duration `json:"path_prefix.total.duration.min"`
	PathPrefixTotalDurationTotal   map[string]time.Duration `json:"path_prefix.total.duration.total"`
	PathPrefixTotalLastCall        map[string]time.Time     `json:"path_prefix.total.last_call"`

	BaseURLSuccessCount         map[string]uint64        `json:"base_url.success.count"`
	BaseURLSuccessDurationAvg   map[string]time.Duration `json:"base_url.success.duration.avg"`
	BaseURLSuccessDurationMax   map[string]time.Duration `json:"base_url.success.duration.max"`
	BaseURLSuccessDurationMin   map[string]time.Duration `json:"base_url.success.duration.min"`
	BaseURLSuccessDurationTotal map[string]time.Duration `json:"base_url.success.duration.total"`
	BaseURLSuccessLastCall      map[string]time.Time     `json:"base_url.success.last_call"`
	BaseURLWarningCount         map[string]uint64        `json:"base_url.warning.count"`
	BaseURLWarningDurationAvg   map[string]time.Duration `json:"base_url.warning.duration.avg"`
	BaseURLWarningDurationMax   map[string]time.Duration `json:"base_url.warning.duration.max"`
	BaseURLWarningDurationMin   map[string]time.Duration `json:"base_url.warning.duration.min"`
	BaseURLWarningDurationTotal map[string]time.Duration `json:"base_url.warning.duration.total"`
	BaseURLWarningLastCall      map[string]time.Time     `json:"base_url.warning.last_call"`
	BaseURLErrorCount           map[string]uint64        `json:"base_url.error.count"`
	BaseURLErrorDurationAvg     map[string]time.Duration `json:"base_url.error.duration.avg"`
	BaseURLErrorDurationMax     map[string]time.Duration `json:"base_url.error.duration.max"`
	BaseURLErrorDurationMin     map[string]time.Duration `json:"base_url.error.duration.min"`
	BaseURLErrorDurationTotal   map[string]time.Duration `json:"base_url.error.duration.total"`
	BaseURLErrorLastCall        map[string]time.Time     `json:"base_url.error.last_call"`
	BaseURLTotalCount           map[string]uint64        `json:"base_url.total.count"`
	BaseURLTotalDurationAvg     map[string]time.Duration `json:"base_url.total.duration.avg"`
	BaseURLTotalDurationMax     map[string]time.Duration `json:"base_url.total.duration.max"`
	BaseURLTotalDurationMin     map[string]time.Duration `json:"base_url.total.duration.min"`
	BaseURLTotalDurationTotal   map[string]time.Duration `json:"base_url.total.duration.total"`
	BaseURLTotalLastCall        map[string]time.Time     `json:"base_url.total.last_call"`

	EndPointSuccessCount         map[string]uint64        `json:"end_point.success.count"`
	EndPointSuccessDurationAvg   map[string]time.Duration `json:"end_point.success.duration.avg"`
	EndPointSuccessDurationMax   map[string]time.Duration `json:"end_point.success.duration.max"`
	EndPointSuccessDurationMin   map[string]time.Duration `json:"end_point.success.duration.min"`
	EndPointSuccessDurationTotal map[string]time.Duration `json:"end_point.success.duration.total"`
	EndPointSuccessLastCall      map[string]time.Time     `json:"end_point.success.last_call"`
	EndPointWarningCount         map[string]uint64        `json:"end_point.warning.count"`
	EndPointWarningDurationAvg   map[string]time.Duration `json:"end_point.warning.duration.avg"`
	EndPointWarningDurationMax   map[string]time.Duration `json:"end_point.warning.duration.max"`
	EndPointWarningDurationMin   map[string]time.Duration `json:"end_point.warning.duration.min"`
	EndPointWarningDurationTotal map[string]time.Duration `json:"end_point.warning.duration.total"`
	EndPointWarningLastCall      map[string]time.Time     `json:"end_point.warning.last_call"`
	EndPointErrorCount           map[string]uint64        `json:"end_point.error.count"`
	EndPointErrorDurationAvg     map[string]time.Duration `json:"end_point.error.duration.avg"`
	EndPointErrorDurationMax     map[string]time.Duration `json:"end_point.error.duration.max"`
	EndPointErrorDurationMin     map[string]time.Duration `json:"end_point.error.duration.min"`
	EndPointErrorDurationTotal   map[string]time.Duration `json:"end_point.error.duration.total"`
	EndPointErrorLastCall        map[string]time.Time     `json:"end_point.error.last_call"`
	EndPointTotalCount           map[string]uint64        `json:"end_point.total.count"`
	EndPointTotalDurationAvg     map[string]time.Duration `json:"end_point.total.duration.avg"`
	EndPointTotalDurationMax     map[string]time.Duration `json:"end_point.total.duration.max"`
	EndPointTotalDurationMin     map[string]time.Duration `json:"end_point.total.duration.min"`
	EndPointTotalDurationTotal   map[string]time.Duration `json:"end_point.total.duration.total"`
	EndPointTotalLastCall        map[string]time.Time     `json:"end_point.total.last_call"`
}

func (m *Metric) Init() {
	m.StartAt = time.Now()

	m.LevelCount = map[string]uint64{}
	m.LevelDurationAvg = map[string]time.Duration{}
	m.LevelDurationMax = map[string]time.Duration{}
	m.LevelDurationMin = map[string]time.Duration{}
	m.LevelDurationTotal = map[string]time.Duration{}
	m.LevelLastCall = map[string]time.Time{}

	m.StatusCodeCount = map[int]uint64{}
	m.StatusCodeDurationAvg = map[int]time.Duration{}
	m.StatusCodeDurationMax = map[int]time.Duration{}
	m.StatusCodeDurationMin = map[int]time.Duration{}
	m.StatusCodeDurationTotal = map[int]time.Duration{}
	m.StatusCodeLastCall = map[int]time.Time{}

	m.MethodCount = map[string]uint64{}
	m.MethodDurationAvg = map[string]time.Duration{}
	m.MethodDurationMax = map[string]time.Duration{}
	m.MethodDurationMin = map[string]time.Duration{}
	m.MethodDurationTotal = map[string]time.Duration{}
	m.MethodLastCall = map[string]time.Time{}

	m.PathPrefixSuccessCount = map[string]uint64{}
	m.PathPrefixSuccessDurationAvg = map[string]time.Duration{}
	m.PathPrefixSuccessDurationMax = map[string]time.Duration{}
	m.PathPrefixSuccessDurationMin = map[string]time.Duration{}
	m.PathPrefixSuccessDurationTotal = map[string]time.Duration{}
	m.PathPrefixSuccessLastCall = map[string]time.Time{}
	m.PathPrefixWarningCount = map[string]uint64{}
	m.PathPrefixWarningDurationAvg = map[string]time.Duration{}
	m.PathPrefixWarningDurationMax = map[string]time.Duration{}
	m.PathPrefixWarningDurationMin = map[string]time.Duration{}
	m.PathPrefixWarningDurationTotal = map[string]time.Duration{}
	m.PathPrefixWarningLastCall = map[string]time.Time{}
	m.PathPrefixErrorCount = map[string]uint64{}
	m.PathPrefixErrorDurationAvg = map[string]time.Duration{}
	m.PathPrefixErrorDurationMax = map[string]time.Duration{}
	m.PathPrefixErrorDurationMin = map[string]time.Duration{}
	m.PathPrefixErrorDurationTotal = map[string]time.Duration{}
	m.PathPrefixErrorLastCall = map[string]time.Time{}
	m.PathPrefixTotalCount = map[string]uint64{}
	m.PathPrefixTotalDurationAvg = map[string]time.Duration{}
	m.PathPrefixTotalDurationMax = map[string]time.Duration{}
	m.PathPrefixTotalDurationMin = map[string]time.Duration{}
	m.PathPrefixTotalDurationTotal = map[string]time.Duration{}
	m.PathPrefixTotalLastCall = map[string]time.Time{}

	m.BaseURLSuccessCount = map[string]uint64{}
	m.BaseURLSuccessDurationAvg = map[string]time.Duration{}
	m.BaseURLSuccessDurationMax = map[string]time.Duration{}
	m.BaseURLSuccessDurationMin = map[string]time.Duration{}
	m.BaseURLSuccessDurationTotal = map[string]time.Duration{}
	m.BaseURLSuccessLastCall = map[string]time.Time{}
	m.BaseURLWarningCount = map[string]uint64{}
	m.BaseURLWarningDurationAvg = map[string]time.Duration{}
	m.BaseURLWarningDurationMax = map[string]time.Duration{}
	m.BaseURLWarningDurationMin = map[string]time.Duration{}
	m.BaseURLWarningDurationTotal = map[string]time.Duration{}
	m.BaseURLWarningLastCall = map[string]time.Time{}
	m.BaseURLErrorCount = map[string]uint64{}
	m.BaseURLErrorDurationAvg = map[string]time.Duration{}
	m.BaseURLErrorDurationMax = map[string]time.Duration{}
	m.BaseURLErrorDurationMin = map[string]time.Duration{}
	m.BaseURLErrorDurationTotal = map[string]time.Duration{}
	m.BaseURLErrorLastCall = map[string]time.Time{}
	m.BaseURLTotalCount = map[string]uint64{}
	m.BaseURLTotalDurationAvg = map[string]time.Duration{}
	m.BaseURLTotalDurationMax = map[string]time.Duration{}
	m.BaseURLTotalDurationMin = map[string]time.Duration{}
	m.BaseURLTotalDurationTotal = map[string]time.Duration{}
	m.BaseURLTotalLastCall = map[string]time.Time{}

	m.EndPointSuccessCount = map[string]uint64{}
	m.EndPointSuccessDurationAvg = map[string]time.Duration{}
	m.EndPointSuccessDurationMax = map[string]time.Duration{}
	m.EndPointSuccessDurationMin = map[string]time.Duration{}
	m.EndPointSuccessDurationTotal = map[string]time.Duration{}
	m.EndPointSuccessLastCall = map[string]time.Time{}
	m.EndPointWarningCount = map[string]uint64{}
	m.EndPointWarningDurationAvg = map[string]time.Duration{}
	m.EndPointWarningDurationMax = map[string]time.Duration{}
	m.EndPointWarningDurationMin = map[string]time.Duration{}
	m.EndPointWarningDurationTotal = map[string]time.Duration{}
	m.EndPointWarningLastCall = map[string]time.Time{}
	m.EndPointErrorCount = map[string]uint64{}
	m.EndPointErrorDurationAvg = map[string]time.Duration{}
	m.EndPointErrorDurationMax = map[string]time.Duration{}
	m.EndPointErrorDurationMin = map[string]time.Duration{}
	m.EndPointErrorDurationTotal = map[string]time.Duration{}
	m.EndPointErrorLastCall = map[string]time.Time{}
	m.EndPointTotalCount = map[string]uint64{}
	m.EndPointTotalDurationAvg = map[string]time.Duration{}
	m.EndPointTotalDurationMax = map[string]time.Duration{}
	m.EndPointTotalDurationMin = map[string]time.Duration{}
	m.EndPointTotalDurationTotal = map[string]time.Duration{}
	m.EndPointTotalLastCall = map[string]time.Time{}

	go m.serve()
}

func (m *Metric) serve() {
	port := fmt.Sprintf("%v", Conf.Metric.Port)
	fmt.Println()
	fmt.Println("Metric available at " + Fmt("http://localhost:"+port, Magenta))
	http.HandleFunc("/", m.serveUI)
	http.HandleFunc("/ws", m.serveWs)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Metric", Fmt(err.Error(), Red))
	}
}

func (m *Metric) serveUI(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}
	fileContent, err := Conf.MetricUI.ReadFile("ui" + filePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", http.DetectContentType(fileContent))
	w.Write(fileContent)
}

func (m *Metric) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	activeConnMu.Lock()
	activeConn[conn] = struct{}{}
	activeConnMu.Unlock()

	defer func() {
		activeConnMu.Lock()
		delete(activeConn, conn)
		activeConnMu.Unlock()

		conn.Close()
	}()
}

func (m *Metric) sendData(data []byte) {
	activeConnMu.Lock()
	defer activeConnMu.Unlock()

	for conn := range activeConn {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (m *Metric) Update(c *Ctx) {
	level := "error"
	if c.Err == nil {
		if c.StatusCode >= http.StatusOK && c.StatusCode < http.StatusBadRequest {
			level = "success"
		} else if c.StatusCode >= http.StatusBadRequest && c.StatusCode < http.StatusInternalServerError {
			level = "warning"
		}
	}
	m.LevelLastCall[level] = c.FinishAt
	if _, ok := m.LevelCount[level]; ok {
		m.LevelCount[level] += 1
	} else {
		m.LevelCount[level] = 1
	}
	if _, ok := m.LevelDurationTotal[level]; ok {
		m.LevelDurationTotal[level] += c.Duration
	} else {
		m.LevelDurationTotal[level] = c.Duration
	}
	m.LevelDurationAvg[level] = m.LevelDurationTotal[level] / time.Duration(m.LevelCount[level])
	if ldm, ok := m.LevelDurationMax[level]; !ok || c.Duration > ldm {
		m.LevelDurationMax[level] = c.Duration
	}
	if ldm, ok := m.LevelDurationMin[level]; !ok || c.Duration < ldm {
		m.LevelDurationMin[level] = c.Duration
	}

	m.StatusCodeLastCall[c.StatusCode] = c.FinishAt
	if _, ok := m.StatusCodeCount[c.StatusCode]; ok {
		m.StatusCodeCount[c.StatusCode] += 1
	} else {
		m.StatusCodeCount[c.StatusCode] = 1
	}
	if _, ok := m.StatusCodeDurationTotal[c.StatusCode]; ok {
		m.StatusCodeDurationTotal[c.StatusCode] += c.Duration
	} else {
		m.StatusCodeDurationTotal[c.StatusCode] = c.Duration
	}
	m.StatusCodeDurationAvg[c.StatusCode] = m.StatusCodeDurationTotal[c.StatusCode] / time.Duration(m.StatusCodeCount[c.StatusCode])
	if ldm, ok := m.StatusCodeDurationMax[c.StatusCode]; !ok || c.Duration > ldm {
		m.StatusCodeDurationMax[c.StatusCode] = c.Duration
	}
	if ldm, ok := m.StatusCodeDurationMin[c.StatusCode]; !ok || c.Duration < ldm {
		m.StatusCodeDurationMin[c.StatusCode] = c.Duration
	}

	m.MethodLastCall[c.Method] = c.FinishAt
	if _, ok := m.MethodCount[c.Method]; ok {
		m.MethodCount[c.Method] += 1
	} else {
		m.MethodCount[c.Method] = 1
	}
	if _, ok := m.MethodDurationTotal[c.Method]; ok {
		m.MethodDurationTotal[c.Method] += c.Duration
	} else {
		m.MethodDurationTotal[c.Method] = c.Duration
	}
	m.MethodDurationAvg[c.Method] = m.MethodDurationTotal[c.Method] / time.Duration(m.MethodCount[c.Method])
	if ldm, ok := m.MethodDurationMax[c.Method]; !ok || c.Duration > ldm {
		m.MethodDurationMax[c.Method] = c.Duration
	}
	if ldm, ok := m.MethodDurationMin[c.Method]; !ok || c.Duration < ldm {
		m.MethodDurationMin[c.Method] = c.Duration
	}

	if level == "success" {
		m.PathPrefixSuccessLastCall[c.PathPrefix] = c.FinishAt
		if _, ok := m.PathPrefixSuccessCount[c.PathPrefix]; ok {
			m.PathPrefixSuccessCount[c.PathPrefix] += 1
		} else {
			m.PathPrefixSuccessCount[c.PathPrefix] = 1
		}
		if _, ok := m.PathPrefixSuccessDurationTotal[c.PathPrefix]; ok {
			m.PathPrefixSuccessDurationTotal[c.PathPrefix] += c.Duration
		} else {
			m.PathPrefixSuccessDurationTotal[c.PathPrefix] = c.Duration
		}
		m.PathPrefixSuccessDurationAvg[c.PathPrefix] = m.PathPrefixSuccessDurationTotal[c.PathPrefix] / time.Duration(m.PathPrefixSuccessCount[c.PathPrefix])
		if ldm, ok := m.PathPrefixSuccessDurationMax[c.PathPrefix]; !ok || c.Duration > ldm {
			m.PathPrefixSuccessDurationMax[c.PathPrefix] = c.Duration
		}
		if ldm, ok := m.PathPrefixSuccessDurationMin[c.PathPrefix]; !ok || c.Duration < ldm {
			m.PathPrefixSuccessDurationMin[c.PathPrefix] = c.Duration
		}
	} else if level == "warning" {
		m.PathPrefixWarningLastCall[c.PathPrefix] = c.FinishAt
		if _, ok := m.PathPrefixWarningCount[c.PathPrefix]; ok {
			m.PathPrefixWarningCount[c.PathPrefix] += 1
		} else {
			m.PathPrefixWarningCount[c.PathPrefix] = 1
		}
		if _, ok := m.PathPrefixWarningDurationTotal[c.PathPrefix]; ok {
			m.PathPrefixWarningDurationTotal[c.PathPrefix] += c.Duration
		} else {
			m.PathPrefixWarningDurationTotal[c.PathPrefix] = c.Duration
		}
		m.PathPrefixWarningDurationAvg[c.PathPrefix] = m.PathPrefixWarningDurationTotal[c.PathPrefix] / time.Duration(m.PathPrefixWarningCount[c.PathPrefix])
		if ldm, ok := m.PathPrefixWarningDurationMax[c.PathPrefix]; !ok || c.Duration > ldm {
			m.PathPrefixWarningDurationMax[c.PathPrefix] = c.Duration
		}
		if ldm, ok := m.PathPrefixWarningDurationMin[c.PathPrefix]; !ok || c.Duration < ldm {
			m.PathPrefixWarningDurationMin[c.PathPrefix] = c.Duration
		}
	} else {
		m.PathPrefixErrorLastCall[c.PathPrefix] = c.FinishAt
		if _, ok := m.PathPrefixErrorCount[c.PathPrefix]; ok {
			m.PathPrefixErrorCount[c.PathPrefix] += 1
		} else {
			m.PathPrefixErrorCount[c.PathPrefix] = 1
		}
		if _, ok := m.PathPrefixErrorDurationTotal[c.PathPrefix]; ok {
			m.PathPrefixErrorDurationTotal[c.PathPrefix] += c.Duration
		} else {
			m.PathPrefixErrorDurationTotal[c.PathPrefix] = c.Duration
		}
		m.PathPrefixErrorDurationAvg[c.PathPrefix] = m.PathPrefixErrorDurationTotal[c.PathPrefix] / time.Duration(m.PathPrefixErrorCount[c.PathPrefix])
		if ldm, ok := m.PathPrefixErrorDurationMax[c.PathPrefix]; !ok || c.Duration > ldm {
			m.PathPrefixErrorDurationMax[c.PathPrefix] = c.Duration
		}
		if ldm, ok := m.PathPrefixErrorDurationMin[c.PathPrefix]; !ok || c.Duration < ldm {
			m.PathPrefixErrorDurationMin[c.PathPrefix] = c.Duration
		}
	}
	m.PathPrefixTotalLastCall[c.PathPrefix] = c.FinishAt
	if _, ok := m.PathPrefixTotalCount[c.PathPrefix]; ok {
		m.PathPrefixTotalCount[c.PathPrefix] += 1
	} else {
		m.PathPrefixTotalCount[c.PathPrefix] = 1
	}
	if _, ok := m.PathPrefixTotalDurationTotal[c.PathPrefix]; ok {
		m.PathPrefixTotalDurationTotal[c.PathPrefix] += c.Duration
	} else {
		m.PathPrefixTotalDurationTotal[c.PathPrefix] = c.Duration
	}
	m.PathPrefixTotalDurationAvg[c.PathPrefix] = m.PathPrefixTotalDurationTotal[c.PathPrefix] / time.Duration(m.PathPrefixTotalCount[c.PathPrefix])
	if ldm, ok := m.PathPrefixTotalDurationMax[c.PathPrefix]; !ok || c.Duration > ldm {
		m.PathPrefixTotalDurationMax[c.PathPrefix] = c.Duration
	}
	if ldm, ok := m.PathPrefixTotalDurationMin[c.PathPrefix]; !ok || c.Duration < ldm {
		m.PathPrefixTotalDurationMin[c.PathPrefix] = c.Duration
	}

	if level == "success" {
		m.BaseURLSuccessLastCall[c.BaseURL] = c.FinishAt
		if _, ok := m.BaseURLSuccessCount[c.BaseURL]; ok {
			m.BaseURLSuccessCount[c.BaseURL] += 1
		} else {
			m.BaseURLSuccessCount[c.BaseURL] = 1
		}
		if _, ok := m.BaseURLSuccessDurationTotal[c.BaseURL]; ok {
			m.BaseURLSuccessDurationTotal[c.BaseURL] += c.Duration
		} else {
			m.BaseURLSuccessDurationTotal[c.BaseURL] = c.Duration
		}
		m.BaseURLSuccessDurationAvg[c.BaseURL] = m.BaseURLSuccessDurationTotal[c.BaseURL] / time.Duration(m.BaseURLSuccessCount[c.BaseURL])
		if ldm, ok := m.BaseURLSuccessDurationMax[c.BaseURL]; !ok || c.Duration > ldm {
			m.BaseURLSuccessDurationMax[c.BaseURL] = c.Duration
		}
		if ldm, ok := m.BaseURLSuccessDurationMin[c.BaseURL]; !ok || c.Duration < ldm {
			m.BaseURLSuccessDurationMin[c.BaseURL] = c.Duration
		}
	} else if level == "warning" {
		m.BaseURLWarningLastCall[c.BaseURL] = c.FinishAt
		if _, ok := m.BaseURLWarningCount[c.BaseURL]; ok {
			m.BaseURLWarningCount[c.BaseURL] += 1
		} else {
			m.BaseURLWarningCount[c.BaseURL] = 1
		}
		if _, ok := m.BaseURLWarningDurationTotal[c.BaseURL]; ok {
			m.BaseURLWarningDurationTotal[c.BaseURL] += c.Duration
		} else {
			m.BaseURLWarningDurationTotal[c.BaseURL] = c.Duration
		}
		m.BaseURLWarningDurationAvg[c.BaseURL] = m.BaseURLWarningDurationTotal[c.BaseURL] / time.Duration(m.BaseURLWarningCount[c.BaseURL])
		if ldm, ok := m.BaseURLWarningDurationMax[c.BaseURL]; !ok || c.Duration > ldm {
			m.BaseURLWarningDurationMax[c.BaseURL] = c.Duration
		}
		if ldm, ok := m.BaseURLWarningDurationMin[c.BaseURL]; !ok || c.Duration < ldm {
			m.BaseURLWarningDurationMin[c.BaseURL] = c.Duration
		}
	} else {
		m.BaseURLErrorLastCall[c.BaseURL] = c.FinishAt
		if _, ok := m.BaseURLErrorCount[c.BaseURL]; ok {
			m.BaseURLErrorCount[c.BaseURL] += 1
		} else {
			m.BaseURLErrorCount[c.BaseURL] = 1
		}
		if _, ok := m.BaseURLErrorDurationTotal[c.BaseURL]; ok {
			m.BaseURLErrorDurationTotal[c.BaseURL] += c.Duration
		} else {
			m.BaseURLErrorDurationTotal[c.BaseURL] = c.Duration
		}
		m.BaseURLErrorDurationAvg[c.BaseURL] = m.BaseURLErrorDurationTotal[c.BaseURL] / time.Duration(m.BaseURLErrorCount[c.BaseURL])
		if ldm, ok := m.BaseURLErrorDurationMax[c.BaseURL]; !ok || c.Duration > ldm {
			m.BaseURLErrorDurationMax[c.BaseURL] = c.Duration
		}
		if ldm, ok := m.BaseURLErrorDurationMin[c.BaseURL]; !ok || c.Duration < ldm {
			m.BaseURLErrorDurationMin[c.BaseURL] = c.Duration
		}
	}
	m.BaseURLTotalLastCall[c.BaseURL] = c.FinishAt
	if _, ok := m.BaseURLTotalCount[c.BaseURL]; ok {
		m.BaseURLTotalCount[c.BaseURL] += 1
	} else {
		m.BaseURLTotalCount[c.BaseURL] = 1
	}
	if _, ok := m.BaseURLTotalDurationTotal[c.BaseURL]; ok {
		m.BaseURLTotalDurationTotal[c.BaseURL] += c.Duration
	} else {
		m.BaseURLTotalDurationTotal[c.BaseURL] = c.Duration
	}
	m.BaseURLTotalDurationAvg[c.BaseURL] = m.BaseURLTotalDurationTotal[c.BaseURL] / time.Duration(m.BaseURLTotalCount[c.BaseURL])
	if ldm, ok := m.BaseURLTotalDurationMax[c.BaseURL]; !ok || c.Duration > ldm {
		m.BaseURLTotalDurationMax[c.BaseURL] = c.Duration
	}
	if ldm, ok := m.BaseURLTotalDurationMin[c.BaseURL]; !ok || c.Duration < ldm {
		m.BaseURLTotalDurationMin[c.BaseURL] = c.Duration
	}

	if level == "success" {
		m.EndPointSuccessLastCall[c.EndPoint+" "+c.Method] = c.FinishAt
		if _, ok := m.EndPointSuccessCount[c.EndPoint+" "+c.Method]; ok {
			m.EndPointSuccessCount[c.EndPoint+" "+c.Method] += 1
		} else {
			m.EndPointSuccessCount[c.EndPoint+" "+c.Method] = 1
		}
		if _, ok := m.EndPointSuccessDurationTotal[c.EndPoint+" "+c.Method]; ok {
			m.EndPointSuccessDurationTotal[c.EndPoint+" "+c.Method] += c.Duration
		} else {
			m.EndPointSuccessDurationTotal[c.EndPoint+" "+c.Method] = c.Duration
		}
		m.EndPointSuccessDurationAvg[c.EndPoint+" "+c.Method] = m.EndPointSuccessDurationTotal[c.EndPoint+" "+c.Method] / time.Duration(m.EndPointSuccessCount[c.EndPoint+" "+c.Method])
		if ldm, ok := m.EndPointSuccessDurationMax[c.EndPoint+" "+c.Method]; !ok || c.Duration > ldm {
			m.EndPointSuccessDurationMax[c.EndPoint+" "+c.Method] = c.Duration
		}
		if ldm, ok := m.EndPointSuccessDurationMin[c.EndPoint+" "+c.Method]; !ok || c.Duration < ldm {
			m.EndPointSuccessDurationMin[c.EndPoint+" "+c.Method] = c.Duration
		}
	} else if level == "warning" {
		m.EndPointWarningLastCall[c.EndPoint+" "+c.Method] = c.FinishAt
		if _, ok := m.EndPointWarningCount[c.EndPoint+" "+c.Method]; ok {
			m.EndPointWarningCount[c.EndPoint+" "+c.Method] += 1
		} else {
			m.EndPointWarningCount[c.EndPoint+" "+c.Method] = 1
		}
		if _, ok := m.EndPointWarningDurationTotal[c.EndPoint+" "+c.Method]; ok {
			m.EndPointWarningDurationTotal[c.EndPoint+" "+c.Method] += c.Duration
		} else {
			m.EndPointWarningDurationTotal[c.EndPoint+" "+c.Method] = c.Duration
		}
		m.EndPointWarningDurationAvg[c.EndPoint+" "+c.Method] = m.EndPointWarningDurationTotal[c.EndPoint+" "+c.Method] / time.Duration(m.EndPointWarningCount[c.EndPoint+" "+c.Method])
		if ldm, ok := m.EndPointWarningDurationMax[c.EndPoint+" "+c.Method]; !ok || c.Duration > ldm {
			m.EndPointWarningDurationMax[c.EndPoint+" "+c.Method] = c.Duration
		}
		if ldm, ok := m.EndPointWarningDurationMin[c.EndPoint+" "+c.Method]; !ok || c.Duration < ldm {
			m.EndPointWarningDurationMin[c.EndPoint+" "+c.Method] = c.Duration
		}
	} else {
		m.EndPointErrorLastCall[c.EndPoint+" "+c.Method] = c.FinishAt
		if _, ok := m.EndPointErrorCount[c.EndPoint+" "+c.Method]; ok {
			m.EndPointErrorCount[c.EndPoint+" "+c.Method] += 1
		} else {
			m.EndPointErrorCount[c.EndPoint+" "+c.Method] = 1
		}
		if _, ok := m.EndPointErrorDurationTotal[c.EndPoint+" "+c.Method]; ok {
			m.EndPointErrorDurationTotal[c.EndPoint+" "+c.Method] += c.Duration
		} else {
			m.EndPointErrorDurationTotal[c.EndPoint+" "+c.Method] = c.Duration
		}
		m.EndPointErrorDurationAvg[c.EndPoint+" "+c.Method] = m.EndPointErrorDurationTotal[c.EndPoint+" "+c.Method] / time.Duration(m.EndPointErrorCount[c.EndPoint+" "+c.Method])
		if ldm, ok := m.EndPointErrorDurationMax[c.EndPoint+" "+c.Method]; !ok || c.Duration > ldm {
			m.EndPointErrorDurationMax[c.EndPoint+" "+c.Method] = c.Duration
		}
		if ldm, ok := m.EndPointErrorDurationMin[c.EndPoint+" "+c.Method]; !ok || c.Duration < ldm {
			m.EndPointErrorDurationMin[c.EndPoint+" "+c.Method] = c.Duration
		}
	}
	m.EndPointTotalLastCall[c.EndPoint+" "+c.Method] = c.FinishAt
	if _, ok := m.EndPointTotalCount[c.EndPoint+" "+c.Method]; ok {
		m.EndPointTotalCount[c.EndPoint+" "+c.Method] += 1
	} else {
		m.EndPointTotalCount[c.EndPoint+" "+c.Method] = 1
	}
	if _, ok := m.EndPointTotalDurationTotal[c.EndPoint+" "+c.Method]; ok {
		m.EndPointTotalDurationTotal[c.EndPoint+" "+c.Method] += c.Duration
	} else {
		m.EndPointTotalDurationTotal[c.EndPoint+" "+c.Method] = c.Duration
	}
	m.EndPointTotalDurationAvg[c.EndPoint+" "+c.Method] = m.EndPointTotalDurationTotal[c.EndPoint+" "+c.Method] / time.Duration(m.EndPointTotalCount[c.EndPoint+" "+c.Method])
	if ldm, ok := m.EndPointTotalDurationMax[c.EndPoint+" "+c.Method]; !ok || c.Duration > ldm {
		m.EndPointTotalDurationMax[c.EndPoint+" "+c.Method] = c.Duration
	}
	if ldm, ok := m.EndPointTotalDurationMin[c.EndPoint+" "+c.Method]; !ok || c.Duration < ldm {
		m.EndPointTotalDurationMin[c.EndPoint+" "+c.Method] = c.Duration
	}
	b, err := json.Marshal(m)
	if err != nil {
		m.sendData(b)
	}
}
