## Get aws billing on slack by aws sam

![](https://raw.githubusercontent.com/maeda1150/get_aws_billing_on_slack_by_aws_sam/master/images/result.png)

## Require

* Docker & docker-compose
* [AWS SAM](https://github.com/awslabs/serverless-application-model)
* awscli

## Create Slack Slash command & Setup AWS Systems Manager Parameter store

### Create Slash command on Slack

* Note `Token`.

![](https://raw.githubusercontent.com/maeda1150/get_aws_billing_on_slack_by_aws_sam/master/images/slash_command.png)

* Method(HTTP method) is `GET`.

![](https://raw.githubusercontent.com/maeda1150/get_aws_billing_on_slack_by_aws_sam/master/images/http-method.png)

### Setup AWS Systems Manager Parameter store

Set parameter name, `SLACK_TOKEN_FOR_AWS_BILLING`.
Then fill in value Token by getting slash command.

```
$ aws ssm put-parameter --overwrite --name SLACK_TOKEN_FOR_AWS_BILLING --type String --value xxxxxxx
```

![](https://github.com/maeda1150/get_aws_billing_on_slack_by_aws_sam/blob/master/images/ssm.png)

## Development

* provisioning

```
$ docker-compose build
```

* dep ensure

```
$ docker-compose run --rm lambda dep ensure
```

* development with vim-go

```
$ docker-compose run --rm lambda
```

* build binary

```
$ docker-compose run --rm -e GOOS=linux -e GOARCH=amd64 lambda go build -o aws_billing .
```

or in container

```
$ GOOS=linux GOARCH=amd64 go build -o aws_billing .
```

* validate template.yml

```
$ sam validate
```

* start api
  * Need setup `SLACK_TOKEN_FOR_AWS_BILLING` on `AWS Systems Manager Parameter store`.
  * Need environment variables of AWS_ACCESS_KEY_ID & AWS_SECRET_ACCESS_KEY (or AWS_PROFILE) which has [CloudWatchFullAccess, AmazonSSMFullAccess] policies.

```
$ sam local start-api
```

* post localhost

  * get

  ```
  $ curl 'http://127.0.0.1:3000/aws_billing?token={{slash_command_token}}'
  ```

## Deployment

* make s3 bucket
  * NOTE: The S3 bucket should unique in whole world. So this is example.

```
$ aws s3 mb s3://sum-aws-billing
```

* upload s3 bucket & create package.yml

```
$ sam package --template-file template.yaml --s3-bucket {{s3-bucket-name}} --output-template-file package.yaml
```

* deploy cloudformation & lambda

```
$ aws cloudformation deploy --template-file package.yaml --stack-name aws-billing --capabilities CAPABILITY_IAM
```

### If encountered `Unable to upload artifact None referenced by CodeUri parameter of AwsBilling resource.`

If you encountered this error, here is workaround.

```
Unable to upload artifact None referenced by CodeUri parameter of AwsBilling resource.
[Errno 2] No such file or directory: '/example/get_aws_billing_on_slack_by_aws_sam/vendor/github.com/json-iterator/go/skip_tests/array/skip_test.go'
```

```
$ rm /example/get_aws_billing_on_slack_by_aws_sam/vendor/github.com/json-iterator/go/skip_tests/array/skip_test.go'
$ vim /example/get_aws_billing_on_slack_by_aws_sam/vendor/github.com/json-iterator/go/skip_tests/array/skip_test.go'
```

Then paste this.
https://github.com/json-iterator/go/blob/master/skip_tests/array_test.go

**This error may occur 3 times. Just do same workaround.**

* Set url on Slack Slash command.

![](https://github.com/maeda1150/get_aws_billing_on_slack_by_aws_sam/blob/master/images/url.png)

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
