FROM masaki1111/golang-vim-dev

USER root
RUN go get -u github.com/golang/dep/cmd/dep && \
    go get -u github.com/golang/lint/golint

WORKDIR /go/src/get_aws_billing_on_slack_by_aws_sam/
