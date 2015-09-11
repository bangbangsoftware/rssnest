package feed

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
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

type ShortURL struct {
	Kind    string
	Id      string
	LongURL string
	Err     error
}

type HTTPTraffic interface {
	GetRSS(feedURL string) *Rss
	DetectContentType(data []byte) string
	Get(url string) (resp *http.Response, err error)
	Shorten(feedURL string, apiKey string) *ShortURL
}

type StdTraffic struct {
}

func (h StdTraffic) DetectContentType(data []byte) string {
	return http.DetectContentType(data)
}

func (h StdTraffic) Get(url string) (resp *http.Response, err error) {
	return http.Get(url)
}

func (h StdTraffic) Shorten(feedURL string, apiKey string) *ShortURL {
	short := new(ShortURL)
	//curl https://www.googleapis.com/urlshortener/v1/url?key=blar -H 'Content-Type: application/json' -d '{"longUrl": "http://superuser.com/questions/149329/what-is-the-curl-command-line-syntax-to-do-a-post-request"}'
	client := &http.Client{}
	url := "https://www.googleapis.com/urlshortener/v1/url?key=" + apiKey
	var jsonStr = []byte(`{"longUrl":"` + feedURL + `"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		log.Println("Cannot shorten, %v\n", err)
		short.Err = err
		return short
	}
	defer response.Body.Close()
	/*
				var body = response.Body
				log.Println("short return is, %v\n", body)

		//		data, err := ioutil.ReadAll(body)
				if err != nil {
					short.Err = err
					log.Println("Cannot read content from shortener, %v\n", err)
					return short
				}
				json.Unmarshal([]byte(data), &short)
	*/

	//ioutil.ReadAll starts at a very small 512
	//it really should let you specify an initial size
	buffer := bytes.NewBuffer(make([]byte, 0, 65536))
	io.Copy(buffer, response.Body)
	temp := buffer.Bytes()
	length := len(temp)
	var data []byte
	//are we wasting more than 10% space?
	if cap(temp) > (length + length/10) {
		data = make([]byte, length)
		copy(data, temp)
	} else {
		data = temp
	}
	json.Unmarshal(data, &short)

	return short

}

func (h StdTraffic) GetRSS(feedURL string) *Rss {

	r, err := http.Get(feedURL)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer r.Body.Close()

	//ioutil.ReadAll starts at a very small 512
	//it really should let you specify an initial size
	buffer := bytes.NewBuffer(make([]byte, 0, 65536))
	io.Copy(buffer, r.Body)
	temp := buffer.Bytes()
	length := len(temp)
	var XMLdata []byte
	//are we wasting more than 10% space?
	if cap(temp) > (length + length/10) {
		XMLdata = make([]byte, length)
		copy(XMLdata, temp)
	} else {
		XMLdata = temp
	}
	//	XMLdata, err := ioutil.ReadAll(response.Body)
	//	if err != nil {
	//		log.Println(err)
	//		return nil
	//	}

	buf := bytes.NewBuffer(XMLdata)
	decoded := xml.NewDecoder(buf)

	rss := new(Rss)
	err = decoded.Decode(rss)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rss
}
