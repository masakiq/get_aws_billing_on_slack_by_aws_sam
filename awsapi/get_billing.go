package awsapi

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

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
		return 0, fmt.Errorf("Datapoint is empty. Sould extends get Datapoint range.")
	}

	bills := []float64{}
	for _, bill := range resp.Datapoints {
		bills = append(bills, float64(*bill.Maximum))
	}

	return max(bills), nil
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
