FROM emiyalee/stream-system:1.0.1 as builder
WORKDIR $GOPATH/src/github.com/emiyalee/stream-system/segment-server
COPY ./ ./
#RUN go build -o /bin/segment_server
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/segment_server

#FROM golang:1.10.2-stretch
FROM scratch
COPY --from=builder /bin/segment_server /bin/segment_server
COPY ./ffmpeg /bin/ffmpeg
EXPOSE 8080
#ENTRYPOINT ["/bin/segment_server"] 
CMD [ "/bin/segment_server", "-h"]

