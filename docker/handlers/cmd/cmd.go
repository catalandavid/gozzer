package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/iancoleman/strcase"
)

// DockerHandler ...
type DockerHandler interface {
	Run() error
	Stop() error
	Pull() error
}

// DockerCMD ...
type DockerCMD struct {
	Label       string
	Image       string
	Options     map[string]interface{}
	Command     string
	containerID string
}

// Run ...
func (d *DockerCMD) Run() error {
	fmt.Println("Run Docker Image")
	optsArr := []string{"docker", "run"}
	optsArr = append(optsArr, serializeOptions(d.Options)...)
	optsArr = append(optsArr, string(d.Image))
	if d.Command != "" {
		optsArr = append(optsArr, string(d.Command))
	}

	if _, err := exec.LookPath("docker"); err != nil {
		log.Panicln("docker command does not exist, cannot continue!")
		return err
	}

	cmd := exec.Command("sh", "-c", strings.Join(optsArr, " "))

	fmt.Println(cmd)

	var out bytes.Buffer
	cmd.Stdout = &out

	var stdErr bytes.Buffer
	cmd.Stderr = &stdErr

	err := cmd.Run()

	fmt.Println(stdErr.String())

	if err != nil {
		return err
	}

	cID := strings.ReplaceAll(out.String(), "\n", "")

	fmt.Printf("Container ID: %s\n", cID)

	d.containerID = cID

	return nil
}

// Stop ...
func (d DockerCMD) Stop() error {
	fmt.Println("Stop Docker Container")
	cmd := exec.Command("docker", "stop", d.containerID)

	fmt.Println(cmd)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	d.containerID = ""

	return nil
}

// Pull ...
func (d DockerCMD) Pull() error {
	fmt.Printf("Pull Docker Image %s\n", d.Image)

	cmd := exec.Command("docker", "pull", d.Image)

	fmt.Println(cmd)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println("=================================")
	fmt.Println(out.String())
	fmt.Println("=================================")
	return nil
}

// NewDockerCMD ...
func NewDockerCMD(label string, image string, options map[string]interface{}, command string) DockerHandler {
	return &DockerCMD{label, image, options, command, ""}
}

func serializeOptions(options map[string]interface{}) []string {
	res := []string{}

	for key, value := range options {
		opt := serializeOption(key, value)
		if opt != "" {
			res = append(res, opt)
		}
	}

	return res
}

func serializeOption(key string, value interface{}) string {
	// fmt.Printf("key[%s] value[%s]\n", key, value)

	oPrefix := ""
	oKey := ""
	oSep := " "

	oKey = strcase.ToKebab(string(key))

	if len(oKey) == 1 {
		oPrefix = "-"
	} else {
		oPrefix = "--"
		oSep = "="
	}

	if bValue, ok := value.(bool); ok {
		if bValue {
			return oPrefix + oKey
		} else {
			return ""
		}
	} else if aValue, ok := value.([]interface{}); ok {
		opts := []string{}
		for _, v := range aValue {
			opts = append(opts, serializeOption(key, v))
		}
		return strings.Join(opts, " ")
	} else if aValue, ok := value.([]string); ok {
		opts := []string{}
		for _, v := range aValue {
			opts = append(opts, serializeOption(key, v))
		}
		return strings.Join(opts, " ")
	}

	res := oPrefix + oKey + oSep + value.(string)
	res = strings.ReplaceAll(res, `"`, "")

	// fmt.Printf(">>>%v<<<\n", res)

	return res
}
