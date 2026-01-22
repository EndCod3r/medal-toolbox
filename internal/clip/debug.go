package clip

import (
	"fmt"
)

// DebugGameNames prints the game names for the first few clips to help diagnose issues
func DebugGameNames(clips map[string]Clip, count int) {
	fmt.Printf("Debug: Showing game names for first %d clips:\n", count)
	i := 0
	for _, clip := range clips {
		if i >= count {
			break
		}
		gameName := clip.GetGameName()
		fmt.Printf("Clip %d: GameTitle='%s', GameName='%s', HasCategory=%v\n", 
			i+1, clip.GameTitle, gameName, clip.Content.Category != nil)
		i++
	}
}