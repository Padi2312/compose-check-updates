package internal

import (
	"flag"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected CCUFlags
	}{
		{
			name: "default values",
			args: []string{},
			expected: CCUFlags{
				Help:        false,
				Update:      false,
				Restart:     false,
				Interactive: false,
				Directory:   ".",
				Full:        false,
				Major:       false,
				Minor:       false,
				Patch:       true,
			},
		},
		{
			name: "update flag",
			args: []string{"-u"},
			expected: CCUFlags{
				Update: true,
				Patch:  true,
			},
		},
		{
			name: "full flag",
			args: []string{"-f"},
			expected: CCUFlags{
				Full:  true,
				Major: true,
				Minor: true,
				Patch: true,
			},
		},
		{
			name: "directory flag",
			args: []string{"-d", "/path/to/dir"},
			expected: CCUFlags{
				Directory: "/path/to/dir",
				Patch:     true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save the original command-line arguments and restore them after the test
			origArgs := os.Args
			defer func() { os.Args = origArgs }()

			// Set the command-line arguments for the test
			os.Args = append([]string{"cmd"}, tt.args...)

			// Reset the flags to their default state
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			// Parse the flags
			exitCode := 0
			flag.CommandLine.Usage = func() {
				exitCode = 2
			}
			err := flag.CommandLine.Parse(os.Args[1:])
			if err != nil {
				exitCode = 2
			}

			result := Parse()
			if exitCode != 0 {
				return
			}

			// Compare the parsed flags with the expected values
			if result != tt.expected {
				t.Errorf("Parse() = %+v, expected %+v", result, tt.expected)
			}
		})
	}
}