package main

import (
	"encoding/json"
	"github.com/bangbangsoftware/feed"
	"io/ioutil"
	"log"
	//"net/http"
	//	_ "net/http/pprof"
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
}

type CastsTest struct {
	Misc Misc
}

type Casts struct {
	Misc  Misc
	Items []Item
}

type GeneralConf struct {
	Sleep     int
	Feedfile  string
	AudioDir  string
	VisualDir string
}

type FtpConf struct {
	Url string
	Usr string
	Pw  string
}

type TweetConf struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessTokenKey    string
	AccessTokenSecret string
}

type PropergateConf struct {
	IncludePrice bool
	Ftp          FtpConf
	Tweet        TweetConf
}

type Conf struct {
	General    GeneralConf
	Propergate PropergateConf
}

//type ItemConf struct {
//
//}

//type Config struct {
//        General  GeneralConf
//        Items []ItemConf
//}
// usage: %s start|stop|restart

// 0. Load config
// 1. set timeout to  minutes in config
// 2. parse opml
// 3. foreach
// 4.   get date
// 5.   parse feed
// 6.   get url
// 7.   do we already have url?
// 8.     compile html
// 9.     download
//10.     put in correct dir based on content type
//11.     gold and silver price
//12.     compile html
//13.     ftp html
//14.     tweet

func main() {
	//log.Println(http.ListenAndServe(":6060", nil))
	file, e := ioutil.ReadFile("./conf.json")
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	var config Conf
	json.Unmarshal(file, &config)

	log.Printf("loading rss list from: %s \n", config.General.Feedfile)
	castsFile, e2 := ioutil.ReadFile(config.General.Feedfile)
	if e2 != nil {
		log.Printf("File error: %v\n", e2)
		os.Exit(1)
	}
	var cs Casts
	json.Unmarshal(castsFile, &cs)
	log.Printf("%s feed created on %s \n", cs.Misc.User, cs.Misc.CreatedOn)

	log.Printf("%v feeds \n", len(cs.Items))
	for i := 0; i < len(cs.Items); i++ {
		item := cs.Items[i]
		//		log.Printf("%s (%s) is described as '%s' and is at %s \n", item.Name, item.Date, item.Desc, item.Url)
		log.Printf("=================================================================\n")
		log.Printf("%s (%s) \n", item.Name, item.Date)
		feed.Process(item.Url, config.General.AudioDir, config.General.VisualDir)
	}

}
