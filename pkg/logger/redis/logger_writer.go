package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

//Options ...
type Options struct {
	Host string
	Port int32
	Key  string
}

//LoggerWriter ...
type LoggerWriter struct {
	client *redis.Client
	key    string
}

//NewLoggerWriter ...
func NewLoggerWriter(options *Options) (*LoggerWriter, error) {
	address := fmt.Sprintf("%s:%d", options.Host, options.Port)
	client := redis.NewClient(
		&redis.Options{
			Addr:     address,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	loggerWriter := &LoggerWriter{
		client: client,
		key:    options.Key,
	}
	return loggerWriter, nil
}

//TestConnection ...
func (w *LoggerWriter) TestConnection() bool {
	_, err := w.client.Ping().Result()
	return err == nil
}

func (w *LoggerWriter) Write(p []byte) (n int, err error) {
	_, err = w.client.RPush(w.key, p).Result()
	return len(p), err
}
