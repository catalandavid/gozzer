package chromedp

import (
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// GetWebPageScreenshot ...
func GetWebPageScreenshot(url string, width int, height int, loadWait time.Duration) []byte {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var buf []byte
	err := chromedp.Run(ctx, screenshot(url, width, height, loadWait, &buf))
	if err != nil {
		log.Fatal(err)
	}

	// // save the screenshot to disk
	// if err = ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
	// 	log.Fatal(err)
	// }

	return buf
}

// GetWebPagePNGScreenshot ...
func GetWebPagePNGScreenshot(url string, width int, height int, loadWait time.Duration) image.Image {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var buf []byte
	err := chromedp.Run(ctx, screenshot(url, width, height, loadWait, &buf))
	if err != nil {
		log.Fatal(err)
	}

	i, err := png.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(buf))))
	if err != nil {
		fmt.Println(err)
	}
	// // save the screenshot to disk
	// if err = ioutil.WriteFile("screenshot.png", buf, 0644); err != nil {
	// 	log.Fatal(err)
	// }

	return i
}

func screenshot(url string, width int, height int, loadWait time.Duration, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		// chromedp.Navigate(url),

		chromedp.ActionFunc(func(ctx context.Context) error {

			success, err := network.SetCookie("login-token", "c102bd1e-d454-4dc5-88d9-0ff51140ed78%3a0122e965-4ece-4f36-9422-48e71a944466_4cc06b026a257eab1d888d8f355e5181%3acrx.default").WithDomain("localhost").Do(ctx)
			if err != nil && success != true {
				fmt.Println("SetCookie")
				fmt.Println(err)
			}

			ID, _, err := browser.GetWindowForTarget().Do(ctx)
			if err != nil {
				fmt.Println(err)
			}

			bounds := browser.Bounds{
				Left:        0,
				Top:         0,
				Width:       int64(width),
				Height:      int64(height),
				WindowState: browser.WindowStateNormal,
			}

			err = browser.SetWindowBounds(ID, &bounds).Do(ctx)
			if err != nil {
				fmt.Println(err)
			}

			// *res = *b
			// // fmt.Printf("%+v\n", res)
			// return nil

			// // _, b,
			// // _, b, err := browser.SetWindowBounds(ID, ).Do(ctx)
			// if err != nil {
			// 	fmt.Println(err)
			// }
			// *res = *b
			// // fmt.Printf("%+v\n", res)
			return nil
		}),
		chromedp.Navigate(url),

		chromedp.Sleep(loadWait),

		// chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.CaptureScreenshot(res),
	}
}
