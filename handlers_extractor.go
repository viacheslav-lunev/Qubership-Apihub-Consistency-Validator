package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type Endpoint struct {
	Method string
	Path   string
}

func ExtractEndpoints(filePath string) ([]Endpoint, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []Endpoint
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()

		path, method := extractPathAndMethod(line)
		if method != "" && path != "" {
			mConst := covertFromGoToConst(method)
			if mConst == "" {
				return nil, fmt.Errorf("failed to conver method %s to const. File %s, line %d", method, filePath, lineNumber)
			}
			path = replaceGoPathParams(path)

			result = append(result, Endpoint{
				Method: mConst,
				Path:   path,
			})
		}
		lineNumber++
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, err
}

func extractPathAndMethod(line string) (string, string) {
	re := regexp.MustCompile(`.*HandleFunc\("(.*)",.*.Methods\((.*)\)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) < 3 {
		return "", ""
	}

	return matches[1], matches[2]
}

func replaceGoPathParams(path string) string {
	var re = regexp.MustCompile(`{.*?}`)
	path = re.ReplaceAllString(path, `*`)
	return path
}

func covertFromGoToConst(goConst string) string {
	switch goConst {
	case "http.MethodGet":
		return "GET"
	case "http.MethodPost":
		return "POST"
	case "http.MethodPut":
		return "PUT"
	case "http.MethodPatch":
		return "PATCH"
	case "http.MethodDelete":
		return "DELETE"
	case "http.MethodHead":
		return "HEAD"
	case "http.MethodOptions":
		return "OPTIONS"
	case "http.MethodConnect":
		return "CONNECT"
	case "http.MethodTrace":
		return "TRACE"
	}
	return ""
}
