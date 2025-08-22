// internal/clip/collections.go
package clip

import (
	"fmt"
	"strings"
)

type CollectionInfo struct {
	ID   string
	Name string
}

func GetAllCollections(clips map[string]Clip) []CollectionInfo {
	collectionMap := make(map[string]string) // ID -> Name
	
	for _, clip := range clips {
		for _, collection := range clip.Content.ContentCollections {
			if collection.CollectionID != "" && collection.Name != "" {
				collectionMap[collection.CollectionID] = collection.Name
			}
		}
	}
	
	var collections []CollectionInfo
	for id, name := range collectionMap {
		collections = append(collections, CollectionInfo{ID: id, Name: name})
	}
	
	return collections
}

func SearchCollectionsByName(collections []CollectionInfo, searchTerm string) []CollectionInfo {
	var results []CollectionInfo
	searchTerm = strings.ToLower(searchTerm)
	
	for _, collection := range collections {
		if strings.Contains(strings.ToLower(collection.Name), searchTerm) {
			results = append(results, collection)
		}
	}
	
	return results
}

func PrintCollections(collections []CollectionInfo) {
	if len(collections) == 0 {
		fmt.Println("No collections found.")
		return
	}
	
	fmt.Println("Collections found:")
	for _, collection := range collections {
		fmt.Printf("  ID: %s, Name: %s\n", collection.ID, collection.Name)
	}
}