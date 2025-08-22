package clip

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Clip struct {
	UUID      string                 `json:"uuid"`
	ClipID    string                 `json:"clipID"`
	Status    string                 `json:"Status"`
	FilePath  string                 `json:"FilePath"`
	GameTitle string                 `json:"GameTitle"`
	Content   Content                `json:"Content"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type Content struct {
	ContentCollections []ContentCollection `json:"contentCollections"`
	Category           Category            `json:"category"`
}

type ContentCollection struct {
	CollectionID string `json:"collectionId"`
	Name         string `json:"name"`
}

type Category struct {
	CategoryName string `json:"categoryName"`
}

func LoadClipsFromFile(filePath string) (map[string]Clip, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var clips map[string]Clip
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&clips)
	if err != nil {
		return nil, err
	}

	return clips, nil
}

func GetDefaultJSONPath() string {
	return filepath.Join(os.Getenv("APPDATA"), "Medal", "store", "clips.json")
}