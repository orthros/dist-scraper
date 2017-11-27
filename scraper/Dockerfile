# Docker image for the display-incrementer

# Use the golang image to compile our application
FROM golang:alpine as builder

RUN apk add --no-cache git

WORKDIR /go/src/app
COPY ./ /go/src/app
# go-wrapper download and install do the go get for our dependencies automatically
RUN go-wrapper download
RUN go-wrapper install
# Build our exe and call it scraper
RUN go build -o scraper

# We'll use this as our running image.
# Todo: figure out what is the best candidate for it
# I feel like alpine is the best option here...
FROM alpine 

WORKDIR /app

# Copy the executable "scraper" from our build image to our 
COPY --from=builder /go/src/app/scraper ./

CMD ["./scraper"]