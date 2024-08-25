package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// Replace these values with your actual Cognito pool ID, client ID, and region
var userPoolID = os.Getenv("UserpoolID")
var clientID = os.Getenv("ClientsID")

// Initialize the Cognito Identity Provider client
var cognitoClient *cognitoidentityprovider.CognitoIdentityProvider

// Initialize gin lambda
var ginLambda *ginadapter.GinLambda

type AuthRequest struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
        RefreshToken    string `json:"refreshToken"`
}

type AuthResponse struct {
        AccessToken string `json:"accessToken"`
        AccessTokenExpiry int64 `json:"accessTokenExpiry"`
        RefreshToken string `json:"refreshToken"`
        RefreshTokenExpiry int64 `json:"refreshTokenExpiry"`
}

type ErrorResponse struct {
        Message string `json:"message"`
}

func init() {

        log.Printf("Initializing Auth Lambda")

        r := gin.Default()
        r.POST("/auth", LambdaHandler)
        
        ginLambda = ginadapter.New(r)

        sess := session.Must(session.NewSessionWithOptions(session.Options{
                Config: aws.Config{Region: aws.String("us-east-1")},
        }))

        cognitoClient = cognitoidentityprovider.New(sess)

}

func initiateAuth(email, password string) (*cognitoidentityprovider.AuthenticationResultType, error) {
        input := &cognitoidentityprovider.AdminInitiateAuthInput{
                UserPoolId:   aws.String(userPoolID),
                ClientId:     aws.String(clientID),
                AuthFlow:     aws.String("ADMIN_NO_SRP_AUTH"),
                AuthParameters: map[string]*string{
                        "USERNAME":  aws.String(email),
                        "PASSWORD": aws.String(password),
                },
                ClientMetadata:map[string]*string{
                        "username":  aws.String(email),
                        "password": aws.String(password),
                },
        }

        result, err := cognitoClient.AdminInitiateAuth(input)
        if err != nil {
                return nil, err
        }

        return result.AuthenticationResult, nil
}

func refreshAuth(refreshToken string) (*cognitoidentityprovider.AuthenticationResultType, error) {
        input := &cognitoidentityprovider.AdminInitiateAuthInput{
                UserPoolId:   aws.String(userPoolID),
                ClientId:     aws.String(clientID),
                AuthFlow:     aws.String("REFRESH_TOKEN_AUTH"),
                AuthParameters: map[string]*string{
                        "REFRESH_TOKEN": aws.String(refreshToken),
                },
        }

        result, err := cognitoClient.AdminInitiateAuth(input)
        if err != nil {
                return nil, err
        }

        return result.AuthenticationResult, nil
}

func LambdaHandler(c *gin.Context) {
        var authRequest AuthRequest
        var authResponse AuthResponse

        e := c.BindJSON(&authRequest)
        if e != nil {
                c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: e.Error()})
                return
        }

        email := authRequest.Email
        password := authRequest.Password
        refreshToken := authRequest.RefreshToken

        var resp *cognitoidentityprovider.AuthenticationResultType
        var err error

        if len(password) > 0 {
                resp, err = initiateAuth(email, password)
        } else if len(refreshToken) > 0 {
                resp, err = refreshAuth(refreshToken)
        }

        if err != nil {
                c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
                return
        }
        
        authResponse = AuthResponse{
                AccessToken: *resp.IdToken,
                AccessTokenExpiry: *resp.ExpiresIn,
                RefreshToken: *resp.RefreshToken,
                RefreshTokenExpiry: *resp.ExpiresIn,
        }

        c.JSON(200, authResponse)
         
}

func Handler(ctx aws.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main(){
	lambda.Start(Handler)
}