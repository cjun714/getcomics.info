package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"getcomics.info/model"
	"github.com/PuerkitoBio/goquery"
	"github.com/cjun714/glog/log"
)

var siteURL = "https://getcomics.info/"
var baseURL = "https://getcomics.info/page/"
var baseDir = "z:/"

func main() {
	// start, e := strconv.Atoi(os.Args[1])
	// if e != nil {
	// 	panic(e)
	// }
	// end, e := strconv.Atoi(os.Args[2])
	// if e != nil {
	// 	panic(e)
	// }

	e := downloadAll(baseURL, 1, 2)
	// e := downloadAll(baseURL, start, end)
	if e != nil {
		panic(e)
	}
}

func downloadAll(baseURL string, start, end int) error {
	for i := start; i <= end; i++ {
		url := baseURL + strconv.Itoa(i)
		byts, e := downloadHTML(url)
		if e != nil {
			log.E("access failed:", url, "error:", e)
			continue
		}

		infos, e := parseIndex(byts)
		if e != nil {
			log.E("parse failed:", url, "error:", e)
			continue
		}

		for _, info := range infos {
			cover := filepath.Base(info.Cover)
			idx := strings.LastIndex(cover, "?")
			if idx != -1 {
				cover = cover[:idx]
			}
			e = downloadImage(info.Cover, filepath.Join(baseDir, cover))
			if e != nil {
				log.E("download failed:", info.Cover, "error:", e)
				continue
			}
			info.Cover = cover
			log.I(info.Name)
			log.I(info.Cover)
		}

	}
	return nil
}

func parseIndexFile(path string) ([]model.ComicInfo, error) {
	byts, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}
	_, e = parseIndex(byts)
	if e != nil {
		return nil, e
	}

	return nil, nil
}

func parseIndex(byts []byte) ([]model.ComicInfo, error) {
	infos := make([]model.ComicInfo, 0, 12)
	r := bytes.NewReader(byts)
	doc, e := goquery.NewDocumentFromReader(r)
	if e != nil {
		return nil, e
	}
	doc.Find(".post-list-posts article").Each(func(i int, s *goquery.Selection) {
		var info model.ComicInfo
		node := s.Find(".post-header-image img").First()
		info.Cover, _ = node.Attr("src")

		node = s.Find(".post-header-image a").First()
		info.PageURL, _ = node.Attr("href")

		node = s.Find(".post-category").First()
		info.Category = node.Text()

		node = s.Find(".post-title").First()
		name := node.Text()
		info.Name = strings.Replace(name, "\n", "", -1)

		node = s.Find("p")
		str := node.Text()
		str = strings.Replace(str, "\n", "", -1)
		str = strings.Trim(str, " ")
		info.StartYear, info.EndYear, info.Size = getYearSize(str)

		infos = append(infos, info)
	})

	return infos, nil
}

func downloadHTML(url string) ([]byte, error) {
	res, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("get html faild,resp code:" + strconv.Itoa(res.StatusCode))
	}

	return ioutil.ReadAll(res.Body)
}

func downloadHTML2(url string, targetPath string) error {
	// get html into bytes[]
	res, e := http.Get(url)
	if e != nil {
		return e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("get html faild,resp code:" + strconv.Itoa(res.StatusCode))
	}
	bs, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return e
	}

	// save html into file
	f, e := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if e != nil {
		return e
	}
	defer f.Close()

	wr := bufio.NewWriter(f)
	_, e = wr.Write(bs)
	if e != nil {
		return e
	}

	return nil
}

// Year : 2016 | Size : 602 MB
// Year : 2016-2018 | Size : 602 MB
func getYearSize(str string) (int, int, int) {
	start, end, size := 0, 0, 0
	n, _ := fmt.Sscanf(str, "Year : %d-%d | Size : %d MB", &start, &end, &size)
	if n != 3 {
		fmt.Sscanf(str, "Year : %d | Size : %d MB", &start, &size)
		end = start
	}

	return start, end, size
}

func downloadImage(url string, targetPath string) error {
	resp, e := http.Get(url)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	f, e := os.Create(targetPath)
	if e != nil {
		return e
	}
	defer f.Close()

	_, e = io.Copy(f, resp.Body)
	if e != nil {
		return e
	}

	return nil
}
