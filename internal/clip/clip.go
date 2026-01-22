package clip

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Clip struct {
	UUID      string                 `json:"uuid"`
	ClipID    string                 `json:"clipID"`
	Status    string                 `json:"Status"`
	FilePath  string                 `json:"FilePath"`
	GameTitle string                 `json:"GameTitle"` // This is the clip title
	Content   Content                `json:"Content"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type Content struct {
	ContentCollections []ContentCollection `json:"contentCollections"`
	Category           *Category           `json:"category"` // Make it a pointer since it might be missing
	CategoryID         string              `json:"categoryId"`
}

type ContentCollection struct {
	CollectionID string `json:"collectionId"`
	Name         string `json:"name"`
}

type Category struct {
	CategoryName string `json:"categoryName"`
}

// GetGameName safely extracts the game name from a clip
func (c *Clip) GetGameName() string {
	// First try to get from Content.Category.CategoryName
	if c.Content.Category != nil && c.Content.Category.CategoryName != "" {
		return c.Content.Category.CategoryName
	}
	
	// If not available, we might need to look in other places
	// For now, return a default or try to extract from file path
	if c.FilePath != "" {
		// Try to extract game name from file path as fallback
		if strings.Contains(c.FilePath, "Rainbow Six Siege") {
			return "Tom Clancy's Rainbow Six Siege"
		}
		// Add more game detection logic as needed
	}
	
	return "Unknown Game"
}

func LoadClipsFromFile(filePath string) (map[string]Clip, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var raw map[string]json.RawMessage
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&raw); err != nil {
		return nil, err
	}

	clips := make(map[string]Clip)

	for key, value := range raw {
		var c Clip
		if err := json.Unmarshal(value, &c); err != nil {
			// Skip non-clip entries (numbers, flags, etc)
			continue
		}

		// Some entries rely on the map key as UUID
		if c.UUID == "" {
			c.UUID = key
		}

		clips[key] = c
	}

	return clips, nil
}


func GetDefaultJSONPath() string {
	return filepath.Join(os.Getenv("APPDATA"), "Medal", "store", "clips.json")
}