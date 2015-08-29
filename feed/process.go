package feed

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bangbangsoftware/config"
)

func processContent(body io.ReadCloser, name string, traffic HTTPTraffic) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println("Cannot read content from rss, %v\n", err)
		return
	}
	var t = traffic.DetectContentType(data)
	log.Printf("%s which is %s content type \n", name, t)
	if strings.Contains(t, "text/html") {
		log.Println("Rss content is html, nothing to download\n")
		return
	}

	if !strings.Contains(name, ".") {
		log.Println("Rss content has no extention, defaulting to mp3\n")
		name = name + ".mp3"
	}

	out, err := os.Create(name)
	defer out.Close()
	if err != nil {
		log.Println("Error creating file for content: %v\n", err)
		return
	}

	buffer := bytes.NewBuffer(data)
	n, err := io.Copy(out, buffer)
	if err != nil {
		log.Println("Error downloading content: %v\n", err)
		return
	}
	log.Println("Downloaded %v bytes", n)
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

func checkAndGet(i int, item Item, store Persist, traffic HTTPTraffic) {
	conf := config.GetConfig()
	dataDir := conf.General.DataDir
	var link = item.Link
	if len(item.Enclosure.Url) > 0 {
		log.Printf("using enclosure url")
		link = item.Enclosure.Url
	}
	log.Printf("[%d] %s link is: %s\n", i, item.Title, link)
	if len(link) == 0 {
		log.Printf("[%d] NO LINK!!? \n", i, item.Title, link)
		//item.Error = "Cant find any link"
		return
	}
	if store.AlreadyHave(link, dataDir) {
		log.Printf("Already have %s\n", link)
		return
	}

	response, err := traffic.Get(link)
	if err != nil {
		log.Printf("Link (%v) error: %v\n", link, err)
		return
	}

	var name = getName(link)
	defer response.Body.Close()
	processContent(response.Body, name, traffic)

	store.Save(item, link, dataDir)
}

var fs Persist = FileStore{itemMap, nil, nil, 0777}
var itemMap = make(map[string]Item)

func Process(url string) {
	processAndPersist(url, new(StdTraffic), fs)
}

func processAndPersist(url string, traffic HTTPTraffic, persist Persist) {
	rss := traffic.GetRSS(url)
	if rss == nil {
		return
	}

	log.Printf("Title : %s\n", rss.Channel.Title)
	log.Printf("Description : %s\n", rss.Channel.Description)
	log.Printf("Link : %s\n", rss.Channel.Link)

	total := len(rss.Channel.Items)

	log.Printf("Total items : %v\n", total)
	if total == 0 {
		log.Printf("No items, all is done here\n")
		return
	}
	i := 0
	//for i := -1; i < total; i++ {
	checkAndGet(i, rss.Channel.Items[i], persist, traffic)
	//}
}
