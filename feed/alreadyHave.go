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

var perm os.FileMode = 0777

func loadItems() {
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
	ioutil.WriteFile("./newData.json", data, perm)

}

func Have(item Item, link string) {
	log.Printf("saving alreadyHave file\n\n")
	gotAlready[link] = item
	jsave, _ := json.Marshal(gotAlready)
	var data []byte = jsave
	ioutil.WriteFile("./alreadyHave.json", data, perm)
	justGot = append(justGot, item)
	jsave, _ = json.Marshal(justGot)
	var data2 []byte = jsave
	ioutil.WriteFile("./newData.json", data2, perm)
}

func AlreadyHave(itemLink string) bool {

	if len(gotAlready) == 0 {
		loadItems()
	}
	if _, ok := gotAlready[itemLink]; ok {
		return true
	}
	return false
}
