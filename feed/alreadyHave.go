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
	Save(item Item, link string, dir string)
}

type FileStore struct {
	//	gotAlready map[string]Item make(map[string]Item)
	GotAlready map[string]Item
	JustGot    []Item
	Prices     [][]GoldMoney

	Perm os.FileMode
}

func (fs FileStore) loadItems(dir string) {
	log.Printf("About to load alreadyHave.json\n")
	filename := dir + "alreadyHave.json"
	file, e := ioutil.ReadFile(filename)
	var data []byte
	fs.Perm = 0777
	if fs.GotAlready == nil {
		fs.GotAlready = make(map[string]Item)
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
	data = []byte("var data = \n")
	json.Unmarshal(file, &fs.GotAlready)
	log.Printf("already have list is: %v \n", len(fs.GotAlready))
	ioutil.WriteFile(dir+"newData.json", data, fs.Perm)
}

func (fs FileStore) Save(item Item, link string, dir string) {
	log.Printf("saving alreadyHave file\n")
	fs.GotAlready[link] = item
	jsave, _ := json.Marshal(fs.GotAlready)
	var data []byte = jsave
	ioutil.WriteFile(dir+"alreadyHave.json", data, fs.Perm)

	log.Printf("saving newData file to ftp to webserver\n")
	fs.JustGot = append(fs.JustGot, item)
	jsave, _ = json.Marshal(fs.JustGot)
	var data2 []byte = jsave
	ioutil.WriteFile(dir+"newData.json", data2, fs.Perm)

	log.Printf("saving prices file to ftp to webserver\n\n\n")
	spot := GetPrices()
	log.Printf("Spot is %v", spot)
	fs.Prices = append(fs.Prices, spot)
	jsave, _ = json.Marshal(fs.Prices)
	var data3 []byte = jsave
	ioutil.WriteFile(dir+"prices.json", data3, fs.Perm)
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
