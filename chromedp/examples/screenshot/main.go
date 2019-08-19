package main

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/catalandavid/gozzer/chromedp"
)

func main() {
	buf := chromedp.GetWebPageScreenshot("http://golang.org", 1200, 666, (3 * time.Second))

	// save the screenshot to disk
	if err := ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
		log.Fatal(err)
	}

}
