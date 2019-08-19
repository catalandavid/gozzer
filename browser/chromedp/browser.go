package chromedp

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"time"

	"github.com/chromedp/cdproto/runtime"

	"github.com/catalandavid/gozzer/misc"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Browser ...
type Browser struct {
	Ctx               context.Context
	NavigationOccured bool
	cancelFunc        context.CancelFunc
}

// New ...
func New(headless bool) *Browser {
	b := &Browser{}

	initCtx := context.Background()

	if headless == false {
		opts := []chromedp.ExecAllocatorOption{
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			// chromedp.Headless,
			chromedp.DisableGPU,
			// chromedp.UserDataDir(dir),
		}

		ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
		initCtx = ctx
	}

	// Init Browser
	ctx, cancel := chromedp.NewContext(initCtx)

	b.Ctx = ctx
	b.cancelFunc = cancel

	b.ListenTarget(func(ev interface{}) {
		switch ev.(type) {
		case *page.EventLoadEventFired:
			// out, _ := json.Marshal(ev)
			b.NavigationOccured = true
			// fmt.Printf("Navigate: %T %s\n", ev, string(out))
		case *page.EventFrameNavigated:
			// out, _ := json.Marshal(ev)
			b.NavigationOccured = true
			// fmt.Printf("Navigate: %T %s\n", ev, string(out))
		}
	})

	b.NavigationOccured = false

	return b
}

// Close ...
func (b *Browser) Close() {
	b.cancelFunc()
}

// ListenTarget ...
func (b *Browser) ListenTarget(fn func(interface{})) {
	chromedp.ListenTarget(b.Ctx, fn)
}

// SetWindowSize ...
func (b *Browser) SetWindowSize(width int, height int) {
	err := chromedp.Run(b.Ctx, func(w int, h int) chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.ActionFunc(func(ctx context.Context) error {
				ID, _, err := browser.GetWindowForTarget().Do(ctx)
				if err != nil {
					fmt.Println(err)
				}

				bounds := browser.Bounds{
					Left:        0,
					Top:         0,
					Width:       int64(w),
					Height:      int64(h),
					WindowState: browser.WindowStateNormal,
				}

				err = browser.SetWindowBounds(ID, &bounds).Do(ctx)
				if err != nil {
					fmt.Println(err)
				}

				return nil
			}),
		}

	}(width, height))
	if err != nil {
		log.Fatal(err)
	}
}

// Navigate ...
func (b Browser) Navigate(url string) {
	b.NavigationOccured = false

	_ = chromedp.Run(b.Ctx, func(url string) chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.Navigate(url),
		}
	}(url))

	misc.PollCheckUntil(func() bool { return b.NavigationOccured }, 3*time.Second)
}

// SetCookie ...
func (b *Browser) SetCookie(name string, value string, url string) {
	_ = chromedp.Run(b.Ctx, func(name string, value string, url string) chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.ActionFunc(func(ctx context.Context) error {
				success, err := network.SetCookie(name, value).WithURL(url).Do(ctx)
				if err != nil && success != true {
					fmt.Println("SetCookie")
					fmt.Println(err)

					return err
				}

				return nil
			}),
		}
	}(name, value, url))
}

// GetNodeForLocation ...
func (b *Browser) GetNodeForLocation(x int, y int) (cdp.BackendNodeID, cdp.NodeID) {
	var (
		backendNodeID cdp.BackendNodeID
		nodeID        cdp.NodeID
		err           error
	)

	_ = chromedp.Run(b.Ctx, func(x int, y int) chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.ActionFunc(func(ctx context.Context) error {
				backendNodeID, nodeID, err = dom.GetNodeForLocation(int64(x), int64(y)).Do(ctx)
				if err != nil {
					fmt.Println("GetNodeForLocation")
					fmt.Println(err)
				}

				return nil
			}),
		}
	}(x, y))

	return backendNodeID, nodeID
}

// GetContentSize ...
func (b *Browser) GetContentSize() (page.VisualViewport, dom.Rect) {
	var (
		// layoutViewport *page.LayoutViewport
		visualViewport *page.VisualViewport
		contentSize    *dom.Rect
		err            error
	)

	_ = chromedp.Run(b.Ctx, func() chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.ActionFunc(func(ctx context.Context) error {

				_, visualViewport, contentSize, err = page.GetLayoutMetrics().Do(ctx)
				if err != nil {
					fmt.Println("GetNodeForLocation")
					fmt.Println(err)
				}

				return nil
			}),
		}
	}())

	return *visualViewport, *contentSize
}

// ExecJS ...
func (b *Browser) ExecJS(js string) []byte {
	// fmt.Println(js)
	var res []byte

	err := chromedp.Run(b.Ctx, chromedp.EvaluateAsDevTools(js, &res))
	if err != nil {
		fmt.Println(err)
	}

	// json.MarshalJSON()
	// jsonByte, _ := res.MarshalJSON()

	// var out bytes.Buffer

	// _ = json.Indent(&out, jsonByte, "", "\t")

	// fmt.Println(out.String())
	// fmt.Println(string(res))
	return res
}

// InjectJS ...
func (b *Browser) InjectJS(js string) {
	_ = chromedp.Run(b.Ctx, func(js string) chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.ActionFunc(func(ctx context.Context) error {
				r1, r2, err := runtime.Evaluate(js).WithThrowOnSideEffect(false).WithSilent(true).WithIncludeCommandLineAPI(true).Do(ctx)

				fmt.Println("===================")
				fmt.Println(r1)
				fmt.Println("===================")
				fmt.Println(r2)
				fmt.Println("===================")
				fmt.Println(err)
				fmt.Println("===================")
				return nil
			}),
		}
	}(js))
}

// MouseClickXY ...
func (b *Browser) MouseClickXY(x, y int64) {
	chromedp.Run(b.Ctx, chromedp.MouseClickXY(x, y))
}

// CapturePNGScreenshot ...
func (b *Browser) CapturePNGScreenshot() image.Image {
	var buf []byte

	err := chromedp.Run(b.Ctx, func(res *[]byte) chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.CaptureScreenshot(res),
		}
	}(&buf))

	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(buf)
	png, err := png.Decode(r)
	if err != nil {
		panic("Bad png")
	}

	// png, err := png.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(buf))))
	// if err != nil {
	// 	fmt.Println(err)
	// }

	return png
}

// CapturePNGBytesBufferScreenshot ...
func (b *Browser) CapturePNGBytesBufferScreenshot() []byte {
	var buf []byte

	err := chromedp.Run(b.Ctx, func(res *[]byte) chromedp.Tasks {
		return chromedp.Tasks{
			chromedp.CaptureScreenshot(res),
		}
	}(&buf))

	if err != nil {
		log.Fatal(err)
	}

	return buf
}

// GetURL ...
func (b *Browser) GetURL() string {
	url := b.ExecJS("window.location.href")

	urlString := string(url)

	return urlString[1 : len(urlString)-1]
}
