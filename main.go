package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
Traffic collections

1. Android: /Android-Mischief-Dataset/AndroidMischiefDataset_v2/
2. Malware: /
3. Normal: /
4. IoT: /IoT-23-Dataset-v2/
*/

const BASE_URL = "https://mcfp.felk.cvut.cz/publicDatasets"

type URL_INFO struct {
	url    string
	prefix string
}

var DATA_SOURCES = map[string]URL_INFO{
	"android": URL_INFO{"Android-Mischief-Dataset/AndroidMischiefDataset_v2/", "RAT"},
	"malware": URL_INFO{"", "CTU-Malware-Capture-Botnet"},
	"normal":  URL_INFO{"", "CTU-Normal"},
	"iot":     URL_INFO{"IoT-23-Dataset-v2/", "CTU-"},
}

func main() {
	/*
		Usage

		$ go run main.go -source=[Data Collection Type] -save_path=[Directory PATH]

		- [Data Collection Type]: Choose one among `android`|`malware`|`normal`|`iot`
		- [Directory PATH]: Directory where the downloaded resource will be placed in
	*/
	source := flag.String("source", "malware", "Choose one among `android`|`malware`|`normal`|`iot`")
	save_path := flag.String("save_path", ".", "Directory where the downloaded resource will be placed in")

	flag.Parse()

	source_info, ok := DATA_SOURCES[*source]

	if !ok {
		log.Fatal("[Error] Choose one of the following data sources: `android`, `malware`, `normal`, `iot`\n")
	}

	source_url, _ := url.JoinPath(BASE_URL, source_info.url)

	base_dir := filepath.Join(*save_path, *source)

	DownloadResource(source_url, base_dir, source_info.prefix, 0)
}

func DownloadResource(url_path string, dir_path string, prefix string, depth int) {
	if depth == 0 {
		fmt.Printf(".%s\n", dir_path)
	} else {
		fmt.Printf("%s├── %s\n", strings.Repeat("    ", depth-1), dir_path[strings.LastIndex(dir_path, "/")+1:])
	}

	resp, err := http.Get(url_path)

	if err != nil {
		log.Fatalf("[Error] Base URL (%s) is not accessible \n%s", url_path, err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		log.Fatalf("[Error] Failed to read a html document (path: %s)\n%s", url_path, err)
	}

	if _, err := os.Stat(dir_path); os.IsNotExist(err) {
		// directory for a certain category does not exist
		if err := os.Mkdir(dir_path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	doc.Find("table td a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")

		if strings.HasSuffix(url_path, link) {
			return
		}

		if !strings.HasPrefix(link, prefix) {
			return
		}

		if s.Text() == "Parent Directory" {
			return
		}

		material_url, _ := url.JoinPath(url_path, link)

		if strings.HasSuffix(material_url, "/") {
			subdir_path := filepath.Join(dir_path, link)
			DownloadResource(material_url, subdir_path, "", depth+1)
		} else {
			file_path := filepath.Join(dir_path, link)
			DownloadFile(material_url, file_path, depth)
		}
	})
}

// Download from a given url to a file.
func DownloadFile(url_path string, file_path string, depth int) {
	fmt.Printf("%s└── %s\n", strings.Repeat("    ", depth), file_path[strings.LastIndex(file_path, "/")+1:])

	if _, err := os.Stat(file_path); err == nil {
		return
	}

	out, err := os.Create(file_path)

	if err != nil {
		log.Fatalf("[Error] Cannot create file %s\n%s", file_path, err)
	}

	resp, err := http.Get(url_path)

	if err != nil {
		log.Fatalf("[Error] Failed to read a html document (path: %s)\n%s", url_path, err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()
}
