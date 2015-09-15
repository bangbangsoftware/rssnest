package feed

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Persist interface {
	AlreadyHave(itemLink string, dir string) bool
	Save(item RssResult, link string, dir string)
	GetLast(i int) []RssResult
}

type FileStore struct {
	GotAlready map[string]RssResult
	Perm       os.FileMode
}

func (fs FileStore) loadItems(dir string) {
	log.Printf("About to load alreadyHave.json\n")
	filename := dir + "alreadyHave.json"
	file, e := ioutil.ReadFile(filename)
	var data []byte
	fs.Perm = 0777
	if fs.GotAlready == nil {
		fs.GotAlready = make(map[string]RssResult)
	}

	if e != nil {
		if strings.Contains(e.Error(), "no such file") {
			log.Printf("No %v file, so creating one\n", filename)
			ioutil.WriteFile(filename, data, fs.Perm)
		} else {
			log.Printf("File error: %v\n", e)
			os.Exit(1)
		}
	}
	json.Unmarshal(file, &fs.GotAlready)
	log.Printf("already have list is: %v \n", len(fs.GotAlready))
}

func (fs FileStore) Save(rssResult RssResult, link string, dir string) {
	log.Printf("saving alreadyHave file\n")
	fs.GotAlready[link] = rssResult
	jsave, _ := json.Marshal(fs.GotAlready)
	var data []byte = jsave
	ioutil.WriteFile(dir+"alreadyHave.json", data, fs.Perm)
}

func reverseMap(m map[string]RssResult, size int) []RssResult {
	n := make([]RssResult, size)
	for _, v := range m {
		n = append(n, v)
		if len(n) == size {
			return n
		}
	}
	return n
}

func (fs FileStore) GetLast(i int) []RssResult {
	return reverseMap(fs.GotAlready, i)
}

func (fs FileStore) AlreadyHave(itemLink string, dir string) bool {
	if len(fs.GotAlready) == 0 {
		fs.loadItems(dir)
	} else {
		log.Printf("No loading needed GotAlready list is: %v \n", len(fs.GotAlready))
	}
	if _, ok := fs.GotAlready[itemLink]; ok {
		return true
	}
	return false
}
