package config

import (
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const DebugIndividualOps log.Level = -8

func NewLogger() *log.Logger {
	var level log.Level

	levelEnv, levelEnvSet := os.LookupEnv("LOG_LEVEL")
	levelEnv = strings.ToUpper(levelEnv)

	goEnv := strings.ToUpper(os.Getenv("GOENV"))
	devMode := goEnv == "DEVELOPMENT"

	if devMode && !levelEnvSet {
		level = log.DebugLevel
	} else {
		switch levelEnv {
		case "DEBUG_PATCHES":
			level = DebugIndividualOps
		case "DEBUG":
			level = log.DebugLevel
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
		ReportCaller: devMode,
	}

	logger := log.NewWithOptions(os.Stderr, opts)

	s := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Foreground(lipgloss.NoColor{})

	logger.SetStyles(&log.Styles{
		Timestamp: lipgloss.NewStyle(),
		Levels: map[log.Level]lipgloss.Style{
			DebugIndividualOps: s.
				SetString("PTCH").
				Background(lipgloss.Color("199")),
			log.DebugLevel: s.
				SetString("DEBU").
				Background(lipgloss.Color("33")),
			log.InfoLevel: s.
				SetString("INFO").
				Background(lipgloss.Color("72")),
			log.WarnLevel: s.
				SetString("WARN").
				Background(lipgloss.Color("202")),
			log.ErrorLevel: s.
				SetString("ERRO").
				Background(lipgloss.Color("196")),
			log.FatalLevel: s.
				SetString("FATA").
				Background(lipgloss.Color("162")),
		},
	})

	return logger
}
