package aem

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	gCDP "github.com/catalandavid/gozzer/browser/chromedp"
)

// AEMBrowser ...
type AEMBrowser struct {
	*gCDP.Browser
}

// Login ...
func (b *AEMBrowser) Login(baseURL, username, password string) {
	resp, err := http.PostForm(baseURL+"/libs/granite/core/content/login.html/j_security_check",
		url.Values{"j_username": {username}, "j_password": {password}, "j_validate": {"true"}})

	if err != nil {
		fmt.Println(err)

		return
	}

	// fmt.Println("resp")
	// fmt.Println(resp)

	loginTokenHeaderValue := resp.Header.Get("Set-Cookie")
	loginToken := strings.Split(loginTokenHeaderValue, ";")[0]
	loginToken = strings.Split(loginToken, "=")[1]

	// fmt.Println(loginToken)

	b.SetCookie("login-token", loginToken, baseURL)
}

// AvoidSurvey ...
func (b *AEMBrowser) AvoidSurvey() {
	fmt.Println("Entering AvoidSurvey")

	// err := chromedp.Run(b.Ctx,
	// 	chromedp.WaitNotPresent("omg_surveyContainer", chromedp.ByID),
	// )

	// fmt.Println("AvoidSurvey Query Done")

	b.ExecJS("if (document.querySelector('#omg_surveyContainer') != null) { location.reload(); }")

	time.Sleep(1 * time.Second)

	// if err != nil {
	// 	fmt.Println("======================")
	// 	fmt.Println("OMG Survey displayed!! Reload window!!")
	// 	fmt.Println("======================")
	// 	b.ExecJS("if (document.querySelector('#omg_surveyContainer') != null) { location.reload(); }")
	// } else {
	// 	fmt.Println("OMG Survey not displayed")
	// }
}
