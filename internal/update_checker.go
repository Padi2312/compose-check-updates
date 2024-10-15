package internal

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

type UpdateChecker struct {
	path string
}

func NewUpdateChecker(path string) *UpdateChecker {
	return &UpdateChecker{path: path}
}

func (u *UpdateChecker) Check(major, minor, patch bool) ([]UpdateInfo, error) {
	updateInfos, err := u.createUpdateInfos()
	if err != nil {
		return nil, err
	}

	for i, updateInfo := range updateInfos {
		version, err := GetSemver(updateInfo.CurrentTag)
		if err != nil {
			slog.Warn(fmt.Sprintf("Skipping (invalid semver) \t Image: %s \t Path: %s", updateInfo.ImageName, updateInfo.FilePath))
			continue
		}

		tags, err := FetchImageTags(updateInfo.ImageName)
		if err != nil {
			slog.Error(fmt.Sprintf("Skipping (failed fetching tags) \t Image: %s \t Path: %s", updateInfo.ImageName, updateInfo.FilePath))
			continue
		}

		latestVersion, err := FindLatestVersion(version, tags, major, minor, patch)
		if err != nil {
			slog.Error(fmt.Sprintf("Skipping (failed finding last version) \t Image: %s \t Path: %s", updateInfo.ImageName, updateInfo.FilePath))
			continue
		}

		if latestVersion != "" {
			updateInfos[i].LatestTag = latestVersion
		}
	}

	return updateInfos, nil
}

func (u *UpdateChecker) createUpdateInfos() ([]UpdateInfo, error) {
	var updateInfos []UpdateInfo
	uniqueImages := make(map[string]struct{})

	file, err := os.Open(u.path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	imageNamePattern := regexp.MustCompile(`^\s*image:\s*(\S+)\s*$`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := imageNamePattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			imageName := matches[1]
			name, tag := u.getNameAndTag(imageName)
			imageKey := name + ":" + tag

			if _, exists := uniqueImages[imageKey]; !exists {
				uniqueImages[imageKey] = struct{}{}
				updateInfos = append(updateInfos, UpdateInfo{
					FilePath:      u.path,
					RawLine:       line,
					FullImageName: imageName,
					ImageName:     name,
					CurrentTag:    tag,
				})
			}
		}
	}

	return updateInfos, nil
}

func (u *UpdateChecker) getNameAndTag(imageName string) (string, string) {
	parts := strings.Split(imageName, ":")
	if len(parts) < 2 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
