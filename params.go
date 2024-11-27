package main

import "strings"

type Params map[string]string
type Processor func(value string) string

func parseParams(url string) Params {
	pipeline := []Processor{
		func(value string) string {
			return strings.ReplaceAll(value, "/", "")
		},
		func(value string) string {
			return strings.ReplaceAll(value, "?", "&")
		},
	}

	var processed string = url
	for _, processor := range pipeline {
		processed = processor(processed)
	}

	params := Params{}
	parts := strings.Split(strings.ReplaceAll(processed, "?", "&"), "&")
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		kv := strings.Split(part, "=")
		params[kv[0]] = kv[1]
	}

	return params
}
