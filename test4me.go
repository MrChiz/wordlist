package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"strings"

	"github.com/google/go-cmp/cmp"
)

// HackerOne struct
type HackerOneInfo struct {
	Scope []struct {
		ASSET      string `json:"asset_identifier"`
		ASSET_TYPE string `json:"asset_type"`
		BOUNTY     bool   `json:"eligible_for_bounty"`
	} `json:"in_scope"`
}

type HackerOne struct {
	NAME   string       `json:"name"`
	STATUS bool         `json:"offers_bounties"`
	TARGET string       `json:"website"`
	INFO   HackerOneInfo `json:"targets"`
}

// BugCrowd struct
type BugCrowdInfo struct {
	Scope []struct {
		TYPE   string `json:"type"`
		TARGET string `json:"target"`
	} `json:"in_scope"`
}

type BugCrowd struct {
	NAME string      `json:"name"`
	INFO BugCrowdInfo `json:"targets"`
}

// Federacy struct
type Federacy struct {
	NAME   string      `json:"name"`
	STATUS bool        `json:"offers_awards"`
	INFO   BugCrowdInfo `json:"targets"`
}

// HackenProof struct
type HackenProof struct {
	NAME   string      `json:"name"`
	STATUS bool        `json:"archived"`
	INFO   BugCrowdInfo `json:"targets"`
}

// Intigriti struct
type IntigritiInfo struct {
	Scope []struct {
		TYPE        string `json:"type"`
		ENDPOINT    string `json:"endpoint"`
		DESCRIPTION string `json:"description"`
		IMPACT      string `json:"impact"`
	} `json:"in_scope"`
}

type Intigriti struct {
	NAME   string       `json:"name"`
	HANDLE string       `json:"handle"`
	STATUS string       `json:"status"`
	INFO   IntigritiInfo `json:"targets"`
}

// YesweHack struct
type YesweHack struct {
	NAME   string      `json:"name"`
	STATUS bool        `json:"disabled"`
	INFO   BugCrowdInfo `json:"targets"`
}

func main() {
	var wg sync.WaitGroup
	baseURL := "https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/"
	list := []string{"hackerone_data.json", "bugcrowd_data.json", "intigriti_data.json", "federacy_data.json", "hackenproof_data.json", "yeswehack_data.json"}
	for {
		wg.Add(len(list))
		for _, name := range list {
			go GetData(fmt.Sprintf("%s%s", baseURL, name), name, strings.Split(name, "_")[0], &wg)
		}
		wg.Wait()
		time.Sleep(30 * time.Minute)
	}
}

func GetData(url, fileName, platform string, wg *sync.WaitGroup) {
	defer wg.Done()
	filePath := fmt.Sprintf("Data/%s", fileName)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var db interface{}
		err = json.Unmarshal(body, &db)
		if err != nil {
			log.Fatal(err)
		}

		jsonData, err := json.MarshalIndent(db, "", "\t")
		if err != nil {
			log.Fatal(err)
		}

		SaveData(filePath, jsonData)
		Comparison(jsonData, platform, filePath)
	}
}

func SaveData(filePath string, data []byte) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	file.Sync()
}

func Comparison(jsonData []byte, platform, filePath string) {
	var currentData interface{}
	err := json.Unmarshal(jsonData, &currentData)
	if err != nil {
		log.Fatal(err)
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			SaveData(filePath, jsonData)
		}
		log.Fatal(err)
	}

	var storedData interface{}
	err = json.Unmarshal(fileData, &storedData)
	if err != nil {
		log.Fatal(err)
	}

	diff := cmp.Diff(storedData, currentData)
	if diff != "" {
		fmt.Printf("Differences found in %s:\n%s\n", platform, diff)
		SaveData(filePath, jsonData)
	}
}
