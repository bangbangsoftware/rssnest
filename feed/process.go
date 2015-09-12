package feed

import (
	//	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bangbangsoftware/config"
)

type RssResult struct {
	Name        string
	Date        time.Time
	Item        Item
	Link        string
	AlreadyHave bool
	Failed      bool
	FailReason  string
	Message     string
}

func processContent(t string, body io.ReadCloser, name string, traffic HTTPTraffic, r RssResult) RssResult {
	defer body.Close()
	log.Printf("%s which is %s content type \n", name, t)
	if strings.Contains(t, "text/html") || strings.Contains(t, "text/xml") {
		r.Message = "Rss content type is text, nothing to download"
		log.Println(r.Message + "\n")
		return r
	}

	if !strings.Contains(name, ".") {
		name = name + ".mp3"
		r.Message = "Rss content had no extention, defaulted to mp3"
		log.Println(r.Message + "\n")
	}

	out, err := os.Create(name)
	defer out.Close()
	if err != nil {
		r.Failed = true
		r.FailReason = "Error creating file for content"
		log.Println(r.FailReason + "\n")
		return r
	}

	n, err := io.Copy(out, body)
	if err != nil {
		r.Failed = true
		r.FailReason = "Error copying content to file"
		log.Println(r.FailReason + "\n")
		return r
	}
	log.Println("Downloaded ", n, " bytes")
	out = nil

	r.Message = fmt.Sprintf("%v Downloaded %v bytes", r.Message, n)
	return r
}

func getName(link string) string {
	conf := config.GetConfig()
	audioDir := conf.General.AudioDir
	visualDir := conf.General.VisualDir
	var bits = strings.Split(link, "/")
	var name = bits[len(bits)-1]
	name = strings.Split(name, "?")[0]
	var dir = audioDir
	if strings.HasSuffix(name, "mp4") || strings.HasSuffix(name, "mv4") {
		log.Println("Visual content\n")
		dir = visualDir
	}
	return dir + name
}

func checkAndGet(i int, item Item, store Persist, traffic HTTPTraffic) RssResult {
	conf := config.GetConfig()
	dataDir := conf.General.DataDir
	var r RssResult
	r.Item = item
	r.Date = time.Now()
	var link = item.Link
	if len(item.Enclosure.Url) > 0 {
		log.Printf("using enclosure url")
		link = item.Enclosure.Url
	}
	log.Printf("[%d] %s link is: %s\n", i, item.Title, link)
	if len(link) == 0 {
		log.Printf("[%d] NO LINK!!? \n", i, item.Title, link)
		r.Failed = true
		r.FailReason = "No Link?"
		return r
	}
	short := traffic.Shorten(link, conf.Propagate.Apikey)
	if short.Err == nil {
		log.Printf("%v shortened to %v", link, short.Id)
		r.Link = short.Id
	} else {
		log.Printf("Failed shortening")
		log.Printf("%v", short.Err)
	}

	if store.AlreadyHave(link, dataDir) {
		log.Printf("Already have %s\n", link)
		r.AlreadyHave = true
		return r
	}

	response, err := traffic.Get(link)
	if err != nil {
		log.Printf("Link (%v) error: %v\n", link, err)
		r.Failed = true
		r.FailReason = "Link error"
		return r
	}

	name := getName(link)
	defer response.Body.Close()
	log.Printf("response (%v) \n", response.Header)
	t := response.Header.Get("Content-Type")
	result := processContent(t, response.Body, name, traffic, r)
	store.Save(item, link, dataDir)
	return result
}

var fs Persist = FileStore{itemMap, nil, nil, 0777}
var itemMap = make(map[string]Item)

func Process(url string, howmany int) []RssResult {
	return processAndPersist(url, new(StdTraffic), fs, howmany)
}

func processAndPersist(url string, traffic HTTPTraffic, persist Persist, howmany int) []RssResult {
	var items []RssResult
	rss := traffic.GetRSS(url)
	if rss == nil {
		log.Printf("No feed, all is done here\n")
		return items
	}

	log.Printf("Title : %s\n", rss.Channel.Title)
	log.Printf("Description : %s\n", rss.Channel.Description)
	log.Printf("Link : %s\n", rss.Channel.Link)

	total := len(rss.Channel.Items)

	log.Printf("Total items : %v\n", total)
	if total == 0 {
		log.Printf("No items, all is done here\n")
		return items
	}
	for i := 0; i < howmany; i++ {
		if i < len(rss.Channel.Items) {
			item := rss.Channel.Items[i]
			result := checkAndGet(i, item, persist, traffic)
			if result.AlreadyHave {
				log.Printf("got it already....")
			} else {
				log.Printf("Item found\n")
				items = append(items, result)
			}
		}
	}
	return items
}
