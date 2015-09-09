package feed

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

func processContent(body io.ReadCloser, name string, traffic HTTPTraffic, r RssResult) RssResult {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		r.Failed = true
		r.FailReason = "Cannot read content from rss"
		log.Println("Cannot read content from rss, %v\n", err)
		return r
	}
	var t = traffic.DetectContentType(data)
	log.Printf("%s which is %s content type \n", name, t)
	if strings.Contains(t, "text/html") {
		log.Println("Rss content is html, nothing to download\n")
		r.Message = "Rss content his html, nothing to download"
		return r
	}

	if !strings.Contains(name, ".") {
		log.Println("Rss content has no extention, defaulting to mp3\n")
		name = name + ".mp3"
		r.Message = "Rss content had no extention, defaulted to mp3"
	}

	out, err := os.Create(name)
	defer out.Close()
	if err != nil {
		log.Println("Error creating file for content: %v\n", err)
		r.Failed = true
		r.FailReason = "Error creating file for content"
		return r
	}

	buffer := bytes.NewBuffer(data)
	n, err := io.Copy(out, buffer)
	if err != nil {
		log.Println("Error downloading content: %v\n", err)
		r.Failed = true
		r.FailReason = "Error downloading content"
		return r
	}
	log.Println("Downloaded %v bytes", n)
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
		//item.Error = "Cant find any link"
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

	var name = getName(link)
	defer response.Body.Close()
	result := processContent(response.Body, name, traffic, r)
	store.Save(item, link, dataDir)
	return result
}

var fs Persist = FileStore{itemMap, nil, nil, 0777}
var itemMap = make(map[string]Item)

func Process(url string, howmany int) []RssResult {
	return processAndPersist(url, new(StdTraffic), fs, howmany)
}

func processAndPersist(url string, traffic HTTPTraffic, persist Persist, howmany int) []RssResult {
	items := make([]RssResult, howmany)
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
