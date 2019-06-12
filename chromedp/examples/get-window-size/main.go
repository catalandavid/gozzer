package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/browser"

	"github.com/chromedp/chromedp"
)

func main() {
	// dir, err := ioutil.TempDir("", "chromedp-example")
	// if err != nil {
	// 	panic(err)
	// }
	// defer os.RemoveAll(dir)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
		// chromedp.UserDataDir(dir),
		chromedp.WindowSize(999, 777),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// run task list
	var winSize browser.Bounds
	err := chromedp.Run(taskCtx, getWindowSize(`https://www.golang.org/`, `#footer`, &winSize))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Window Bounds: %+v\n", winSize)
}

func getWindowSize(urlstr, sel string, res *browser.Bounds) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, b, err := browser.GetWindowForTarget().Do(ctx)
			if err != nil {
				fmt.Println(err)
			}
			*res = *b
			// fmt.Printf("%+v\n", res)
			return nil
		}),
		chromedp.Sleep(10 * time.Second),
	}
}
