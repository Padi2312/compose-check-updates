package main

import (
	"log/slog"
	"os"

	"github.com/padi2312/compose-check-updates/internal"
	"github.com/padi2312/compose-check-updates/internal/logger"
	"github.com/padi2312/compose-check-updates/internal/modes"
)

var version = "0.1.0"

func main() {
	// Set colorized logger
	logger := slog.New(logger.NewCustomHandler(slog.LevelInfo, os.Stdout))
	slog.SetDefault(logger)

	ccuFlags := internal.Parse(version)
	root := ccuFlags.Directory
	composeFilePaths, err := internal.GetComposeFilePaths(root)
	if err != nil {
		slog.Error("Error getting compose file paths", "error", err)
		os.Exit(1)
		return
	}

	updateInfos := []internal.UpdateInfo{}
	for _, path := range composeFilePaths {
		updateChecker := internal.NewUpdateChecker(path, internal.NewRegistry(""))
		info, err := updateChecker.Check(ccuFlags.Major, ccuFlags.Minor, ccuFlags.Patch)
		if err != nil {
			slog.Error("Error checking for updates", "error", err)
			continue
		}
		updateInfos = append(updateInfos, info...)
	}

	if ccuFlags.Interactive {
		modes.Interactive(updateInfos)
		return
	} else {
		modes.Default(updateInfos, ccuFlags)
	}

}
