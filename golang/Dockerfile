FROM golang:latest as builder

WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .
RUN go build ./cmd/demo
RUN cp demo /usr/bin/

# FROM alpine
# RUN apk add libc6-compat
# 
# COPY --from=builder /src/demo /usr/bin

# ENTRYPOINT [ "/usr/bin/demo" ]