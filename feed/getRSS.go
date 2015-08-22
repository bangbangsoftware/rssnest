package feed

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Enclosure struct {
	Url string `xml:"url,attr"`
}

type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Enclosure   Enclosure `xml:"enclosure"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

func GetRSS(feedURL string) *Rss {

	response, err := http.Get(feedURL)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer response.Body.Close()

	XMLdata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	buffer := bytes.NewBuffer(XMLdata)
	decoded := xml.NewDecoder(buffer)

	rss := new(Rss)
	err = decoded.Decode(rss)
	if err != nil {
		log.Println(err)
		return nil
	}

	return rss
}
