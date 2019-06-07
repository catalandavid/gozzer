package http

import (
	"io/ioutil"
	"net/http"

	"github.com/valyala/fastjson"
)

// GetJSON ...
func GetJSON(url string) (fastjson.Value, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
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

	return *v, nil
}
