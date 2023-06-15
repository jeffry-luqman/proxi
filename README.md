# Proxi
Proxi is a simple reverse proxy, allows you to forward HTTP requests from multiple endpoints to different targets based on the provided path.

## Features
- Simple HTTP reverse proxy
- Configuration using YAML file
- Request and response logging
- Realtime metric

## Installation

To install proxi, you can download a prebuilt binary, or you can use go install if you have [Go](https://go.dev) installed.

### Windows
TODO

### Mac
TODO

### Linux
TODO

### With [Go Install](https://go.dev/ref/mod#go-install)
```
go install github.com/jeffry-luqman/proxi@latest
```

## Usage
You can also see this information by running `proxi --help` from the command line.
```
Usage:
  proxi [flags]

Flags:
  -h, --help             help for proxi
  -c, --config string    Configuration file name, (default proxi.yml)
                         Sample config file: https://raw.githubusercontent.com/jeffry-luqman/proxi/main/proxi.yml
  -p, --port int         Port (default 8181)
  -t, --targets string   Target URL for each prefix, delimited with semicolon.
                         Ex: proxi -t "/ https://example.com; /api https://api.example.com"
  -d, --debug            Print a request log to the terminal without waiting for a response
  -q, --quiet            Silence output on the terminal
  -l, --log string       Specify log file
  -m, --metric int       Specify metric port
      --use-stdlib       Use net/http instead of fasthttp
```

## License
Proxi is free and open-source software licensed under the [MIT License](https://github.com/jeffry-luqman/proxi/blob/main/LICENSE).