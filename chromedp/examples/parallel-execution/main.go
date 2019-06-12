package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/chromedp/cdproto/browser"

	"github.com/chromedp/chromedp"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func worker(id int, jobs <-chan int, results chan<- int) {
	// Delay startup
	time.Sleep(time.Duration(id*2500) * time.Millisecond)

	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		openWebSite(id)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {
	jobs := make(chan int, 1000)
	results := make(chan int, 1000)

	for w := 1; w <= 16; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= 1000; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= 1000; a++ {
		<-results
	}
}

func openWebSite(i int) {
	fmt.Printf("%d openWebSite()\n", i)
	// dir, err := ioutil.TempDir("", "chromedp-example")
	// if err != nil {
	// 	panic(err)
	// }
	// defer os.RemoveAll(dir)

	// r := rand.New(rand.NewSource(99))
	w := 100 + rand.Intn(899)
	h := 100 + rand.Intn(677)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoSandbox,
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		// chromedp.DisableGPU,
		// chromedp.UserDataDir(dir),
		chromedp.WindowSize(w, h),
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

	fmt.Printf("[%d, %d] Window Bounds: %+v\n", w, h, winSize)

	// cancel()
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
		// chromedp.Sleep(1 * time.Millisecond),
	}
}
