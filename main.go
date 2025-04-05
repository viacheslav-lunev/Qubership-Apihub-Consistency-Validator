package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/rbretecher/go-postman-collection"
)

type Config struct {
	PostmanCollections     []string `json:"postman_collections"`
	SourceCodeHTTPHandlers []string `json:"source_code_http_handlers"`
	ApihubURL              string   `json:"apihub_url"`
	PackageID              string   `json:"package_id"`
	Version                string   `json:"version"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the configuration file as an argument.")
	}
	filePath := os.Args[1]

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	var serverEndpoints []Endpoint
	for _, file := range config.SourceCodeHTTPHandlers {
		parsedEndpoints, err := ExtractEndpoints(file)
		if err != nil {
			log.Fatalf("Error extracting serverEndpoints: %v", err)
		}
		serverEndpoints = append(serverEndpoints, parsedEndpoints...)
	}
	serverEndpoints = removeEndpointDuplicates(serverEndpoints)

	var requestEndpoints []Endpoint
	for _, fp := range config.PostmanCollections {
		file, err := os.Open(fp)
		if err != nil {
			if file != nil {
				file.Close()
			}
			log.Fatalf("Error opening file %s: %v", fp, err)
		}

		c, err := postman.ParseCollection(file)
		if err != nil {
			log.Fatalf("Error parsing collection: %v", err)
		}
		if c != nil {
			re, err := ExtractRequests(c)
			if err != nil {
				log.Fatalf("Error extracting requests: %v", err)
			}
			requestEndpoints = append(requestEndpoints, re...)
		}
		file.Close()
	}

	for _, rep := range requestEndpoints {
		found := false
		for _, sep := range serverEndpoints {
			if sep.Method != rep.Method {
				continue
			}
			serverPathRegex := strings.Replace(sep.Path, "*", ".*?", -1)
			re, err := regexp.Compile(serverPathRegex)
			if err != nil {
				fmt.Printf("warn, regexp not compiled: %s\n", err)
				continue
			}

			if re.MatchString(rep.Path) {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Endpoint %+v not found in implementation\n", rep)
		}
	}

	for _, sep := range serverEndpoints {
		serverPathRegex := strings.Replace(sep.Path, "*", ".*?", -1)
		re, err := regexp.Compile(serverPathRegex)
		if err != nil {
			fmt.Printf("warn, regexp not compiled: %s\n", err)
			continue
		}

		found := false
		for _, rep := range requestEndpoints {
			if sep.Method != rep.Method {
				continue
			}
			if re.MatchString(rep.Path) {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Endpoint %+v not found in postman collections\n", sep)
		}
	}
}

func matchPostmanReqToGoEndpoint() {

}

func removeEndpointDuplicates(input []Endpoint) []Endpoint {
	var result []Endpoint
	serverEndpointsSet := make(map[string]struct{}, 0)
	for _, e := range input {
		key := e.Method + "|" + e.Path
		if _, exists := serverEndpointsSet[key]; !exists {
			serverEndpointsSet[key] = struct{}{}
			result = append(result, e)
		}
	}
	return result
}
