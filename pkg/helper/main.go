package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	RG()
}

func RG() {
	jsonPath := filepath.Join("shake256_hlist.json")

	// 1) Read the file's contents
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	// 2) Unmarshal into a slice of strings (your keccak hashes)
	var contents []string
	if err := json.Unmarshal(data, &contents); err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	fmt.Printf("Data from %s:\n%v\n", jsonPath, contents)

	client := &http.Client{}

	// We'll collect *all* responses in this slice
	var allResults []map[string]interface{}

	// 3) Loop over each hash in the slice and POST it to the server
	for _, hashStr := range contents {
		// Build the JSON body with "hashed_credential": <hashStr>
		requestBodyMap := map[string]string{"hashed_credential": hashStr}
		requestBodyBytes, err := json.Marshal(requestBodyMap)
		if err != nil {
			log.Printf("Error marshaling request body: %v\n", err)
			continue
		}

		// Create a new POST request
		req, err := http.NewRequest(
			http.MethodPost,
			"http://localhost:8080/generate-tree-object",
			bytes.NewBuffer(requestBodyBytes),
		)
		if err != nil {
			log.Printf("Error creating request for hash %q: %v\n", hashStr, err)
			continue
		}

		// Set headers
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		// Execute the request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making POST request for hash %q: %v\n", hashStr, err)
			continue
		}

		// We must close the body *after* reading it
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// (Optional) print to stdout for debugging
		fmt.Printf("Response for hash %q:\n%s\n", hashStr, string(respBody))

		// Attempt to parse the server response as JSON
		// so we don't get escaped strings in the final file
		var parsed interface{}
		if err := json.Unmarshal(respBody, &parsed); err != nil {
			// If it's not valid JSON, we’ll just storage the raw string
			parsed = string(respBody)
		}

		// Build a single entry that includes the hash and the parsed response
		entry := map[string]interface{}{
			"hash":     hashStr,
			"response": parsed,
		}

		// Append this entry to our master slice
		allResults = append(allResults, entry)
	}

	// 4) After the loop, write one big JSON array to "response.json"
	fileName := "response.json"
	output, err := json.MarshalIndent(allResults, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling final results: %v", err)
	}

	err = os.WriteFile(fileName, output, 0644)
	if err != nil {
		log.Fatalf("Error writing %s: %v", fileName, err)
	}

	fmt.Printf("\nAll JSON responses successfully written to %s\n", fileName)
}
