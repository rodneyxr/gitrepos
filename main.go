package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Command line arguments
var (
	RepoListFilePath  string
	RepoSaveDirectory string
	Username          string
	Token             string
	MaxRepoSize       int
)

func init() {
	flag.StringVar(&RepoListFilePath, "f", "data/repos.txt", "File containing a list of GitHub repositories on each line.")
	flag.StringVar(&RepoSaveDirectory, "o", "data/repos/", "Directory in which to save the download repositories.")
	flag.StringVar(&Username, "u", "", "Username for basic authentication for the GitHub REST API v3.")
	flag.StringVar(&Token, "p", "", "Password or token for basic authentication.")
	flag.IntVar(&MaxRepoSize, "m", 50000, "Maximum repository size that you wish to download.")
	flag.Parse()

	RepoSaveDirectory = strings.Replace(RepoSaveDirectory, `\`, `/`, -1)
	if !strings.HasSuffix(RepoSaveDirectory, "/") {
		RepoSaveDirectory += "/"
	}
}

func main() {
	// Open the repo list file
	file, err := os.Open(RepoListFilePath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer file.Close()

	// For each line in the repo list file download the archive of the repo
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Construct URL to download from
		// https://github.com/{user}/{repo_name}/archive/master.zip
		repoName := scanner.Text()
		url := `https://github.com/` + repoName + `/archive/master.zip`
		fmt.Println(url)

		// Get info from the repo to check the size
		repo := RepoInfo{}
		if err := getJSON(`https://api.github.com/repos/`+repoName, &repo); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to download repo: %s\n", repoName)
			continue
		}

		// If repo size is too large then log and skip
		if repo.Size > MaxRepoSize {
			fmt.Printf("Size of %dKb too large; skipping '%s'\n", repo.Size, repoName)
			continue
		}

		// Create directory to save our repo downloads to
		err := os.MkdirAll(RepoSaveDirectory, 0777)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// Get the zip file via HTTP GET request
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
			continue
		}
		defer resp.Body.Close()

		// Create zip file to save reponse to
		out, err := os.Create(RepoSaveDirectory + strings.Replace(repoName, "/", "_", -1) + ".zip")
		if err != nil {
			log.Fatal(err)
			continue
		}
		defer out.Close()

		// Save the response to the zip file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatal(err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
