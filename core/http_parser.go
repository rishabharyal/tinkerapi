package core

import (
	"encoding/json"
	"fmt"
	"strings"
)

type HTTPParser struct {
    requestString string
}

type HTTPRequest struct {
	URL     string
	Headers map[string]string
	Params  map[string]interface{}
}

func (f *HTTPParser) parse() (*HTTPRequest, error) {
	lines := strings.Split(f.requestString, "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid HTTP request")
	}

	// Extract URL from the request line
	requestLine := lines[0]
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid request line")
	}
	url := parts[1]

	// Parse headers
	headers := make(map[string]string)
	for _, line := range lines[1:] {
		if line == "" {
			break // Reached the end of headers
		}
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) != 2 {
			return nil, fmt.Errorf("invalid header format")
		}
		key := strings.TrimSpace(headerParts[0])
		value := strings.TrimSpace(headerParts[1])
		headers[key] = value
	}

	// Parse body (params) if present
	var params map[string]interface{}
	bodyIndex := -1
	for i, line := range lines {
		if line == "" {
			bodyIndex = i + 1
			break
		}
	}
	if bodyIndex != -1 && bodyIndex < len(lines) {
		body := strings.Join(lines[bodyIndex:], "\n")
		if len(body) > 0 {
			params = make(map[string]interface{})
			err := json.Unmarshal([]byte(body), &params)
			if err != nil {
				return nil, fmt.Errorf("failed to parse body: %v", err)
			}
		}
	}

	return &HTTPRequest{
		URL:     url,
		Headers: headers,
		Params:  params,
	}, nil
}
