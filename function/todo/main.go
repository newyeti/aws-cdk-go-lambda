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
	
	r.POST("/todo", CreateTodo)

	r.GET("/todo/:id", GetTodo)

	r.PUT("/todo/:id", UpdateTodo)

	r.DELETE("/todo/:id", DeleteTodo)

	ginLambda = ginadapter.New(r)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main(){
	lambda.Start(Handler)
}