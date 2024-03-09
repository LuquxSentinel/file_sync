package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileInfo stores information about a file.
type FileInfo struct {
	Path        string
	LastHash    string
	LastModTime time.Time
}

// calculateFileHash computes the MD5 hash of a file.
func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// scanDirectory scans a directory for file changes.
func scanDirectory(directoryPath string, files map[string]*FileInfo) error {
	return filepath.Walk(directoryPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if the file has been modified
		lastInfo, exists := files[filePath]
		if exists {
			lasthash, _ := calculateFileHash(filePath)
			if lastInfo.LastModTime.Equal(info.ModTime()) && lastInfo.LastHash == lasthash {
				return nil // No change
			}

			// Handle modification
			fmt.Printf("File %s has been modified!\n", filePath)
			writeBytes(filePath)
			// Perform actions to sync changes bidirectionally here
		} else {
			// Handle addition
			fmt.Printf("File %s has been added!\n", filePath)
			writeBytes(filePath)
			// Perform actions to sync changes bidirectionally here
		}

		lastHash, _ := calculateFileHash(filePath)
		// Update file information
		files[filePath] = &FileInfo{
			Path:        filePath,
			LastHash:    lastHash,
			LastModTime: info.ModTime(),
		}

		return nil
	})
}

// handleDeletedFiles checks for deleted files.
func handleDeletedFiles(files map[string]*FileInfo) {
	for filePath, _ := range files {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// Handle deletion
			fmt.Printf("File %s has been deleted!\n", filePath)

			// Perform actions to sync changes bidirectionally here

			// Remove file from the map
			delete(files, filePath)
		}
	}
}

func writeBytes(filePath string) error {
	connection, err := net.Dial("tcp", ":4000")

	paths := strings.Split(filePath, "/")
	var filename string = paths[len(paths)-1]

	log.Println(filename)

	err = gob.NewEncoder(connection).Encode(filename)

	// binary.Write(connection, binary.LittleEndian, filename)

	if err != nil {
		log.Panicln(err)
		return err
	}
	log.Println("----")
	log.Println(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Panicln(err)
	}

	//	file information [stats]
	fileStats, _ := file.Stat()

	// sizeBuffer := make([]byte, 8)
	binary.Write(connection, binary.LittleEndian, fileStats.Name())

	// binary.LittleEndian.AppendUint16(sizeBuffer, uint16(fileStats.Size()))

	// _, err = c.connection.Write(sizeBuffer)
	// if err != nil {
	// 	return err
	// }

	//	write file to connection [network ("tcp") server]
	_, err = io.Copy(connection, file)
	return err
}

// func uploadFile(filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Create a buffer to store the file contents
// 	var buffer bytes.Buffer
// 	writer := multipart.NewWriter(&buffer)

// 	// Create a form field for the file
// 	fileWriter, err := writer.CreateFormFile("file", filename)
// 	if err != nil {
// 		fmt.Println("Error creating form file:", err)
// 		return
// 	}

// 	// Use TeeReader to simultaneously read from the file and write to the form field
// 	teeReader := io.TeeReader(file, fileWriter)

// 	// Read the file contents and write to the form field
// 	if _, err := io.Copy(fileWriter, teeReader); err != nil {
// 		fmt.Println("Error copying file contents:", err)
// 		return
// 	}

// 	// Close the multipart writer to finalize the request
// 	writer.Close()

// 	// Make the HTTP POST request
// 	response, err := http.Post("http://20.20.18.60:3000/upload", writer.FormDataContentType(), &buffer)
// 	if err != nil {
// 		fmt.Println("Error making POST request:", err)
// 		return
// 	}
// 	defer response.Body.Close()

// 	// Print the server response
// 	fmt.Println("Server Response:", response.Status)
// }

// func main() {
// 	// Directory to monitor
// 	directoryPath := "/home/wethinkcode_/Documents/test/"

// 	// Map to store file information
// 	files := make(map[string]*FileInfo)

// 	// Periodically scan for file changes
// 	ticker := time.NewTicker(5 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			// Scan for modifications and additions
// 			err := scanDirectory(directoryPath, files)
// 			if err != nil {
// 				fmt.Printf("Error scanning directory: %v\n", err)
// 			}

// 			// Check for deleted files
// 			handleDeletedFiles(files)
// 		}
// 	}
// }
