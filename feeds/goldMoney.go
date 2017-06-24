package feeds

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GoldMoney struct {
  Rates Rate
}

type Rate struct {
  Spots  []Spot
}

type Spot struct {
  quoteCurrency string
  baseCurrency string
  bid float64
  avg float64
  ask float64
}

//var feedURL = "http://ws.goldmoney.com/metal/prices/currentSpotPrices?currency=gbp&units="
//var feedURL = "https://wealth.goldmoney.com/api/prices/currentSpotPrices/?currency=gbp&price=ask&units="
var feedURL = "https://wealth-api.goldmoney.com/public/markets/summary"

func getIt(url string) GoldMoney {
	var price GoldMoney
	response, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return price
	}
	defer response.Body.Close()
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
	json.Unmarshal(data, &price)

	return price
}

func GetPrices() []GoldMoney {
	var price []GoldMoney
	goldFeed := feedURL + "ounces"
	silverFeed := feedURL + "grams"
	price = append(price, getIt(goldFeed))
	log.Printf("Gold : %s\n", price)
	price = append(price, getIt(silverFeed))
	log.Printf("Silver : %s\n", price)
	return price
}
