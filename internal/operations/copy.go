// internal/operations/copy.go
package operations

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/EndCod3r/medal-toolbox/internal/clip"
)

type CopyResult struct {
	SuccessCount int
	ErrorCount   int
	Successes    []CopySuccess
	Errors       []CopyError
}

type CopySuccess struct {
	SourcePath string
	DestPath   string
}

type CopyError struct {
	SourcePath string
	Error      error
}

func CopyClips(clips []clip.Clip, destDir string) CopyResult {
	result := CopyResult{}
	
	// Ensure destination directory exists
	if err := os.MkdirAll(destDir, 0755); err != nil {
		result.Errors = append(result.Errors, CopyError{
			SourcePath: "",
			Error:      fmt.Errorf("failed to create destination directory: %w", err),
		})
		result.ErrorCount++
		return result
	}

	for _, c := range clips {
		destPath, err := copyFile(c.FilePath, destDir)
		if err != nil {
			result.Errors = append(result.Errors, CopyError{
				SourcePath: c.FilePath,
				Error:      err,
			})
			result.ErrorCount++
			fmt.Printf("Error: %v\n", err)
		} else {
			result.Successes = append(result.Successes, CopySuccess{
				SourcePath: c.FilePath,
				DestPath:   destPath,
			})
			result.SuccessCount++
			fmt.Printf("Copied: %s\n", filepath.Base(c.FilePath))
		}
	}
	return result
}

func copyFile(src, destDir string) (string, error) {
	// Check if source file exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return "", fmt.Errorf("source file does not exist")
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer sourceFile.Close()

	destPath := filepath.Join(destDir, filepath.Base(src))
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return destPath, err
}

func WriteLogFile(destDir string, result CopyResult) error {
	logPath := filepath.Join(destDir, "copy_log_"+time.Now().Format("20060102_150405")+".txt")
	file, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header
	file.WriteString("Medal Clip Copy Log\n")
	file.WriteString("===================\n")
	file.WriteString(fmt.Sprintf("Date: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	file.WriteString(fmt.Sprintf("Total clips processed: %d\n", result.SuccessCount+result.ErrorCount))
	file.WriteString(fmt.Sprintf("Successfully copied: %d\n", result.SuccessCount))
	file.WriteString(fmt.Sprintf("Errors: %d\n", result.ErrorCount))
	file.WriteString("\n")

	// Write successful copies
	if len(result.Successes) > 0 {
		file.WriteString("SUCCESSFUL COPIES:\n")
		file.WriteString("==================\n")
		for _, success := range result.Successes {
			file.WriteString(fmt.Sprintf("Source: %s\n", success.SourcePath))
			file.WriteString(fmt.Sprintf("Destination: %s\n", success.DestPath))
			file.WriteString("\n")
		}
	}

	// Write errors
	if len(result.Errors) > 0 {
		file.WriteString("ERRORS:\n")
		file.WriteString("=======\n")
		for _, copyError := range result.Errors {
			file.WriteString(fmt.Sprintf("File: %s\n", copyError.SourcePath))
			file.WriteString(fmt.Sprintf("Error: %v\n", copyError.Error))
			file.WriteString("\n")
		}
	}

	return nil
}