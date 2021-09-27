package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"nemesis.sh/go-khinsider/khinsider"
)

// FN to search for albums
func search(query string) {
	fmt.Printf("🔍 Query: %v \n", query)
	links := khinsider.Search(query)
	fmt.Println("-- -- -- -- -- -- -- --")
	for i, l := range links {
		fmt.Printf("[%d] ~+> %v \n", i, l)
	}
	fmt.Println("-- -- -- -- -- -- -- --")
	fmt.Println("[?] use the 'download' subcommand to download an album.")
}

// FN to download a specific album
func download(album string, downloadPath string) {
	fmt.Printf("🔍 Album to download: %v \n", album)
	khinsider.Download(album, downloadPath)
	os.Exit(0)
}

// FN for flag management
func flags() {
	downloadCmd := flag.NewFlagSet("download", flag.ExitOnError)
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	// Download SubCMD
	album := downloadCmd.String("album", "", "Specify the album to download. [default: nil]")
	downloadPath := downloadCmd.String("download-path", "downloads", "Specify the download directory. [default: ./downloads]")

	// Search SubCMD
	query := searchCmd.String("query", "", "Specify the query to search through khinsider's albums.")

	// check if enough arguments are passed
	if len(os.Args) < 2 {
		fmt.Println("[!] Please specify the 'search' or 'download' subcommands.")
	}

	// subcommands handling
	switch os.Args[1] {
	case "search":
		searchCmd.Parse(os.Args[2:])
		if len(*query) != 0 {
			search(*query)
		} else {
			fmt.Println("[!] Please specify a string value for the 'query' parameter.")
			os.Exit(1)
		}

	case "download":
		downloadCmd.Parse(os.Args[2:])
		fmt.Println(*album)
		if len(*album) != 0 && len(*downloadPath) != 0 {
			if _, err := os.Stat(*downloadPath); os.IsNotExist(err) {
				err = os.Mkdir(*downloadPath, 0700)
				if err != nil {
					log.Fatal(err)
				}
			}
			download(*album, *downloadPath)
		} else {
			fmt.Println("[!] Please specify a string value for the 'album' and 'download-path' parameter.")
			os.Exit(1)
		}

	default:
		fmt.Println("[!] Expected 'search' or 'download' subcommands.")
		os.Exit(1)
	}
}

// FN main
func main() {
	flags()
}
