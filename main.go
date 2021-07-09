package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var zoneDirs = []string{
	"/usr/share/zoneinfo/",
	"/usr/share/lib/zoneinfo/",
	"/usr/lib/locale/TZ/",
}
var zoneDir string

var Timezones []string

type GroupObject struct {
	GroupId      string   `json:"groupId"`
	Date         string   `json:"date"`
	TimezoneName string   `json:"timezoneName"`
	UserIds      []string `json:"userIds"`
}

func main() {
	lambda.Start(sqsHandler)
}

func sqsHandler(event events.SQSEvent) {
	var g GroupObject
	for _, record := range event.Records {
		str := record.Body
		_ = json.Unmarshal([]byte(str), &g)
		fmt.Printf("The body was \n %+v\n", g)
	}
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Incoming request is: %#v\n", req)
	switch req.HTTPMethod {
	case "GET":
		return FetchUser(req)
	case "POST":
		return CreateUser(req)
	case "PUT":
		return UpdateUser(req)
	case "DELETE":
		return DeleteUser(req)
	default:
		log.Printf("Requested for unhandled method %#v", req.HTTPMethod)
		response := events.APIGatewayProxyResponse{StatusCode: 405}
		return response, nil
	}
}

func printTimezones(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	for _, zoneDir = range zoneDirs {
		ReadFile("")
	}
	for i, zone := range Timezones {
		fmt.Println(i, zone)
	}
	zones, _ := json.Marshal(Timezones)
	response := events.APIGatewayProxyResponse{StatusCode: 200, Body: string(zones)}
	return response, nil
}

func ReadFile(path string) {
	files, _ := ioutil.ReadDir(zoneDir + path)
	for _, f := range files {
		if f.Name() != strings.ToUpper(f.Name()[:1])+f.Name()[1:] {
			continue
		}
		if f.IsDir() {
			ReadFile(path + "/" + f.Name())
		} else {
			Timezones = append(Timezones, (path + "/" + f.Name())[1:])
		}
	}
}
