package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	log.Printf("Gin starting")

	r := gin.Default()
	
	r.POST("/todo", CreateHandler)

	r.GET("/todo", GetAllHandler)

	r.GET("/todo/:id", GetHandler)

	r.PUT("/todo/:id", UpdateHandler)

	r.DELETE("/todo/:id", DeleteHandler)

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main(){
	lambda.Start(Handler)
}