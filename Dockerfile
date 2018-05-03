FROM masaki1111/golang-vim-dev

RUN go get -u github.com/golang/dep/cmd/dep && \
    go get -u github.com/golang/lint/golint

RUN go get -u github.com/aws/aws-lambda-go/events && \
    go get -u github.com/aws/aws-lambda-go/lambda && \
    go get -u github.com/awslabs/aws-lambda-go-api-proxy/... && \
	  go get -u github.com/gin-gonic/gin && \
	  go get -u github.com/aws/aws-sdk-go

WORKDIR /go/src/get_aws_billing_on_slack_by_aws_sam/
