package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

// download downloads file from the url.
func download(url string) error {
	file, err := fileName(url)
	if err != nil {
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// fileName extracts filename from the URL.
func fileName(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	_, file := path.Split(u.Path)
	return file, nil
}

var api_url = "https://api.github.com"

// getDownloadUrls returns URLs for downloading assets from the latest repo release.
func getDownloadUrls(repo string) ([]string, error) {
	url := api_url + "/repos/" + repo + "/releases/latest"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := struct {
		Assets []struct {
			BrowserDownloadUrl string `json:"browser_download_url"`
		}
	}{}
	if err := json.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	var urls []string
	for _, a := range r.Assets {
		urls = append(urls, a.BrowserDownloadUrl)
	}
	return urls, nil
}
