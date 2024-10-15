package modes

import (
	"fmt"
	"log/slog"

	"github.com/padi2312/compose-check-updates/internal"
)

func Default(updateInfos []internal.UpdateInfo, ccuFlags internal.CCUFlags) {
	for _, i := range updateInfos {
		if i.HasNewVersion(ccuFlags.Major, ccuFlags.Minor, ccuFlags.Patch) {
			if !ccuFlags.Update && !ccuFlags.Restart {
				// If no flags are provided, just print the new version
				slog.Info(fmt.Sprintf("New version for %s: %s -> %s", i.ImageName, i.CurrentTag, i.LatestTag))
			}

			if ccuFlags.Update {
				if err := i.Update(); err != nil {
					slog.Error(fmt.Sprintf("error updating file: %v", err))
					continue
				}
				slog.Info(fmt.Sprintf("File [%s] | Image %s has new version %s", i.FilePath, i.ImageName, i.LatestTag))
			}

			if ccuFlags.Restart {
				if err := i.Restart(); err != nil {
					slog.Error(fmt.Sprintf("error restarting service: %v", err))
					continue
				}
				slog.Info(fmt.Sprintf("Compose file [%s] restarted", i.FilePath))
			}
		}
	}

}
