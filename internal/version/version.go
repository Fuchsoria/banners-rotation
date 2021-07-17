package version

import (
	"encoding/json"
	"fmt"
	"os"
)

type Version struct {
	Release   string
	BuildDate string
	GitHash   string
}

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func PrintVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(Version{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
