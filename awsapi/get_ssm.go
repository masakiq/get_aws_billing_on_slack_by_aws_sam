package awsapi

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func GetSsm(name string, region string) (string, error) {
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
