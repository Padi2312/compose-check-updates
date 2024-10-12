package main

import (
	"log/slog"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("Please provide the root directory as an argument.")
		os.Exit(1)
		return
	}

	root := os.Args[1]

	composeFilePaths, err := GetComposeFilePaths(root)
	if err != nil {
		slog.Error("Error getting compose file paths", "error", err)
		os.Exit(1)
		return
	}

	for _, path := range composeFilePaths {
		updateChecker := NewUpdateChecker(path)
		info, err := updateChecker.Check()
		if err != nil {
			slog.Error("Error checking for updates", "error", err)
			continue
		}

		for _, i := range info {
			slog.Info("", "image", i.ImageName, "current", i.CurrentTag, "latest", i.LatestTag)
		}
	}

}
