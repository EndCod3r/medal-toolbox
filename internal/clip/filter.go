package clip

import (
	"strings"
)

type FilterOptions struct {
	PathSearch    string
	Title         string
	CollectionID  string
	CollectionName string
	Game          string
}

func FilterClips(clips map[string]Clip, options FilterOptions) []Clip {
	var filtered []Clip

	for _, clip := range clips {
		if matchesFilter(clip, options) {
			filtered = append(filtered, clip)
		}
	}

	return filtered
}

func matchesFilter(clip Clip, options FilterOptions) bool {
	if options.PathSearch != "" && !strings.Contains(strings.ToLower(clip.FilePath), strings.ToLower(options.PathSearch)) {
		return false
	}

	if options.Title != "" && !strings.Contains(strings.ToLower(clip.GameTitle), strings.ToLower(options.Title)) {
		return false
	}

	if options.Game != "" && clip.Content.Category.CategoryName != "" &&
		!strings.Contains(strings.ToLower(clip.Content.Category.CategoryName), strings.ToLower(options.Game)) {
		return false
	}

	if options.CollectionID != "" || options.CollectionName != "" {
		collectionMatch := false
		for _, collection := range clip.Content.ContentCollections {
			if options.CollectionID != "" && collection.CollectionID == options.CollectionID {
				collectionMatch = true
				break
			}
			if options.CollectionName != "" && strings.Contains(strings.ToLower(collection.Name), strings.ToLower(options.CollectionName)) {
				collectionMatch = true
				break
			}
		}
		if !collectionMatch && (options.CollectionID != "" || options.CollectionName != "") {
			return false
		}
	}

	return true
}