package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/bangbangsoftware/config"
	"github.com/bangbangsoftware/feed"
	"github.com/dutchcoders/goftp"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

type Misc struct {
	User      string
	CreatedOn string
}

type Item struct {
	Name  string
	Desc  string
	Url   string
	Date  string
	Error string
	Dir   string
}

type CastsTest struct {
	Misc Misc
}

type Casts struct {
	Misc  Misc
	Items []Item
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	f, err := os.OpenFile("rssnest.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	var configFile = flag.String("conf", "./conf.json", "The path to the configuration file")
	flag.Parse()
	log.Printf("loading config from: %s \n", *configFile)
	config := config.LoadConfig(*configFile)

	log.Printf("loading rss list from: %s \n", config.General.Feedfile)
	castsFile, e2 := ioutil.ReadFile(config.General.Feedfile)
	if e2 != nil {
		log.Printf("RSS list file error: %v\n", e2)
		os.Exit(1)
	}
	var cs Casts
	json.Unmarshal(castsFile, &cs)
	log.Printf("%s feed created on %s \n", cs.Misc.User, cs.Misc.CreatedOn)

	total := len(cs.Items)
	log.Printf("%v feeds \n", total)

	prices := feed.GetPrices()
	var newItems []feed.RssResult
	for i := 0; i < len(cs.Items); i++ {
		item := cs.Items[i]
		name := item.Name
		//		log.Printf("%s (%s) is described as '%s' and is at %s \n", item.Name, item.Date, item.Desc, item.Url)
		log.Printf("=================================================================\n")
		log.Printf("[%v/%v] %s (%s) \n", i, total, name, item.Date)
		items := feed.Process(name, item.Url, item.Dir, 1)

		log.Printf("%v more items found ", len(items))
		for _, e := range items {
			title := e.Item.Title
			if !e.Failed && !e.AlreadyHave {
				tweet(title, e.Link, config)
			}
			newItems = append(newItems, e)
		}
		log.Printf("%v total items found ", len(newItems))

		log.Printf("\n")
	}
	if len(newItems) == 0 {
		log.Printf("No new data overall\n")
		return
	}
	lastOnes := feed.GetLast(config.Propagate.QtyPerPage)
	saveAndFtp(lastOnes, prices, config)
}

func saveAndFtp(newItems []feed.RssResult, prices []feed.GoldMoney, config config.Settings) {
	var perm os.FileMode = 0777
	var err error
	var ftp *goftp.FTP

	log.Printf("saving newData (%v) file to ftp to %v\n", len(newItems), config.Propagate.Ftp.Url)
	var newData = []byte("var data = \n")
	jsave, _ := json.Marshal(newItems)
	for _, e := range jsave {
		newData = append(newData, e)
	}
	ioutil.WriteFile(config.General.DataDir+"newData.json", newData, perm)
	var file *os.File
	if file, err = os.Open(config.General.DataDir + "newData.json"); err != nil {
		panic(err)
	}
	defer file.Close()
	addr := fmt.Sprintf("%v:%v", config.Propagate.Ftp.Url, config.Propagate.Ftp.Port)
	if ftp, err = goftp.Connect(addr); err != nil {
		panic(err)
	}
	defer ftp.Close()
	if err = ftp.Login(config.Propagate.Ftp.Usr, config.Propagate.Ftp.Pw); err != nil {
		panic(err)
	}
	if err := ftp.Stor("newData.json", file); err != nil {
		panic(err)
	}

	log.Printf("saving prices (%v) file to ftp to webserver\n\n\n", len(prices))
	var data3 = []byte("var prices = \n")
	jsave2, _ := json.Marshal(prices)
	for _, e := range jsave2 {
		data3 = append(data3, e)
	}
	ioutil.WriteFile(config.General.DataDir+"prices.json", data3, perm)
	var pricesFile *os.File
	if pricesFile, err = os.Open(config.General.DataDir + "prices.json"); err != nil {
		panic(err)
	}
	defer pricesFile.Close()
	if err := ftp.Stor("prices.json", pricesFile); err != nil {
		panic(err)
	}

}

func tweet(msg string, u string, config config.Settings) {
	anaconda.SetConsumerKey(config.Propagate.Tweet.ConsumerKey)
	anaconda.SetConsumerSecret(config.Propagate.Tweet.ConsumerSecret)
	api := anaconda.NewTwitterApi(config.Propagate.Tweet.AccessTokenKey, config.Propagate.Tweet.AccessTokenSecret)
	//link, err := url.Parse(u)
	full := u + " " + msg
	if len(full) > 140 {
		full = full[0:137] + "..."
	}
	log.Printf("Full message to tweet is '%v'\n", full)
	api.PostTweet(full, nil) // "https://twitter.com/rssnest")

}
