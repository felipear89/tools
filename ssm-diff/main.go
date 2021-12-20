package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rfyiamcool/go-shell"
)

type SSM struct {
	Parameters []Parameters `json:"Parameters"`
}

type Parameters struct {
	Name             string    `json:"Name"`
	Type             string    `json:"Type"`
	Value            string    `json:"Value"`
	Version          int       `json:"Version"`
	LastModifiedDate time.Time `json:"LastModifiedDate"`
	ARN              string    `json:"ARN"`
	DataType         string    `json:"DataType"`
}

func main() {

	leftParam := flag.String("left", "", "")
	rightParam := flag.String("right", "", "")
	flag.Parse()

	leftString := getParam(*leftParam)
	rightString := getParam(*rightParam)

	left := &SSM{}
	err := json.Unmarshal([]byte(leftString), &left)
	if err != nil {
		log.Fatal(err, " unmarshal left failed")
	}

	right := &SSM{}
	err = json.Unmarshal([]byte(rightString), &right)
	if err != nil {
		log.Fatal(err, " unmarshal right failed")
	}

	leftMap := make(map[string]Parameters)
	for _, v := range left.Parameters {
		leftMap[strings.ReplaceAll(v.Name, *leftParam, "")] = v
	}

	rightMap := make(map[string]Parameters)
	for _, v := range right.Parameters {
		rightMap[strings.ReplaceAll(v.Name, *rightParam, "")] = v
	}

	fmt.Println("\nComparing", *leftParam, "and", *rightParam)

	changed := diff(rightMap, leftMap)
	addedRight := missing(rightMap, leftMap)
	addedLeft := missing(leftMap, rightMap)

	fmt.Println("\n##### Missing on " + *leftParam + " #####")
	for _, v := range addedRight {
		fmt.Println("added:", v.Name+" = "+v.Value)
	}

	fmt.Println("\n##### Missing on " + *rightParam + " #####")
	for _, v := range addedLeft {
		fmt.Println("added:", v.Name+" = "+v.Value)
	}
	fmt.Println("\n##### Differences #####")
	for _, p := range changed {
		fmt.Println("diff:", p.Name+" = "+p.Value)
	}
}

func missing(rightMap, leftMap map[string]Parameters) map[string]Parameters {
	added := make(map[string]Parameters)
	for k, v := range rightMap {
		if _, ok := leftMap[k]; !ok {
			added[k] = v
		}
	}
	return added
}

func diff(rightMap, leftMap map[string]Parameters) []Parameters {
	changed := []Parameters{}
	for k, v := range rightMap {
		if _, ok := leftMap[k]; ok {
			if v.Value != leftMap[k].Value {
				changed = append(changed, v)
				changed = append(changed, leftMap[k])
			}
		}
	}
	return changed
}

func getParam(path string) string {
	awsCommand := "aws ssm get-parameters-by-path --path " + path + " --recursive"
	fmt.Println("$", awsCommand)
	cmd := shell.NewCommand(awsCommand)
	err := cmd.Run()
	output := cmd.Status.Output

	if err != nil || strings.Contains(output, "error") {
		log.Fatal(err, " unable to run 'aws ssm get-parameters-by-path' ", output)
	}
	return output
}
