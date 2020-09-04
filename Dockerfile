# Start by building the application.
FROM golang:1.15-buster as build
WORKDIR /opt
ADD . .
RUN make build

# Now copy it into our base image.
#FROM gcr.io/distroless/base-debian10
# sadly, git and its dependencies is too large,
# in the future we can use https://github.com/src-d/go-git
FROM debian
COPY --from=build /opt/app /
COPY --from=build /opt/templ /templ
RUN apt update && apt install -y git
ENV TZ=Asia/Chongqing
ENTRYPOINT ["/app", "/conf.yaml"]
