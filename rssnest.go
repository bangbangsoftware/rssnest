// rssnest project
package main

import (
	"flag"
	"fmt"
	"github.com/bangbangsoftware/config"
	"github.com/bangbangsoftware/feeds"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func setup() {
	//Switches
	var memoryPort = flag.Int("mem", -1, "Turn memory listener on this port")
	var configFile = flag.String("conf", "./conf.json", "The path to the configuration file")
	var sendAssets = flag.Bool("ftp", false, "FTP the assets to the webserver")
	flag.Parse()

	// memory
	if *memoryPort > -1 {
		go func() {
			var listen = fmt.Sprintf("localhost:%v", memoryPort)
			log.Println(http.ListenAndServe(listen, nil))
		}()
	}

	// Config
	log.Printf("loading config from: %s \n", *configFile)
	config.LoadConfig(*configFile)

	// Ftp Assets
	if *sendAssets {
		log.Printf("Ftping assets to webserver")
		//ftpAssets(config)
	}
	return
}

func main() {
	// Logging
	f, err := os.OpenFile("rssnest.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	setup()

	var store = new(FileStore)
	var source = new(HttpSource)
	feeds.StoreNewContent(store, source)

	var target = new(FtpTarget)
	var short = new(GoogleShort)
	Propagate(store, target, short)
}
