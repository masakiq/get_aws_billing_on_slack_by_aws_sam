## Require

* Docker & docker-compose
* [AWS SAM](https://github.com/awslabs/serverless-application-model)
* awscli

## Create Slack Slash command & AWS Systems Manager Parameter store

### Create Slash command on Slach

Note Token.

![](https://raw.githubusercontent.com/maeda1150/get_aws_billing_on_slack_by_aws_sam/master/images/slash_command.png)

### Setup AWS Systems Manager Parameter store

Set parameter name, `SLACK_TOKEN_FOR_AWS_BILLING`.
Then fill in value by getting slash command token.

![](https://github.com/maeda1150/get_aws_billing_on_slack_by_aws_sam/blob/master/images/ssm.png)

## Development

* provisioning

```
$ docker-compose build
```

* development on golang-vim-dev

```
$ docker run --rm -tiv `pwd`:/go/src/get_aws_billing_on_slack_by_aws_sam get_aws_billing_on_slack_by_aws_sam_lambda

# fish shell
$ docker run --rm -tiv (pwd):/go/src/get_aws_billing_on_slack_by_aws_sam get_aws_billing_on_slack_by_aws_sam_lambda
```

* build binary(in container)

```
$ GOOS=linux GOARCH=amd64 go build -o aws_billing .
```

* validate template.yml

```
$ sam validate
```

* start api

```
$ sam local start-api
```

  * Need setup `SLACK_TOKEN_FOR_AWS_BILLING` on `AWS Systems Manager Parameter store`.
  * Need environment variables of AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY (or AWS_PROFILE) which has [CloudWatchFullAccess, AmazonSSMFullAccess] policies.

* post localhost

  * get

  ```
  $ curl 'http://127.0.0.1:3000/aws_billing?token={{slash_command_token}}'
  ```

## Deployment

* make s3 bucket

```
$ aws s3 mb s3://sum-aws-billing
```

* upload s3 bucket & create package.yml

```
$ sam package --template-file template.yaml --s3-bucket sum-aws-billing --output-template-file package.yaml
```

* modify pachage.yml Environment Variables

  * TODO auto aupdate

* deploy cloudformation & lambda

```
$ aws cloudformation deploy --template-file package.yaml --stack-name aws-billing --capabilities CAPABILITY_IAM
```

## Test deployed api

* get apigateway id

```
$ aws apigateway get-rest-apis --output json --query 'items[?name==`aws-billing`].id'
```

* post api

  * get

  ```
  $ curl 'https://{{id}}.execute-api.ap-northeast-1.amazonaws.com/Prod/aws_billing?token={{slash_command_token}}'
  ```

REF: https://github.com/awslabs/aws-lambda-go-api-proxy
