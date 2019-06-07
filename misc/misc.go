package misc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// GetIntsFromInput ...
func GetIntsFromInput(maxInputs int) []int {
	// Variables
	var input string
	var inputs []int

	for i := 0; i < maxInputs; i++ {
		// Prompt user to enter an int and wait for the input
		fmt.Printf("Please enter an int ('x' to interrupt sequence of inputs and continue program) %02d/%02d: ", i+1, maxInputs)
		fmt.Scan(&input)

		// Handle "quit" criteria
		if input == "x" {
			break
		}

		if intInput, err := strconv.Atoi(input); err == nil {
			inputs = append(inputs, intInput)
		} else {
			fmt.Println("Cannot evaluate your input as an int, please try again.")
			i--
		}
	}

	return inputs
}

// LoadJSONFromFile ...
func LoadJSONFromFile(path string, target interface{}) error {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	// var data interface{}

	err = json.Unmarshal([]byte(file), target)

	if err != nil {
		return err
	}

	return nil
}

// Config ...
type Config struct {
	*viper.Viper
}

// LoadConfig ...
func LoadConfig(filePath string) (Config, error) {
	v, err := readConfig(filePath, map[string]interface{}{
		"parallelWorkers":  1,
		"workingDirectory": ".",
		"itemsFile":        "items.txt"})

	if err != nil {
		return Config{}, err
	}

	// TODO
	// v.WatchConfig()
	// v.OnConfigChange(func(e fsnotify.Event) {
	// 	log.Infof("Config file changed: %s", e.Name)
	// })

	return Config{v}, nil
}

// PollCheckUntil ...
func PollCheckUntil(checkFn func() bool, timeout time.Duration) error {
	interval := timeout / 20

	// Init Ticker with computed interval
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var checkDuration time.Duration

	for {
		select {
		case <-ticker.C:
			res := checkFn()

			if res == true {
				return nil
			} else if checkDuration >= timeout {
				return errors.New("Reached timeout")
			} else {
				checkDuration += interval
			}
		}
	}
}

func getCurrentPath() string {
	ex, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return ex
}

func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()

	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	v.SetConfigFile(filename)
	err := v.ReadInConfig()
	return v, err
}
