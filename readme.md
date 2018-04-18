# sls-go-crud

Example Serverless CRUD application using Serverless framework with Golang

## Setup
In order to deploy this sample project you need an AWS account, Serverless framework installed globablly and Golang.

To build, run the make file first.

```make``` 

Then deploy using Serverless-Cli

```serverless deploy```

---

After the deployment, copy the url from the console and replace your-url.com


### to Create an Item
```
curl -X POST \
  https://your-url.com/dev/dogs \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "Woof",
	"age": 5,
	"weight": 40,
	"race": "Wolf",
	"favfood": "Chicken"
}'
```

### to Get All the Items
```
curl -X GET \
  https://your-url.com/dev/dogs \
  -H 'Cache-Control: no-cache'
```

### to Get Single Item with Id
```
curl -X GET \
  https://your-url.com/dev/dogs/rmzMSTZiR \
  -H 'Cache-Control: no-cache'
```

### To Update Single Item
```
curl -X PUT \
  https://your-url.com/dev/dogs/rmzMSTZiR \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -d '{
	"name": "Woof",
	"age": 5,
	"weight": 40,
	"race": "Wolf",
	"favfood": "Chicken"
}'
```

### To Delete Single Item
```
curl -X DELETE \
  https://your-url.com/dev/dogs/rmzMSTZiR \
  -H 'Cache-Control: no-cache' 
```  