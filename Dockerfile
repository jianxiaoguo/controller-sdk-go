FROM drycc/go-dev:v0.22.0
# This Dockerfile is used to bundle the source and all dependencies into an image for testing.

ADD https://codecov.io/bash /usr/local/bin/codecov
RUN chmod +x /usr/local/bin/codecov

COPY glide.yaml /go/src/github.com/drycc/controller-sdk-go/
COPY glide.lock /go/src/github.com/drycc/controller-sdk-go/

WORKDIR /go/src/github.com/drycc/controller-sdk-go

RUN glide install --strip-vendor

COPY ./_scripts /usr/local/bin

COPY . /go/src/github.com/drycc/controller-sdk-go
