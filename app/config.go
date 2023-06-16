package app

import "embed"

var Conf = Config{
	Port:    8181,
	Targets: map[string]string{"/": "http://localhost:3000"},
	Log: LogConfig{
		Console: LogConsoleConfig{
			Enable: true,
		},
		File: LogFileConfig{
			UseLocalTime: true,
			MaxSize:      100,
			MaxAge:       7,
		},
	},
}

type Config struct {
	Port      int               `yaml:"port"`
	UseStdlib bool              `yaml:"use_stdlib"`
	TargetStr string            `yaml:"-"`
	Targets   map[string]string `yaml:"targets"`
	Log       LogConfig         `yaml:"log"`
	Metric    MetricConfig      `yaml:"metric"`
	MetricUI  embed.FS          `yaml:"-"`
}

type LogConfig struct {
	Console LogConsoleConfig `yaml:"console"`
	File    LogFileConfig    `yaml:"file"`
}

type LogConsoleConfig struct {
	Enable                  bool `yaml:"enable"`
	Disable                 bool `yaml:"-"`
	PrintRequestImmediately bool `yaml:"print_request_immediately"`
}

type LogFileConfig struct {
	Enable                 bool   `yaml:"enable"`
	Filename               string `yaml:"filename"`
	UseLocalTime           bool   `yaml:"use_local_time"`
	MaxSize                int    `yaml:"max_size"`
	MaxAge                 int    `yaml:"max_age"`
	MaxBackups             int    `yaml:"max_backups"`
	IncludeRequestHeaders  bool   `yaml:"include_request_headers"`
	IncludeRequestBody     bool   `yaml:"include_request_body"`
	IncludeResponseHeaders bool   `yaml:"include_response_headers"`
	IncludeResponseBody    bool   `yaml:"include_response_body"`
}

type MetricConfig struct {
	Enable bool `yaml:"enable"`
	Port   int  `yaml:"port"`
}
