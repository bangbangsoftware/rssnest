package feed

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func processContent(body io.ReadCloser, name string) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
		return
	}
	var t = http.DetectContentType(data)
	log.Printf("%s which is %s content type \n", name, t)
	if strings.Contains(t, "text/html") {
		return
	}

	if !strings.Contains(name, ".") {
		name = name + ".mp3"
	}

	out, err := os.Create(name)
	defer out.Close()
	if err != nil {
		log.Println(err)
		return
	}

	buffer := bytes.NewBuffer(data)
	n, err := io.Copy(out, buffer)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(n)
}

func checkAndGet(i int, item Item, audioDir string, visualDir string) {
	var link = item.Link
	if len(item.Enclosure.Url) > 0 {
		link = item.Enclosure.Url
	}
	log.Printf("[%d] %s link is: %s\n", i, item.Title, link)

	if AlreadyHave(link) {
		log.Printf("Already have %s\n", link)
		return
	}
	if len(link) > 0 {
		var bits = strings.Split(link, "/")
		var name = bits[len(bits)-1]
		name = strings.Split(name, "?")[0]
		var dir = audioDir
		if strings.HasSuffix(name, "mp4") || strings.HasSuffix(name, "mv4") {
			dir = visualDir
		}
		response, err := http.Get(link)
		if err != nil {
			log.Printf("File error: %v\n", err)
			return
		} else {
			defer response.Body.Close()
			processContent(response.Body, dir+name)
		}
	} else {
		log.Printf("[%d] NO LINK!!? \n", i, item.Title, link)
		//item.Error = "Cant find any link"
	}
	Have(item, link)
}

func Process(url string, audioDir string, visDir string) {

	rss := GetRSS(url)
	if rss == nil {
		return
	}

	log.Printf("Title : %s\n", rss.Channel.Title)
	log.Printf("Description : %s\n", rss.Channel.Description)
	log.Printf("Link : %s\n", rss.Channel.Link)

	total := len(rss.Channel.Items)

	log.Printf("Total items : %v\n", total)
	spot := GetPrices()
	log.Printf("Spot is %v", spot)

	if total == 0 {
		return
	}
	for i := 0; i < total; i++ {
		checkAndGet(i, rss.Channel.Items[i], audioDir, visDir)
	}
}
