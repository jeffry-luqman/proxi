package main

import (
	"embed"
	"os"

	"github.com/jeffry-luqman/proxi/app"
	"github.com/spf13/cobra"
)

//go:embed all:ui
var ui embed.FS

func main() {
	app.Conf.MetricUI = ui

	cmd := &cobra.Command{
		Use:   "proxi",
		Short: "Simple reverse proxy.",
		Long:  `Proxi is a simple reverse proxy, allows you to forward HTTP requests from multiple endpoints to different targets based on the provided path.`,
		Run: func(c *cobra.Command, args []string) {
			app.Run()
		},
	}
	cmd.Flags().StringVarP(&app.ConfigFile, "config", "c", "", "Configuration file name, (default proxi.yml)\nSample config file: https://raw.githubusercontent.com/jeffry-luqman/proxi/main/proxi.yml")
	cmd.Flags().IntVarP(&app.Conf.Port, "port", "p", app.Conf.Port, "Port")
	cmd.Flags().StringVarP(&app.Conf.TargetStr, "targets", "t", "", "Target URL for each prefix, delimited with semicolon.\nEx: proxi -t \"/ https://example.com; /api https://api.example.com\"")
	cmd.Flags().BoolVarP(&app.Conf.UseStdlib, "use-stdlib", "", app.Conf.UseStdlib, "Use net/http instead of fasthttp")
	cmd.Flags().BoolVarP(&app.Conf.Log.Console.Disable, "quiet", "q", app.Conf.Log.Console.Disable, "Silence output on the terminal")
	cmd.Flags().BoolVarP(&app.Conf.Log.Console.PrintRequestImmediately, "debug", "d", app.Conf.Log.Console.PrintRequestImmediately, "Print a request log to the terminal without waiting for a response")
	cmd.Flags().BoolVarP(&app.Conf.PrintFullURL, "print-full-url", "f", app.Conf.PrintFullURL, "Print an full url instead of path only")
	cmd.Flags().StringVarP(&app.Conf.Log.File.Filename, "log", "l", "", "Specify log file")
	cmd.Flags().IntVarP(&app.Conf.Metric.Port, "metric", "m", app.Conf.Metric.Port, "Specify metric port")
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
