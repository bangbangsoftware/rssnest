package feed

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var gotAlready map[string]Item = make(map[string]Item)

func loadItems() {
	log.Printf("About to load alreadyHave.json\n")
	file, e := ioutil.ReadFile("./alreadyHave.json")
	if e != nil {
		if strings.Contains(e.Error(), "no such file") {
			log.Printf("No alreadyHave.json file, so creating one\n")
			var data []byte
			var perm os.FileMode = 0777
			ioutil.WriteFile("./alreadyHave.json", data, perm)
		} else {
			log.Printf("File error: %v\n", e)

			os.Exit(1)
		}
	}
	json.Unmarshal(file, &gotAlready)
	log.Printf("already have list is: %n \n", len(gotAlready))

}

func Have(item Item, link string) {
	log.Printf("saving alreadyHave file\n\n")
	gotAlready[link] = item
	jsave, _ := json.Marshal(gotAlready)
	var data []byte = jsave
	var perm os.FileMode = 0777
	ioutil.WriteFile("./alreadyHave.json", data, perm)
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
