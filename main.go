package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var serverKey string
var executeEnabled bool

func generateRandomKey() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal("Error generating random key:", err)
	}
	return hex.EncodeToString(bytes)
}

func writeFile(filePath string, content string) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening or creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write([]byte(content))
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func keyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		key := req.Header.Get("Authorization")
		if key != serverKey {
			res.WriteHeader(http.StatusForbidden)
			res.Write([]byte("Invalid server key"))
			return
		}
		next.ServeHTTP(res, req)
	})
}

func main() {
	serverKey = generateRandomKey()
	fmt.Println("Server Key:", serverKey)

	var dir string
	var port int
	flag.StringVar(&dir, "dir", ".", "The directory you want the stakz server to run in.")
	flag.BoolVar(&executeEnabled, "execute", false, "Enable the /execute endpoint.")
	flag.IntVar(&port, "port", 3001, "The port the server will listen on.")
	flag.Parse()

	err := os.Chdir(dir)

	if err != nil {
		log.Fatal("Error changing working directory:", err)
	}

	http.Handle("/", keyMiddleware(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
	})))

	http.Handle("/content", keyMiddleware(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("Content Request received")
		host := req.Header.Get("Origin")
		if strings.HasPrefix(host, "http://localhost") || strings.HasPrefix(host, "https://stakz.dev") {
			res.Header().Set("Access-Control-Allow-Origin", host)
		} else {
			res.WriteHeader(http.StatusForbidden)
			res.Write([]byte("invalid host: " + host))
			return
		}
		if req.Method == "OPTIONS" {
			res.WriteHeader(http.StatusOK)
			return
		} else if req.Method == "POST" {
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
	})))

	http.Handle("/echo", keyMiddleware(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		req.Write(res)
	})))

	http.Handle("/execute", keyMiddleware(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if !executeEnabled {
			res.WriteHeader(http.StatusForbidden)
			res.Write([]byte("Execute endpoint disabled. Please enable it with the --execute flag."))
			return
		}
		res.Header().Set("Access-Control-Allow-Origin", "*")
		res.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")
		script, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
			return
		}
		scriptStr := string(script)
		log.Println(scriptStr)
		cmd := exec.Command("/bin/sh", "-c", scriptStr)
		out, err := cmd.Output()
		log.Println(string(out))
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(err.Error()))
		} else {
			res.Write(out)
		}
	})))

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Server listening on %s\n", addr)
	log.Fatalln(http.ListenAndServe(addr, nil))
}
