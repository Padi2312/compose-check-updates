package main

import (
	"os"
	"os/exec"
	"strings"
)

type UpdateInfo struct {
	FilePath      string
	RawLine       string
	ImageName     string
	FullImageName string
	CurrentTag    string
	LatestTag     string
}

func (u *UpdateInfo) HasNewVersion() bool {
	if u.CurrentTag == "" || u.LatestTag == "" {
		return false
	}

	current, err := GetSemver(u.CurrentTag)
	if err != nil {
		return false
	}

	latest, err := GetSemver(u.LatestTag)
	if err != nil {
		return false
	}

	return latest.GreaterThan(current)
}

func (u *UpdateInfo) Update() error {
	input, err := os.ReadFile(u.FilePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, u.RawLine) {
			lines[i] = strings.Replace(line, u.CurrentTag, u.LatestTag, 1)
		}
	}

	output := strings.Join(lines, "\n")
	err = os.WriteFile(u.FilePath, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (u *UpdateInfo) Restart() error {
	// Execute docker-compose up -d here with os
	cmd := exec.Command("docker-compose", "-f", u.FilePath, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
