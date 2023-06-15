package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Config struct voor het opslaan van configuratie-instellingen
type Config struct {
	ServerPort     string `json:"server_port"`
	FilePath       string `json:"file_path"`
	LogFile        string `json:"log_file"`
	WelcomeMessage string `json:"welcome_message"`
	ErrorMessage   string `json:"error_message"`
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(config.LogFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logFile, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stderr, logFile)) // Aanpassing

	err = helloWorld(config.WelcomeMessage)
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			content, err := ioutil.ReadFile(config.FilePath)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err)
				return
			}

			w.Header().Set("Content-Type", "text/html")
			w.Write(content)
		} else {
			http.NotFound(w, r)
		}
	})

	go func() {
		err := http.ListenAndServe(":"+config.ServerPort, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Server gestart op http://localhost:" + config.ServerPort)
	select {}
}

func helloWorld(message string) error {
	fmt.Println(message)
	return fmt.Errorf("Er is een fout opgetreden")
}

// Functie voor het laden van de configuratie vanuit een JSON-bestand
func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
