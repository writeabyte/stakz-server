package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func writeFile(filePath string, content string) {
	// Attempt to open the file for writing. If the file doesn't exist, create it.
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening or creating file:", err)
		return
	}
	defer file.Close()

	// Write content to the file
	_, err = file.Write([]byte(content))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func main() {
	// Define a function to handle HTTP requests
	// TODO -
	var dir string
	flag.StringVar(&dir, "dir", ".", "The directory you want the stakz server to run in.")
	flag.Parse()

	os.Chdir(dir)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		files := []string{}
		err := filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					files = append(files, path)
				}
				if err != nil {
					return err
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
		body, err := json.Marshal(map[string]interface{}{
			"files": files,
		})
		fmt.Fprintf(res, string(body))
	})

	http.HandleFunc("/content", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "*")
		if req.Method == "OPTIONS" {
			res.WriteHeader(http.StatusOK)
			return
		} else if req.Method == "POST" {
			// Create a new decoder
			var data struct {
				Path        string `json:"path"`
				FileContent string `json:"fileContent"`
			}
			err := json.NewDecoder(req.Body).Decode(&data)
			fmt.Println(err)
			writeFile(data.Path, data.FileContent)
		} else if req.Method == "GET" {
			params := req.URL.Query()
			content, err := os.ReadFile(params.Get("file"))
			if content != nil {
				res.Write(content)
			} else {
				res.WriteHeader(http.StatusInternalServerError)
				res.Write([]byte(err.Error()))
			}
		}
	})

	// Basic echo http endpoint
	http.HandleFunc("/echo", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		req.Write(res)
	})

	http.HandleFunc("/execute", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		// Execute the command
		script, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}
		scriptStr := string(script)
		log.Println(scriptStr)
		cmd := exec.Command("/bin/bash", "-c", scriptStr)
		out, err := cmd.CombinedOutput()
		log.Println(string(out))
		res.Write(out)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		}
	})

	// Start the server and listen on port 8080
	log.Fatalln(http.ListenAndServe(":3001", nil))
	fmt.Println("Server listening on 3001")
}
