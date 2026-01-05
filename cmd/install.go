/*
Copyright © 2025 ZOLLIDAN zollidan@aol.com
*/
package cmd

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

type LiveryJSON struct {
	Rev            int    `json:"rev"`
	ProductId      int    `json:"productId"`
	ProductPackage string `json:"productPackage"`
	Title          string `json:"title"`
	Airline        string `json:"airline"`
	AirlineIcao    string `json:"airlineIcao"`
	ATCID          string `json:"atcId"`
	LiveryID       string `json:"liveryId"`
	Version        int    `json:"version"`
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [zip-file-full-path] [community-folder-path]",
	Short: "A brief description of your command",
	Long:  `A longer description `,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		zipPath := args[0]
		communityPath := args[1]

		_, err := findLayoutGenerator(communityPath)
		if err != nil {
			log.Fatal("Could not find layout generator")
			return
		}

		planeFolder := "pmdg-aircraft-77f-liveries"
		fullPath := filepath.Join(communityPath, planeFolder, "SimObjects", "Airplanes")

		// standard folder permissions
		err = os.MkdirAll(fullPath, 0755)
		if err != nil {
			log.Printf("Error creating directory: %v\n", err)
			return
		}

		_, err = os.Stat(zipPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Printf("Error: File does not exist, %s", err.Error())
				return
			}
			log.Printf("Error accessing the file: %s", err.Error())
			return
		}

		r, err := zip.OpenReader(zipPath)
		if err != nil {
			log.Printf("Error opening zip file: %s", err.Error())
			return
		}

		defer r.Close()

		var found bool
		var data io.ReadCloser

		// check if livery.json exists and open it
		for _, f := range r.File {
			if f.Name == "livery.json" {
				found = true
				var err error
				data, err = f.Open()
				if err != nil {
					fmt.Printf("Error opening livery.json: %s\n", err)
					return
				}
				defer data.Close()
				break
			}
		}

		if !found {
			fmt.Println("Error: livery.json not found in zip file")
			return
		}

		// read file to []byte
		body, err := io.ReadAll(data)
		if err != nil {
			fmt.Printf("Error reading livery.json: %s\n", err)
			return
		}

		// read json from []byte
		var livery LiveryJSON
		if err := json.Unmarshal(body, &livery); err != nil {
			fmt.Printf("Error parsing livery.json: %s\n", err)
			return
		}

		targetDir := filepath.Join(fullPath, livery.LiveryID)
		err = os.MkdirAll(targetDir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}

		// for _, f := range r.File {
		// 	fpath := filepath.Join(targetDir, f.Name)

		// 	if f.FileInfo().IsDir() {
		// 		os.MkdirAll(fpath, f.Mode())
		// 		continue
		// 	}

		// 	if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
		// 		fmt.Printf("Error creating sub-dirs: %v\n", err)
		// 		return
		// 	}

		// 	rc, err := f.Open()
		// 	if err != nil {
		// 		fmt.Printf("Error opening zip content [%s]: %v\n", f.Name, err)
		// 		return
		// 	}

		// 	outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		// 	if err != nil {
		// 		rc.Close()
		// 		fmt.Printf("Error creating file on disk: %v\n", err)
		// 		return
		// 	}

		// 	_, err = io.Copy(outFile, rc)

		// 	outFile.Close()
		// 	rc.Close()

		// 	if err != nil {
		// 		fmt.Printf("Error copying file %s: %v\n", f.Name, err)
		// 		return
		// 	}
		// }

		liveryFolderContent, err := os.ReadDir(targetDir)
		if err != nil {
			log.Fatalf("Error reading folder: %v", err.Error())
		}

		for _, f := range liveryFolderContent {
			options := "options.ini"
			if f.Name() == options {
				err := os.Rename(options, fmt.Sprintf("%s.ini", livery.ATCID))
				if err != nil {
					log.Fatalf("Could not rename: %s", options)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func findLayoutGenerator(communityPath string) (string, error) {
	generatorPath := filepath.Join(communityPath, "MSFSLayoutGenerator.exe")

	if _, err := os.Stat(generatorPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("MSFSLayoutGenerator.exe not found in current directory: %s", generatorPath)
		}
		return "", fmt.Errorf("error checking MSFSLayoutGenerator.exe: %w", err)
	}

	return generatorPath, nil
}
