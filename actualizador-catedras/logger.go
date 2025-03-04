package main

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func newLogger() {
	var level log.Level

	levelEnv, levelEnvSet := os.LookupEnv("LOG_LEVEL")
	goEnv := strings.ToUpper(os.Getenv("GO_ENV"))

	if levelEnv := strings.ToUpper(levelEnv); (goEnv == "DEVELOPMENT" && !levelEnvSet) ||
		levelEnv == "DEBUG" {
		level = log.DebugLevel
	} else {
		switch levelEnv {
		case "WARN", "WARNING":
			level = log.WarnLevel
		case "ERROR":
			level = log.ErrorLevel
		default:
			level = log.InfoLevel
		}
	}

	opts := log.Options{
		Level:        level,
		ReportCaller: level == log.DebugLevel,
	}

	logger := log.NewWithOptions(os.Stderr, opts)

	commonStyles := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(lipgloss.NoColor{})

	logger.SetStyles(&log.Styles{
		Timestamp: lipgloss.NewStyle(),
		Levels: map[log.Level]lipgloss.Style{
			log.DebugLevel: commonStyles.
				SetString("DEBU").
				Background(lipgloss.Color("33")),
			log.InfoLevel: commonStyles.
				SetString("INFO").
				Background(lipgloss.Color("72")),
			log.WarnLevel: commonStyles.
				SetString("WARN").
				Background(lipgloss.Color("202")),
			log.ErrorLevel: commonStyles.
				SetString("ERRO").
				Background(lipgloss.Color("196")),
			log.FatalLevel: commonStyles.
				SetString("FATA").
				Background(lipgloss.Color("162")),
		},
	})

	log.SetDefault(logger)
}
