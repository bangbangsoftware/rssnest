package feed

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var gotAlready map[string]Item = make(map[string]Item)
var justGot []Item
var prices [][]GoldMoney

var perm os.FileMode = 0777

func loadItems(dir string) {
	log.Printf("About to load alreadyHave.json\n")
	file, e := ioutil.ReadFile("./alreadyHave.json")
	var data []byte

	if e != nil {
		if strings.Contains(e.Error(), "no such file") {
			log.Printf("No alreadyHave.json file, so creating one\n")
			ioutil.WriteFile("./alreadyHave.json", data, perm)
		} else {
			log.Printf("File error: %v\n", e)

			os.Exit(1)
		}
	}
	json.Unmarshal(file, &gotAlready)
	log.Printf("already have list is: %n \n", len(gotAlready))
	ioutil.WriteFile(dir+"newData.json", data, perm)

}

func Have(item Item, link string, dir string) {
	log.Printf("saving alreadyHave file\n\n")
	gotAlready[link] = item
	jsave, _ := json.Marshal(gotAlready)
	var data []byte = jsave
	ioutil.WriteFile(dir+"alreadyHave.json", data, perm)

	justGot = append(justGot, item)
	jsave, _ = json.Marshal(justGot)
	var data2 []byte = jsave
	ioutil.WriteFile(dir+"newData.json", data2, perm)

	spot := GetPrices()
	log.Printf("Spot is %v", spot)
	prices = append(prices, spot)
	jsave, _ = json.Marshal(prices)
	var data3 []byte = jsave
	ioutil.WriteFile(dir+"prices.json", data3, perm)

}

func AlreadyHave(itemLink string, dir string) bool {

	if len(gotAlready) == 0 {
		loadItems(dir)
	}
	if _, ok := gotAlready[itemLink]; ok {
		return true
	}
	return false
}
