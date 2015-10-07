package main

import (
	"github.com/bangbangsoftware/config"
	"github.com/bangbangsoftware/feeds"
	"log"
)

func Propagate(store feeds.Store, target Target, shortener Shortener) {
	log.Printf("\nPropagate the rss feeds results...")
	conf := config.GetConfig()
	qty := conf.Propagate.QtyPerPage
	filename := conf.General.DataDir + conf.General.StoreName
	log.Printf("Propagate last %v results...", qty)
	var newItems = store.GetLast(qty, filename)
	var prices = feeds.GetPrices()
	target.Send(newItems, prices)
	sendList := store.GetToMessage()
	for i := range sendList {
		msg := sendList[i].Item.Title
		url := sendList[i].Link
		target.Message(new(GoogleShort), msg, url)
	}
}
