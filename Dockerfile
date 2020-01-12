# Start by building the application.
FROM golang:1.13-buster as build
WORKDIR /opt
ADD . .
RUN make build

# Now copy it into our base image.
#FROM gcr.io/distroless/base-debian10
FROM alpine
COPY --from=build /opt/app /
COPY --from=build /opt/templ /templ
RUN apk add tzdata git && \
    cp /usr/share/zoneinfo/Asia/Chongqing /etc/localtime && \
    echo "Asia/Chongqing" > /etc/timezone && \
    apk del tzdata
ENTRYPOINT ["/app", "/conf.yaml"]
