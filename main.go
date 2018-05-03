package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
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
	billing, err := GetBilling()
	if err != nil {
		c.String(400, "Get Billing error")
		return
	}
	c.String(200, strconv.FormatFloat(billing, 'f', 4, 64)+" $")
}

const (
	region         = "us-east-1"
	namespace      = "AWS/Billing"
	metricName     = "EstimatedCharges"
	dimensionName  = "Currency"
	dimensionValue = "USD"
	period6hours   = 21600
	periodDay      = 86400
)

func GetBilling() (float64, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		log.Println(err)
		return 0, err
	}

	svc := cloudwatch.New(sess)

	params := &cloudwatch.GetMetricStatisticsInput{
		Namespace:  aws.String(namespace),
		MetricName: aws.String(metricName),
		Period:     aws.Int64(period6hours),
		StartTime:  aws.Time(time.Now().Add(time.Duration(periodDay) * time.Second * -1)),
		EndTime:    aws.Time(time.Now()),
		Statistics: []*string{
			aws.String(cloudwatch.StatisticMaximum),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String(dimensionName),
				Value: aws.String(dimensionValue),
			},
		},
		Unit: aws.String(cloudwatch.StandardUnitNone),
	}

	resp, err := svc.GetMetricStatistics(params)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	log.Println(*resp)
	if len(resp.Datapoints) < 1 {
		return 0, fmt.Errorf("Datapoint is 0. Sould extends get Datapoint start time.")
	}

	bills := []float64{}
	for _, bill := range resp.Datapoints {
		bills = append(bills, float64(*bill.Maximum))
	}

	return max(bills), nil
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

func max(a []float64) float64 {
	max := a[0]
	for _, i := range a {
		if i > max {
			max = i
		}
	}
	return max
}
