package main

import (
	"github.com/rbretecher/go-postman-collection"
	"strings"
)

func ExtractRequests(c *postman.Collection) ([]Endpoint, error) {
	var endpoints []Endpoint
	if c.Items != nil {
		for _, item := range c.Items {
			ep, err := processItem(item)
			if err != nil {
				return nil, err
			}
			if ep != nil {
				endpoints = append(endpoints, ep...)
			}
		}
	}
	return endpoints, nil
}

func processItem(item *postman.Items) ([]Endpoint, error) {
	var endpoints []Endpoint
	for _, it := range item.Items {
		ep, err := processItem(it)
		if err != nil {
			return nil, err
		}
		if ep != nil {
			endpoints = append(endpoints, ep...)
		}
	}

	if item.Request != nil && item.Request.URL != nil {
		p := "/" + strings.Join(item.Request.URL.Path, "/")
		//p = replacePostmanPathParams(p)
		endpoints = append(endpoints, Endpoint{
			Method: string(item.Request.Method),
			Path:   p,
		})
	}

	return endpoints, nil
}

/*func replacePostmanPathParams(path string) string {
	var re = regexp.MustCompile(`{{.*?}}`)
	path = re.ReplaceAllString(path, `*`)
	return path
}*/
