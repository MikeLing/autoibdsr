package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/ngaut/log"
)

type Config struct {
	Port string `toml:"port"`

	JSONFilePath string `toml:"jsonFilePath"`
}

var cfg *Config

// parseConfig reads info from config file
func parseConfig(filePath string) *Config {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("Failed to read config file in %s", filePath)
		return nil
	}

	if err = toml.Unmarshal(data, &cfg); err != nil {
		log.Errorf("Failed to decode config file in %s", filePath)
		return nil
	}

	return cfg
}

// deliverJSON warp json data add provide a api could access it
func deliverJSON(w http.ResponseWriter, r *http.Request) {
	// Parse parameters from url
	r.ParseForm()
	fileName := r.Form.Get("fileName")
	if len(fileName) <= 0 {
		log.Errorf("Must give a file name to continue")
		w.WriteHeader(http.StatusNotFound)
	}

	// join file path together
	JSONFilePath := fmt.Sprintf("%s%s", cfg.JSONFilePath, fileName)

	data, err := ioutil.ReadFile(JSONFilePath)
	if err != nil {
		log.Errorf("Failed to read json file in %s", JSONFilePath)
	}

	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Errorf("Failed to redent data to json")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonData)

}

func main() {
	parseConfig("config.toml")
	fmt.Printf("Server started successfully! Running on %s and store json file in %s", cfg.Port, cfg.JSONFilePath)
	http.HandleFunc("/search", deliverJSON)
	http.ListenAndServe(cfg.Port, nil)
}
