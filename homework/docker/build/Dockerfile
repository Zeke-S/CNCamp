FROM golang:1.17.6-alpine AS build
RUN apk add --no-cache git
# RUN go get github.com/Zeke-S/CNCamp

WORKDIR /
RUN git clone https://github.com/Zeke-S/CNCamp.git
WORKDIR /CNCamp/homework/httpserver
RUN CGO_ENABLED=0 GOOS=linux go build 


FROM scratch
COPY --from=build /CNCamp/homework/httpserver/httpserver /bin/httpserver
ENTRYPOINT ["/bin/httpserver"]