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
	
	"github.com/teris-io/shortid"
)

type Dog struct {
	Id string	`json:"Id"`
	Name string	`json:"Name"`
	Age int	`json:"Age"`
	Weight float32	`json:"Weight"`
	Race string	`json:"Race"`
	FavFood string	`json:"Favfood"`
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
	
	id, _ := shortid.Generate()
	data := &Dog{
		Id: id,
	}
	json.Unmarshal([]byte(request.Body), data)
	dynamoItem, _ := dynamodbattribute.MarshalMap(data)
	fmt.Printf("dynamoItem %s", dynamoItem)
	params := &dynamodb.PutItemInput{
		Item: dynamoItem,
		TableName: aws.String(config.Table),
	}
	_, errDb := dynamo.PutItem(params);
	if errDb!= nil{
		fmt.Printf("errDb %s", errDb.Error())
		return events.APIGatewayProxyResponse{
			Body: string(errDb.Error()),
			StatusCode: 500,
			}, nil
			}else{
				body, _ := json.Marshal(data)
				return events.APIGatewayProxyResponse{
					Body: string(body),
					StatusCode: 200,
					}, nil
				}
			}
			
			func main() {
				lambda.Start(Handler)
			}