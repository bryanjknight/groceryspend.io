# Start by building the application.
FROM golang:1.16-buster as build

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

ENV GIN_MODE=release
RUN go build -o /go/bin/app

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /

# we need to include the source to do database migrations
COPY --from=build /go/src/app /
CMD ["/app"]