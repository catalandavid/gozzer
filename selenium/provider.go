package selenium

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Provider ...
type Provider int

// Provider
const (
	Docker Provider = iota
	Saucelabs
	BrowserStack
	DockerCompose
)

var providerNames = []string{"docker", "saucelabs", "browserstack", "docker-compose"}

// var providerIDs = []Provider{Docker, Saucelabs, BrowserStack, DockerCompose}

var providerIDs = map[string]Provider{
	"docker":         Docker,
	"saucelabs":      Saucelabs,
	"browserstack":   BrowserStack,
	"docker-compose": DockerCompose,
}

// ParseProvider ...
func ParseProvider(b string) Provider {
	idx := indexOf(b, providerNames)
	return Provider(idx)
}

// String returns the name of the day
func (p Provider) String() string {
	return providerNames[p]
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

// MarshalJSON ...
func (p *Provider) MarshalJSON() ([]byte, error) {
	fmt.Println("---------------------")
	fmt.Println("---------------------")
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(p.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON ...
func (p *Provider) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*p = providerIDs[j]
	return nil
}
