package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	gojenkins "github.com/yosida95/golang-jenkins"
)

var jenkinsURL string
var jobNamePattern string
var buildsFolder string
var forceDelete bool
var auth *gojenkins.Auth

func init() {
	auth = &gojenkins.Auth{
		Username: os.Getenv("JENKINS_API_USERNAME"),
		ApiToken: os.Getenv("JENKINS_API_TOKEN"),
	}
}

func main() {
	// Parse CLI Flags
	flag.StringVar(&jenkinsURL, "j", "", "Root URL of the jenkins server")
	flag.StringVar(&jobNamePattern, "p", "", "String pattern to filter job names")
	flag.StringVar(&buildsFolder, "b", "", "Folder containing Jenkins builds")
	flag.BoolVar(&forceDelete, "y", false, "Confirm oprhans deletion")
	flag.Parse()

	// Init Jenkins client
	jenkins := gojenkins.NewJenkins(auth, jenkinsURL)

	// Get all Jenkins pipeline jobs
	jobs, err := jenkins.GetJobs()
	if err != nil {
		log.Fatalln(err)
	}

	for _, job := range jobs {
		if strings.Contains(job.Class, "Folder") || (jobNamePattern != "" && !strings.Contains(job.Name, jobNamePattern)) {
			continue
		}

		var jenkinsBuilds []string
		var folderBuilds []string

		// Get job details
		j, err := jenkins.GetJob(job.Name)
		if err != nil {
			log.Println(err)
			continue
		}

		// Get job builds
		for _, job2 := range j.Jobs {
			jenkinsBuilds = append(jenkinsBuilds, job2.Name)
		}

		// Get list of build folders
		files, err := ioutil.ReadDir(buildsFolder + "/" + job.Name)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, f := range files {
			folderBuilds = append(folderBuilds, f.Name())
		}

		// Get orphans build folders (compared to Jenkins builds)
		orphanBuilds := difference(folderBuilds, jenkinsBuilds)

		log.Printf("Found %d orphan(s) for %s:\n", len(orphanBuilds), job.Name)
		log.Println(orphanBuilds)

		if forceDelete == true {
			for _, orphan := range orphanBuilds {
				out, err := exec.Command("du", "-h", "-d", "1", buildsFolder+"/"+job.Name+"/"+orphan).Output()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("Folder size is %s", out)
			}
		}
	}
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
