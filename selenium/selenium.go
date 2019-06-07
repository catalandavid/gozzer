package selenium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mssola/user_agent"
	"github.com/valyala/fastjson"

	"github.com/catalandavid/gozzer/browser"
	gHttp "github.com/catalandavid/gozzer/http"
	"github.com/catalandavid/gozzer/misc"

	"github.com/tebeka/selenium"
)

// ServerModes ...
type ServerModes uint

// Server Modes
const (
	UNKNOWN ServerModes = iota
	STANDALONE
	HUB
)

var serverModes = []string{"unknown", "standalone", "hub"}

// String returns the name of the day
func (m ServerModes) String() string {
	return serverModes[m]
}

// ServerDetails ...
type ServerDetails struct {
	Mode        ServerModes
	Version     string
	JavaVersion string
}

// BrowserDetails ...
type BrowserDetails struct {
	Name    string
	Version string
}

// PollCheckServerIsReady ...
func PollCheckServerIsReady(serverURL string, t time.Duration) error {
	err := misc.PollCheckUntil(func() bool {
		fmt.Println("poll")

		json, err := gHttp.GetJSON(serverURL)
		if err != nil {
			return false
		}

		// fmt.Println(json)

		return json.Get("value").GetBool("ready")
	}, t)

	return err
}

// GetServerDetails ...
func GetServerDetails(serverURL string) (ServerDetails, error) {
	json, err := gHttp.GetJSON(serverURL)
	if err != nil {
		return ServerDetails{}, err
	}

	// fmt.Println(json)

	message := string(json.Get("value").GetStringBytes("message"))
	message = strings.ToLower(message)

	mode := UNKNOWN

	if strings.Contains(message, "hub") {
		mode = HUB
	} else if strings.Contains(message, "server") {
		mode = STANDALONE
	}

	return ServerDetails{
		mode,
		string(json.Get("value").Get("build").GetStringBytes("version")),
		string(json.Get("value").Get("java").GetStringBytes("version")),
	}, nil
}

// GetBrowserSessionDetailsAsJSON ...
func GetBrowserSessionDetailsAsJSON(serverURL string, browser browser.Family) (fastjson.Value, error) {
	// fmt.Println(browser.String())

	// caps := map[string]interface{}{
	// 	"desiredCapabilities": map[string]interface{}{
	// 		"browserName": browser.String(),
	// 	}}

	// data, err := json.Marshal(caps)
	// if err != nil {
	// 	return fastjson.Value{}, err
	// }

	// resp, err := http.Post(serverURL+"/session", "application/json", bytes.NewBuffer(data))
	// if err != nil {
	// 	return fastjson.Value{}, err
	// }
	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return fastjson.Value{}, err
	// }

	// var p fastjson.Parser

	// v, err := p.Parse(string(body))
	// if err != nil {
	// 	return fastjson.Value{}, err
	// }

	// fmt.Println(v)

	sessionJSON, err := CreateSession(serverURL, browser, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return *sessionJSON.Get("value"), nil
}

// caps := fmt.Sprintf(`{ "desiredCapabilities": { "browserName": "%s" }}`, browser)

// fmt.Println(caps)

// bytesRepresentation, err := json.Marshal(caps)
// if err != nil {
// 	log.Fatalln(err)
// }

// j, err := json.Marshal(caps)
// if err != nil {
// 	return fastjson.Value{}, err
// }

// buf := new(bytes.Buffer)
// json.NewEncoder(buf).Encode(caps)

// caps := selenium.Capabilities{"browserName": browser}

// CreateSession ...
func CreateSession(serverURL string, browser browser.Family, capabilities interface{}) (fastjson.Value, error) {
	// fmt.Println(browser.String())

	caps := map[string]interface{}{
		"desiredCapabilities": map[string]interface{}{
			"browserName": browser.String(),
		}}

	data, err := json.Marshal(caps)
	if err != nil {
		return fastjson.Value{}, err
	}

	resp, err := http.Post(serverURL+"/session", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fastjson.Value{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fastjson.Value{}, err
	}

	var p fastjson.Parser

	v, err := p.Parse(string(body))
	if err != nil {
		return fastjson.Value{}, err
	}

	// fmt.Println(v)
	return *v, nil
}

// GetBrowserDetailsFromUserAgent ...
func GetBrowserDetailsFromUserAgent(serverURL string, browser browser.Family) (BrowserDetails, error) {
	caps := selenium.Capabilities{"browserName": browser.String()}
	wd, err := selenium.NewRemote(caps, serverURL)
	if err != nil {
		return BrowserDetails{}, err
	}
	defer wd.Quit()

	uaString, err := wd.ExecuteScript("return window.navigator.userAgent;", nil)
	if err != nil {
		return BrowserDetails{}, err
	}

	// fmt.Println(uaString)
	ua := user_agent.New(uaString.(string))

	// fmt.Println(ua)
	// fmt.Println(ua.Browser())

	name, version := ua.Browser()
	// fmt.Printf("%v\n", name)
	// fmt.Printf("%v\n", version)

	return BrowserDetails{name, version}, nil
}
