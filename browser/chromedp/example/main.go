package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/catalandavid/gozzer/browser/chromedp"
)

func main() {
	b := chromedp.New(false)

	b.SetWindowSize(666, 777)

	b.Navigate("http://golang.org")

	screenshot := b.CapturePNGScreenshot()

	f, err := os.Create("screen.png")
	if err != nil {
		fmt.Println(err)
	}

	png.Encode(f, screenshot)

	b.Close()
}
