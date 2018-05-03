package main

import (
	"get_aws_billing_on_slack_by_aws_sam/awsapi"
	"log"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	gin.SetMode(gin.ReleaseMode)
	if ginLambda == nil {
		log.Printf("Gin cold start")
		r := gin.Default()
		r.GET("/aws_billing", getAwsCosts)

		ginLambda = ginadapter.New(r)
	}

	return ginLambda.Proxy(req)
}

func main() {
	lambda.Start(Handler)
}

func getAwsCosts(c *gin.Context) {
	token, err := SsmGet("SLACK_TOKEN_FOR_AWS_BILLING", "ap-northeast-1")
	if err != nil {
		c.String(400, "Ssm error")
		return
	}
	if token != c.Query("token") {
		c.String(400, "Incorrect token")
		return
	}
	billing, err := awsapi.GetBilling()
	if err != nil {
		c.String(400, "Get Billing error")
		return
	}
	c.String(200, strconv.FormatFloat(billing, 'f', 4, 64)+" $")
}

func SsmGet(name string, region string) (string, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Println(err)
		return "", err
	}
	svc := ssm.New(sess, &aws.Config{})
	input := &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(false),
	}
	param, err := svc.GetParameter(input)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return *param.Parameter.Value, nil
}
