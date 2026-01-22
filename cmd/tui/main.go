// cmd/tui/main.go
package main

import (
	"flag"

	"github.com/EndCod3r/medal-toolbox/internal/clip"
	"github.com/EndCod3r/medal-toolbox/internal/tui"
)

func main() {
	jsonPath := flag.String("json", clip.GetDefaultJSONPath(), "Path to clips JSON file")
	flag.Parse()

	tui.StartTUI(*jsonPath)
}