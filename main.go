package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
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

	updateInfos := []UpdateInfo{}
	reader := bufio.NewReader(os.Stdin)
	for _, path := range composeFilePaths {
		updateChecker := NewUpdateChecker(path)
		info, err := updateChecker.Check()
		if err != nil {
			slog.Error("Error checking for updates", "error", err)
			continue
		}
		updateInfos = append(updateInfos, info...)
	}

	for _, i := range updateInfos {
		if i.HasNewVersion() {
			// Ask user if they want to update the file with y/n
			fmt.Printf("New version for %s: current=%s, latest=%s\n", i.FullImageName, i.CurrentTag, i.LatestTag)
			fmt.Printf("Do you want to update the file? (y/n)")
			text, _ := reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if text == "y" {
				if err := i.Update(); err != nil {
					slog.Error("Error updating file", "error", err)
					continue
				}

				fmt.Printf("File updated. Image %s has new version %s\n", i.ImageName, i.LatestTag)
				fmt.Printf("Do you want to restart the service? (y/n)")
				text, _ = reader.ReadString('\n')
				text = strings.TrimSpace(text)
				if text != "y" {
					continue
				} else {
					if err := i.Restart(); err != nil {
						slog.Error("Error restarting service", "error", err)
						continue
					}
					slog.Info("Service restarted", "image", i.FullImageName)
				}
			}

		}
	}
}
