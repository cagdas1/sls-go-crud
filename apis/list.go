package main

import (
	"fmt"
	"encoding/json"
	"flag"
	"os"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Dog struct {
	DogId string	 `json:"DogId"`
	Name string      `json:"Name"`
	Age int          `json:"Age"`
	Weight float32   `json:"Weight"`
	Race string 	 `json:"Race"`
	FavFood string 	 `json:"Favfood"`
}

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

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Init dynamodb
	config := Config{}
	if err:= config.Load(); err != nil{
		exitWithError(err)
	}
	awsConfig := &aws.Config{}
	awsConfig.WithRegion(config.Region)
	sess := session.Must(session.NewSession(awsConfig))
	dynamo := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String(config.Table),
	}
	result, dbErr := dynamo.Scan(params)
	if dbErr != nil{
		fmt.Printf("dbError %s", dbErr.Error())
		exitWithError(dbErr)
	}
	dogs := []Dog{}
	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &dogs)
	if err != nil{
		exitWithError(err)
	}
	body, _ := json.Marshal(dogs)
	return events.APIGatewayProxyResponse{
		Body: string(body),
		StatusCode: 200,
	}, nil
}
			
func main() {
	lambda.Start(Handler)
}