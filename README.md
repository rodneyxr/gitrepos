# gitrepos
Automatically download zipped repositories from GitHub.

# Requirements
* [Golang](https://golang.org/dl/)

# Setup
1. Place your repo list in `data/repos.txt`
2. Run `go build`
3. Run `./gitrepos` or `./gitrepos -f 'data/repos.txt' -o 'data/repos/'`
4. Zipped repos will be downloaded and saved in `data/repos/`

# Help
```
$ ./gitrepos --help

Usage of gitrepos.exe:
  -f string
        File containing a list of GitHub repositories on each line. (default "data/repos.txt")
  -m int
        Maximum repository size that you wish to download. (default 50000)
  -o string
        Directory in which to save the download repositories. (default "data/repos/")
  -p string
        Password or token for basic authentication.
  -u string
        Username for basic authentication for the GitHub REST API v3.
```
