package modes

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/padi2312/compose-check-updates/internal"
)

func Interactive(updateInfos []internal.UpdateInfo) {
	reader := bufio.NewReader(os.Stdin)
	// group update infos by compose file path
	groupedUpdateInfos := make(map[string][]internal.UpdateInfo)

	for _, i := range updateInfos {
		groupedUpdateInfos[i.FilePath] = append(groupedUpdateInfos[i.FilePath], i)
	}

	for path, infos := range groupedUpdateInfos {
		slog.Info(fmt.Sprintf("File: %s", path))
		for _, i := range infos {
			if i.HasNewVersion() {
				// Ask user if they want to update the file with y/n
				slog.Info(fmt.Sprintf("New version for %s: current=%s, latest=%s", i.ImageName, i.CurrentTag, i.LatestTag))
				fmt.Printf("Do you want to update %s: %s -> %s (y/n): ", i.ImageName, i.CurrentTag, i.LatestTag)
				text, _ := reader.ReadString('\n')
				text = strings.TrimSpace(text)
				if text == "y" {
					if err := i.Update(); err != nil {
						fmt.Printf("Error updating file: %v\n", err)
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
							fmt.Printf("Error restarting service: %v\n", err)
							continue
						}
						fmt.Printf("Service restarted: image=%s\n", i.FullImageName)
					}
				}
				fmt.Println()
			}
		}

	}
}
