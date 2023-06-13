package main

import (
	"os"

	"github.com/jeffry-luqman/proxi/app"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "proxi",
		Short: "Simple reverse proxy.",
		Long:  `Proxi is a simple reverse proxy, allows you to forward HTTP requests from multiple endpoints to different targets based on the provided path.`,
		Run: func(c *cobra.Command, args []string) {
			app.Run()
		},
	}
	cmd.Flags().StringVarP(&app.ConfigFile, "file", "f", app.ConfigFile, "Proxi configuration file")
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
