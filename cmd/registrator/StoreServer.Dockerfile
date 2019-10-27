FROM emiyalee/stream-system:1.0.1 as builder
WORKDIR $GOPATH/src/github.com/emiyalee/stream-system/registrator
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/registrator

FROM nginx:1.14.0
COPY --from=builder /bin/registrator /bin
CMD ["registrator", "-h"] 

