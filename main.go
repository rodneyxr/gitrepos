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
)

func init() {
	flag.StringVar(&RepoListFilePath, "f", "data/repos.txt", "File containing a list of GitHub repositories on each line in the format 'user/repo_name'. (default='data/repos.txt'")
	flag.StringVar(&RepoSaveDirectory, "o", "data/repos/", "Directory in which to save the download repositories. (default='data/repos/')")
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
