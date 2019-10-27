FROM emiyalee/stream-system:1.0.1 as builder
WORKDIR $GOPATH/src/github.com/emiyalee/stream-system/stream-agent-server
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/stream_agent_server

FROM nginx:1.14.0
WORKDIR /usr/share/nginx/html
RUN rm -rf *
COPY --from=builder /bin/stream_agent_server ./
COPY ./conf ./conf
COPY ./static ./static
COPY ./views ./views
COPY nginx.conf /etc/nginx/nginx.conf
CMD ["./stream_agent_server", "-h"] 
