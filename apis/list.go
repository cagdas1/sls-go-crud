package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Dog struct {
	Id      string  `json:"Id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Weight  float32 `json:"weight"`
	Race    string  `json:"race"`
	FavFood string  `json:"favfood"`
}

type Config struct {
	Region string
	Table  string
}

func (c *Config) Load() error {
	flag.StringVar(&c.Table, "table", "Dogs", "Table to Query on")
	flag.StringVar(&c.Region, "region", "us-west-1", "AWS Region the table is in")
	flag.Parse()
	if len(c.Table) == 0 {
		flag.PrintDefaults()
		return fmt.Errorf("table name is required.")
	}
	return nil
}

func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	panic(err)
}

var dynamo *dynamodb.DynamoDB
var config Config

func setup() {
	config = Config{}
	if len(config.Table) == 0 {
		if err := config.Load(); err != nil {
			exitWithError(err)
		}
	}
	awsConfig := &aws.Config{}
	awsConfig.WithRegion(config.Region)
	sess := session.Must(session.NewSession(awsConfig))
	dynamo = dynamodb.New(sess)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if dynamo == nil || len(config.Table) == 0 {
		setup()
	}

	params := &dynamodb.ScanInput{
		TableName: aws.String(config.Table),
	}
	result, dbErr := dynamo.Scan(params)
	if dbErr != nil {
		fmt.Printf("dbError %s", dbErr.Error())
		exitWithError(dbErr)
	}
	dogs := []Dog{}
	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &dogs)
	if err != nil {
		exitWithError(err)
	}
	body, _ := json.Marshal(dogs)
	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
