package internal

import (
	"flag"
	"os"
)

type CCUFlags struct {
	Help        bool   // Show help message
	Update      bool   // Update the Docker Compose files with the new image tags
	Restart     bool   // Restart the services after updating the Docker Compose files
	Interactive bool   // Interactively choose which docker images to update
	Directory   string // Root directory to search for Docker Compose files
	Full        bool   // Update to the latest semver version
	Major       bool   // Update to the latest major version
	Minor       bool   // Update to the latest minor version
	Patch       bool   // Update to the latest patch version
	Version     bool   // Version of ccu
}

func Parse(version string) CCUFlags {
	args := CCUFlags{}

	flag.BoolVar(&args.Help, "h", false, "Show help message")
	flag.BoolVar(&args.Update, "u", false, "Update the Docker Compose files with the new image tags")
	flag.BoolVar(&args.Restart, "r", false, "Restart the services after updating the Docker Compose files")
	flag.BoolVar(&args.Interactive, "i", false, "Interactively choose which docker images to update")
	flag.StringVar(&args.Directory, "d", ".", "Root directory to search for Docker Compose files")
	flag.BoolVar(&args.Full, "f", false, "Update to the latest major version")
	flag.BoolVar(&args.Major, "major", false, "Update to the latest semver version")
	flag.BoolVar(&args.Minor, "minor", false, "Update to the latest minor version")
	flag.BoolVar(&args.Patch, "patch", true, "Update to the latest patch version")
	flag.BoolVar(&args.Version, "v", false, "Show version information")

	flag.Parse()

	if args.Version {
		println("Version:", version)
		os.Exit(0)
	}

	if args.Help {
		flag.Usage()
		os.Exit(0)
	}

	if args.Full {
		args.Major = true
		args.Minor = true
		args.Patch = true
	}

	return args
}
