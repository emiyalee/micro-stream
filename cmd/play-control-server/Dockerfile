FROM emiyalee/stream-system:1.0.1 as builder
WORKDIR $GOPATH/src/github.com/emiyalee/stream-system/play-control-server
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/play_control_server

#FROM golang:1.10.2-stretch
#FROM scratch
FROM golang:1.10.2
COPY --from=builder /bin/play_control_server /bin/play_control_server
EXPOSE 8080
#ENTRYPOINT ["/bin/play_control_server"]
CMD ["/bin/play_control_server", "-h"]