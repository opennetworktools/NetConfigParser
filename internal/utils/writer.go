package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opennetworktools/NetConfigParser/internal/parser"
)

func WriteParserStructToJSON(data *parser.Parser, outFileName string) error {
	// Create the "out" directory if it doesn't exist
	outDir := "out"
	err := os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Ensure outFileName is within the "out" directory
	if !filepath.HasPrefix(outFileName, outDir) {
		outFileName = filepath.Join(outDir, outFileName)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// Create a file to save the JSON data
	jsonFile, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Write the JSON data to the file
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
