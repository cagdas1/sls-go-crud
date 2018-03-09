package main

import (
	"fmt"
	"flag"
	"os"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)


type Config struct{
	Region string
	Table string
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

func setup(){
	config = Config{}
	if len(config.Table) == 0{
		if err:= config.Load(); err != nil{
			exitWithError(err)
		}
	}
	awsConfig := &aws.Config{}
	awsConfig.WithRegion(config.Region)
	sess := session.Must(session.NewSession(awsConfig))
	dynamo = dynamodb.New(sess)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if dynamo == nil || len(config.Table) == 0{
		setup()
	}
	id := request.PathParameters["id"]

	params := &dynamodb.DeleteItemInput{
		TableName: aws.String(config.Table),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {
				S: aws.String(id),
			},
		},
	}
	result, dbErr := dynamo.DeleteItem(params)
	if dbErr != nil{
		fmt.Printf("dbError %s", dbErr.Error())
		exitWithError(dbErr)
	}
	return events.APIGatewayProxyResponse{
		Body: result.String(),
		StatusCode: 200,
	}, nil
}
			
func main() {
	lambda.Start(Handler)
}