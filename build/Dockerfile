FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
WORKDIR /src/
COPY . .

RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags="-w -s" -o /bin/helmproj

# -----------------------------------------------

FROM alpine
COPY --from=builder /bin/helmproj /bin/helmproj
ENTRYPOINT ["/bin/helmproj"]