# Start by building the application.
FROM golang:1.13-buster as build
WORKDIR /opt
ADD . .
RUN make build

# Now copy it into our base image.
FROM gcr.io/distroless/base-debian10
COPY --from=build /opt/app /
COPY --from=build /opt/templ /templ
COPY --from=busybox /bin/busybox /busybox/busybox
RUN ["/busybox/busybox", "--install", "/bin"]
ENV TZ=Asia/Chongqing
ENTRYPOINT ["/app", "/conf.yaml"]
