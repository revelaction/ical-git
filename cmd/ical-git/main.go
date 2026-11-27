package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/caasmo/go-daemon-runner/run"
	"github.com/revelaction/ical-git/ical"
)

const usage = `Usage:
     ical-git [-c CONF_FILE] 

Options:
     -c, --config                load the configuration file at CONF_FILE instead of default
     -v, --version               Print the version 
     -h, --help                  Show this

CONF_FILE is the toml configuration file 

ical-git will react to a SIGHUP signal reloading the configuration file.

Examples:
     $ ical-git --config /path/to/config/file.toml # start the daemon with the configuration file
     $ ical-git -v  # print version`

// Version can be set at link time
var BuildTag string

// configFile is the config file path
const configPathDefault = "icalgit.toml"

func main() {

	var configPath string
	var versionFlag, helpFlag bool

	flag.Usage = func() { fmt.Fprintf(os.Stderr, "%s\n", usage) }
	flag.BoolVar(&versionFlag, "v", false, "print the version")
	flag.BoolVar(&versionFlag, "version", false, "print the version")
	flag.StringVar(&configPath, "c", configPathDefault, "the config file")
	flag.StringVar(&configPath, "config", configPathDefault, "the config file")
	flag.BoolVar(&helpFlag, "h", false, "print the version")
	flag.BoolVar(&helpFlag, "help", false, "Show this")
	flag.Parse()

	if versionFlag {
		if BuildTag != "" {
			fmt.Println(BuildTag)
			return
		}
		fmt.Println("(unknown)")
		return
	}

	if helpFlag {
		flag.Usage()
		return
	}

	initializeLogger()
	slog.Info("🏁 app:", "Version", BuildTag)

	d := ical.NewDaemon(configPath)

	r, err := run.NewRunner()
	if err != nil {
		slog.Error("Failed to create runner", "error", err)
		os.Exit(1)
	}
	r.Add(d)

	err = r.Run()
	if err != nil {
		slog.Error("Daemon run failed", "error", err)
		os.Exit(1)
	}
}

func initializeLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))

	slog.SetDefault(logger)
}