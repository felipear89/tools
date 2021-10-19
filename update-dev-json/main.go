package main

import (
	"encoding/json"
	"fmt"
	shell "github.com/rfyiamcool/go-shell"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	checkAwsEnvVariable("echo $AWS_PROFILE")
	checkAwsEnvVariable("echo $AWS_DEFAULT_REGION")
	kid := getActiveKid()
	token := getActiveKidToken(kid)
	path, err := getDevFilePath()
	if err != nil {
		log.Fatal(err)
	}
	d := getDevJson(path)

	l := d.Get("media-store-secrets")
	lMap := l.(*OrderedMap)
	lMap.Set("active-kid",  kid)
	lMap.Set("localKeyValue",  token)

	file, _ := json.MarshalIndent(d, "", "  ")
	fmt.Println("Writing on file ", path)

	err = ioutil.WriteFile(path, file, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func getDevJson(path string) *OrderedMap {

	var err error
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	orderedMap := NewOrderedMap()
	err = orderedMap.UnmarshalJSON(file)
	if err != nil {
		log.Fatal(err)
	}
	return orderedMap
}

func checkAwsEnvVariable(awsCommand string) {
	fmt.Println(awsCommand)
	cmd := shell.NewCommand(awsCommand)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
	output := cmd.Status.Output
	output = strings.TrimSuffix(output, "\n")
	if len(output) == 0{
		log.Fatal("$AWS_PROFILE and $AWS_DEFAULT_REGION must be set")
	}
	fmt.Println("Value:", output)
}


func getActiveKid() string {
	awsCommand := "aws ssm get-parameter --name \"/production/media-store-secrets/active-kid\" --output json"
	fmt.Println(awsCommand)
	cmd := shell.NewCommand(awsCommand +
		" | grep Value | sed 's/Value//' | sed 's/\"//g' | sed 's/://' | sed 's/,//' | sed 's/ //g'")
	err := cmd.Run()
	output := cmd.Status.Output
	output = strings.TrimSuffix(output, "\n")
	if err != nil || strings.Contains(output, "error"){
		log.Fatal(err, output)
	}
	fmt.Println("value =", output)
	return output
}

func getActiveKidToken(kid string) string {
	awsCommand := "aws ssm get-parameter --name \"/production/media-store-secrets/%s\" --output json"
	fmt.Println(awsCommand)
	cmd := shell.NewCommand(fmt.Sprintf(awsCommand+
		" | grep Value | sed 's/Value//' | sed 's/\"//g' | sed 's/://' | sed 's/,//' | sed 's/ //g'", kid))
	err := cmd.Run()
	output := cmd.Status.Output
	output = strings.TrimSuffix(output, "\n")
	if err != nil || strings.Contains(output, "error") {
		log.Fatal(err, output)
	}
	fmt.Println("value =", output)
	return output
}

func getDevFilePath() (string, error) {
	if _, err := os.Stat("./dev.json"); err == nil {
		return "./dev.json", nil
	} else if _, err := os.Stat("./config/dev.json"); err == nil {
		return "./config/dev.json", nil
	} else if _, err := os.Stat("./service/config/dev.json"); err == nil {
		return "./service/config/dev.json", nil
	}
	return "", fmt.Errorf("file dev.json not found")
}


