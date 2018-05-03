## Require

* Docker
* [AWS SAM](https://github.com/awslabs/serverless-application-model)
* awscli

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

* invoke lambda on local

```
$ echo '{"hoge": "fuga"}' | sam local invoke AwsBilling
```

or

```
$ sam local invoke AwsBilling -e event.json
```

* validate template.yml

```
$ sam validate
```

* start api

```
$ sam local start-api
```

* start api with env

```
$ sam local start-api --env-vars env.json
or
$ TOKEN=token sam local start-api
```

* post localhost

  * get

  ```
  $ curl 'http://127.0.0.1:3000/aws_billing?token=token&text=text'
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
  $ curl 'https://{{id}}.execute-api.ap-northeast-1.amazonaws.com/Prod/aws_billing?token=token&text=text'
  ```

REF: https://github.com/awslabs/aws-lambda-go-api-proxy
