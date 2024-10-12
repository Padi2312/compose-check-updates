package main

import (
	"bufio"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

type UpdateInfo struct {
	RawLine       string
	ImageName     string
	FullImageName string
	CurrentTag    string
	LatestTag     string
}

type UpdateChecker struct {
	path string
}

func NewUpdateChecker(path string) *UpdateChecker {
	return &UpdateChecker{path: path}
}

func (u *UpdateChecker) Check() ([]UpdateInfo, error) {
	// Get the image name from the compose file
	imageNames, err := u.getImageNames()
	if err != nil {
		return nil, err
	}

	updateInfos, err := u.getUpdateInfos(imageNames)
	if err != nil {
		return nil, err
	}

	return updateInfos, nil
}

func (u *UpdateChecker) getUpdateInfos(imageNames []string) ([]UpdateInfo, error) {
	var updateInfos []UpdateInfo
	for _, imageName := range imageNames {
		parts := strings.Split(imageName, ":")
		if len(parts) < 2 {
			slog.Warn("Skipping image %s because it does not have a tag", imageName, imageName)
			continue
		}

		imageName := parts[0]
		imageTag := parts[1]
		version, err := GetSemver(imageTag)
		if err != nil {
			slog.Warn("Skipping image %s because it does not have a semver tag", imageName, imageName)
			continue
		}

		tags, err := FetchImageTags(imageName)
		if err != nil {
			slog.Error("Skipping image %s because could not fetch tags", imageName, imageName)
			continue
		}

		latestVersion, err := FindLatestVersion(version, tags)
		if err != nil {
			slog.Error("Skipping image %s because could not find latest version", imageName, imageName)
			continue
		}

		if latestVersion != "" {
			updateInfos = append(updateInfos, UpdateInfo{
				RawLine:    imageName + ":" + imageTag,
				ImageName:  imageName,
				CurrentTag: version.Original(),
				LatestTag:  latestVersion,
			})
		}
	}
	return updateInfos, nil
}

func (u *UpdateChecker) getImageTagMap(imageNames []string) (map[string]string, error) {
	imageTagMap := make(map[string]string)
	for _, imageName := range imageNames {
		parts := strings.Split(imageName, ":")
		if len(parts) < 2 {
			slog.Warn("Skipping image %s because it does not have a tag", imageName, imageName)
			continue
		}

		imageName := parts[0]
		imageTag := parts[1]
		imageTagMap[imageName] = imageTag
	}
	return imageTagMap, nil
}

func (u *UpdateChecker) initUpdateInfos() ([]UpdateInfo, error) {
	var updateInfos []UpdateInfo

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
			updateInfos = append(updateInfos, UpdateInfo{
				RawLine:       imageName,
				FullImageName: imageName,
			})
		}
	}

	return updateInfos, nil
}

func (u *UpdateChecker) getImageNames() ([]string, error) {

	file, err := os.Open(u.path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var imageNames []string
	imageNamePattern := regexp.MustCompile(`^\s*image:\s*(\S+)\s*$`)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := imageNamePattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			imageName := matches[1]
			imageNames = append(imageNames, imageName)
		}
	}

	return imageNames, nil
}
