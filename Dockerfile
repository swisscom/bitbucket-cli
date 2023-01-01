FROM golang:1.18 as build

WORKDIR /go/src/bitbucket-cli
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/src/bitbucket-cli ./cmd/bitbucket-cli

# Now copy it into our base image.
FROM gcr.io/distroless/static-debian11
COPY --from=build /go/src/bitbucket-cli /
ENTRYPOINT ["/bitbucket-cli"]
