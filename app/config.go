package app

var config = Config{
	Port:    8181,
	Targets: map[string]string{"/": "http://localhost:3000"},
	Log: LogConfig{
		Console: LogConsoleConfig{
			Enable:                  true,
			PrintRequestImmediately: false,
		},
		File: LogFileConfig{
			Enable:                 true,
			Filename:               "proxi.log",
			UseLocalTime:           true,
			MaxSize:                100,
			MaxAge:                 7,
			MaxBackups:             0,
			IncludeRequestHeaders:  false,
			IncludeRequestBody:     false,
			IncludeResponseHeaders: false,
			IncludeResponseBody:    false,
		},
	},
}

type Config struct {
	Port      int               `yaml:"port"`
	UseStdlib bool              `yaml:"us_stdlib"`
	Targets   map[string]string `yaml:"targets"`
	Log       LogConfig         `yaml:"log"`
	Metric    MetricConfig      `yaml:"metric"`
}

type LogConfig struct {
	Console LogConsoleConfig `yaml:"console"`
	File    LogFileConfig    `yaml:"file"`
}

type LogConsoleConfig struct {
	Enable                  bool `yaml:"enable"`
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
