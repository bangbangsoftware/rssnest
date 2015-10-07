package feeds

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type GoldMoney struct {
	SpotPrices []Spot
}

type Spot struct {
	Metal     string
	Timestamp uint64
	SpotPrice float64
	Trend     string
	Units     string
}

var feedURL = "http://ws.goldmoney.com/metal/prices/currentSpotPrices?currency=gbp&units="

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
	price = append(price, getIt(silverFeed))
	return price
}
