FROM alpine
RUN apk --update add ca-certificates
ADD radish /
ADD VERSION /
ADD html/dist html/dist
WORKDIR /
EXPOSE 8080
ENTRYPOINT ./radish
