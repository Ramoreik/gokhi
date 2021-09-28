package scraper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

const KhinsiderHost = "downloads.khinsider.com"
const KhinsiderOSTUri = "game-soundtracks/album"

var Album, DownloadDir string

/*
A download worker, it has a queue of urls to download,
It is also in a waitGroup.
<- *WaitGroup, downloadQueue
*/
func dlWorker(wg *sync.WaitGroup, downloadQueue chan string) {
	defer wg.Done()
	for u := range downloadQueue {
		downloadFile(u)
	}
}

/*
Instantiates a colly instance with the given callback and goquery.
Mostly exists to save a few lines.
I should look into how colly can be managed, with a singleton maybe ?
<- goquery, callback
-> colly collector
*/
func khrawler(goquery string, cb colly.HTMLCallback) *colly.Collector {
	// Initiate the collector
	Crawler := colly.NewCollector(
		colly.AllowedDomains(KhinsiderHost),
		colly.Async(true),
	)
	err := Crawler.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 3})
	if err != nil {
		log.Fatal(err)
	}

	Crawler.OnHTML(goquery, cb)
	return Crawler
}

/*
Searching an album page for all the song links.
<- album name: string
-> song link list: []string
*/
func findSongLinks() []string {
	var songs []string
	c := khrawler("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, ".mp3") {
			songs = append(songs, e.Request.AbsoluteURL(link))
		}
	})

	url := url.URL{
		Scheme: "https",
		Host:   KhinsiderHost,
		Path:   "game-soundtracks/album/" + url.QueryEscape(Album),
	}

	err := c.Visit(url.String())
	if err != nil {
		log.Println(err)
	}
	c.Wait()
	return songs
}

/*
find for download links within a song's page.
<- songUls: list of song page links.
-> list of download links, one for each song
*/
func findDownloadLinks(songUrls []string) []string {
	var downloadLinks []string
	c := khrawler("audio[src]", func(e *colly.HTMLElement) {
		downloadLinks = append(downloadLinks, e.Attr("src"))
	})
	for _, l := range songUrls {
		err := c.Visit(l)
		if err != nil {
			log.Prinln(err)
		}
	}
	c.Wait()
	return downloadLinks
}

/*
Prepares the download of a file by calculating its filepath and creating
the parent dir if needed.
<- u : url representing the download link of a song.
-> filepath: relative path to the download folder
*/
func prepareDownload(u string) (string, error) {
	var err error
	dir := path.Join(DownloadDir, strings.Split(u, "/")[4])
	filename := strings.Split(u, "/")[6]
	filepath := path.Join(dir, filename)
	if _, e := os.Stat(dir); os.IsNotExist(e) {
		err = os.Mkdir(dir, 0700)
	}
	return filepath, err
}

/*
Downloads a file from the given download link.
<- u : url representing the download link of a song.
*/
func downloadFile(u string) {
	filepath, err := prepareDownload(u)
	fmt.Printf("🔽 Downloading %v\n", filepath)
	// Initiate empty file
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Download the file
	res, err := http.Get(u)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// Copy reponse body (file) to the local file
	_, err = io.Copy(file, res.Body)
	if err != nil {
		log.Fatal(err)
	}
}

/*
Downloads a full album from khinsider.
<- album : string reprensenting the album name to download.

TODO: control over download directory and worker numbers.
*/
func Download(album string, downloadDir string) {
	DownloadDir = downloadDir
	Album = album

	fmt.Println("🔍 Finding songs ...")
	songUrls := findSongLinks()
	if len(songUrls) == 0 {
		fmt.Println("❌ Could not find any songs for the specified album, does it exist ?")
		os.Exit(1)
	}

	fmt.Println("🔍 Finding download links ...")
	downloadUrls := findDownloadLinks(songUrls)

	fmt.Println("🔍 Downloading found songs ...")
	downloadQueue := make(chan string, 10)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go dlWorker(&wg, downloadQueue)
	}
	for _, dl := range downloadUrls {
		dl, _ = url.QueryUnescape(dl)
		downloadQueue <- dl
	}
	close(downloadQueue)
	wg.Wait()
	fmt.Println("✅ Done, please check the specified download directory !")
}

/*
Searches the site using the specified query.
<- query : search query
-> found albums list
*/
func Search(query string) []string {
	var links []string

	c := khrawler("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, KhinsiderOSTUri) {
			links = append(links, strings.Split(link, "/")[3])
		}
	})

	// Construct a vaid url
	url := url.URL{
		Scheme:   "https",
		Host:     KhinsiderHost,
		Path:     "search",
		RawQuery: "search=" + url.QueryEscape(query),
	}

	// Initiate and then return links
	c.Visit(url.String())
	c.Wait()
	return links
}
