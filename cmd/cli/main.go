package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/EndCod3r/medal-toolbox/internal/clip"
	"github.com/EndCod3r/medal-toolbox/internal/operations"
)

func main() {
	jsonPath := flag.String("json", clip.GetDefaultJSONPath(), "Path to clips JSON file")
	copyDir := flag.String("copy-dir", "", "Directory to copy clips to")
	pathSearch := flag.String("path-search", "", "Filter by text in file path")
	title := flag.String("title", "", "Filter by clip title")
	collectionID := flag.String("collection-id", "", "Filter by collection ID")
	collectionName := flag.String("collection-name", "", "Filter by collection name")
	game := flag.String("game", "", "Filter by game name")
	listCollections := flag.String("list-collections", "", "Search for collections by name")

	flag.Parse()

	// If list-collections flag is provided, just list collections and exit
	if *listCollections != "" {
		listCollectionsMode(*jsonPath, *listCollections)
		return
	}

	// Interactive mode if no flags provided
	if flag.NFlag() == 0 {
		interactiveMode()
		return
	}

	// Validate required parameters
	if *copyDir == "" {
		fmt.Println("Error: --copy-dir is required")
		os.Exit(1)
	}

	// Load clips
	clips, err := clip.LoadClipsFromFile(*jsonPath)
	if err != nil {
		fmt.Printf("Error loading clips: %v\n", err)
		os.Exit(1)
	}

	// Apply filters
	filterOptions := clip.FilterOptions{
		PathSearch:     *pathSearch,
		Title:          *title,
		CollectionID:   *collectionID,
		CollectionName: *collectionName,
		Game:           *game,
	}

	filteredClips := clip.FilterClips(clips, filterOptions)

	if len(filteredClips) == 0 {
		fmt.Println("No clips match the specified filters")
		return
	}

	fmt.Printf("Found %d matching clips\n", len(filteredClips))

	// Copy clips with error handling
	result := operations.CopyClips(filteredClips, *copyDir)
	
	fmt.Printf("\nCopy operation completed with %d successes and %d errors\n", 
		result.SuccessCount, result.ErrorCount)
	
	if result.ErrorCount > 0 {
		fmt.Println("\nErrors encountered:")
		for _, copyErr := range result.Errors {
			fmt.Printf("  - %s: %v\n", copyErr.SourcePath, copyErr.Error)
		}
	}

	// Write log file
	if err := operations.WriteLogFile(*copyDir, result); err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
	} else {
		logPath := filepath.Join(*copyDir, "copy_log_*.txt")
		fmt.Printf("Detailed log file created in the destination directory: %s\n", logPath)
	}
}

func listCollectionsMode(jsonPath, searchTerm string) {
	clips, err := clip.LoadClipsFromFile(jsonPath)
	if err != nil {
		fmt.Printf("Error loading clips: %v\n", err)
		os.Exit(1)
	}
	
	allCollections := clip.GetAllCollections(clips)
	
	if searchTerm != "" {
		matchingCollections := clip.SearchCollectionsByName(allCollections, searchTerm)
		clip.PrintCollections(matchingCollections)
	} else {
		clip.PrintCollections(allCollections)
	}
}

func interactiveMode() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== Medal Clip Manager - Interactive Mode ===")

	// Get JSON path
	fmt.Printf("JSON file path [%s]: ", clip.GetDefaultJSONPath())
	jsonPath, _ := reader.ReadString('\n')
	jsonPath = strings.TrimSpace(jsonPath)
	if jsonPath == "" {
		jsonPath = clip.GetDefaultJSONPath()
	}

	// Load clips
	clips, err := clip.LoadClipsFromFile(jsonPath)
	if err != nil {
		fmt.Printf("Error loading clips: %v\n", err)
		os.Exit(1)
	}

	// Offer to list collections first
	fmt.Print("Do you want to search for collections first? (y/N): ")
	searchCollections, _ := reader.ReadString('\n')
	searchCollections = strings.TrimSpace(strings.ToLower(searchCollections))
	
	if searchCollections == "y" || searchCollections == "yes" {
		fmt.Print("Enter collection name to search for (leave empty to list all): ")
		collectionSearch, _ := reader.ReadString('\n')
		collectionSearch = strings.TrimSpace(collectionSearch)
		
		allCollections := clip.GetAllCollections(clips)
		if collectionSearch != "" {
			matchingCollections := clip.SearchCollectionsByName(allCollections, collectionSearch)
			clip.PrintCollections(matchingCollections)
		} else {
			clip.PrintCollections(allCollections)
		}
		
		fmt.Println() // Add a blank line
	}

	// Get filters
	fmt.Print("Search text in path (press Enter to skip): ")
	pathSearch, _ := reader.ReadString('\n')
	pathSearch = strings.TrimSpace(pathSearch)

	fmt.Print("Search text in title (press Enter to skip): ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Search by game name (press Enter to skip): ")
	game, _ := reader.ReadString('\n')
	game = strings.TrimSpace(game)

	fmt.Print("Search by collection ID (press Enter to skip): ")
	collectionID, _ := reader.ReadString('\n')
	collectionID = strings.TrimSpace(collectionID)

	fmt.Print("Search by collection name (press Enter to skip): ")
	collectionName, _ := reader.ReadString('\n')
	collectionName = strings.TrimSpace(collectionName)

	// Apply filters
	filterOptions := clip.FilterOptions{
		PathSearch:     pathSearch,
		Title:          title,
		CollectionID:   collectionID,
		CollectionName: collectionName,
		Game:           game,
	}

	filteredClips := clip.FilterClips(clips, filterOptions)

	if len(filteredClips) == 0 {
		fmt.Println("No clips match the specified filters")
		return
	}

	fmt.Printf("Found %d matching clips\n", len(filteredClips))

	// Get copy directory
	fmt.Print("Directory to copy clips to: ")
	copyDir, _ := reader.ReadString('\n')
	copyDir = strings.TrimSpace(copyDir)

	if copyDir == "" {
		fmt.Println("Error: Copy directory is required")
		os.Exit(1)
	}

	// Copy clips with error handling
	result := operations.CopyClips(filteredClips, copyDir)
	
	fmt.Printf("\nCopy operation completed with %d successes and %d errors\n", 
		result.SuccessCount, result.ErrorCount)
	
	if result.ErrorCount > 0 {
		fmt.Println("\nErrors encountered:")
		for _, copyErr := range result.Errors {
			fmt.Printf("  - %s: %v\n", copyErr.SourcePath, copyErr.Error)
		}
	}

	// Write log file
	if err := operations.WriteLogFile(copyDir, result); err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
	} else {
		logPath := filepath.Join(copyDir, "copy_log_*.txt")
		fmt.Printf("Detailed log file created in the destination directory: %s\n", logPath)
	}
}