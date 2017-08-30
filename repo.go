package main

import (
	"encoding/json"
	"net/http"
)

// RepoInfo holds information about a GitHub repository.
type RepoInfo struct {
	Size int `json:"size"`
}

// https://api.github.com/repos/git/git
func getJSON(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if Username != "" && Token != "" {
		req.SetBasicAuth(Username, Token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(target)
}
