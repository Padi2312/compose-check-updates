package internal

import (
	"flag"
	"os"
)

type CCUFlags struct {
	Help        bool
	Update      bool
	Restart     bool
	Interactive bool
	Directory   string
}

func Parse() CCUFlags {
	args := CCUFlags{}

	flag.BoolVar(&args.Help, "h", false, "Show help message")
	flag.BoolVar(&args.Update, "u", false, "Update the Docker Compose files with the new image tags")
	flag.BoolVar(&args.Restart, "r", false, "Restart the services after updating the Docker Compose files")
	flag.BoolVar(&args.Interactive, "i", false, "Interactively choose which docker images to update")
	flag.StringVar(&args.Directory, "d", ".", "Root directory to search for Docker Compose files")

	flag.Parse()

	if args.Help {
		flag.Usage()
		os.Exit(0)
	}

	return args
}
