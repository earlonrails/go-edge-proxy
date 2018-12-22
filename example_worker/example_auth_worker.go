package example_worker

import (
  "log"
  "github.com/go-redis/redis"
  "github.com/valyala/fasthttp"
  "encoding/json"
)

const PreFilterChannelName string = "pre-filter"
const RespondChannelName string = "respond"

type User struct {
    SID         string
    Authorized  bool
}

type Request struct {
    User  User
    Raw   string
}

// MarshalBinary -
func (r *Request) MarshalBinary() ([]byte, error) {
  return json.Marshal(e)
}

// UnmarshalBinary -
func (r *Request) UnmarshalBinary(data []byte) error {
  if err := json.Unmarshal(data, &r); err != nil {
    return err
  }

  return nil
}

var (
  redisdb   *redis.Client
  preFilter *redis.PubSub
  respond   *redis.PubSub
  err       error
  request   Request
  httpCliet *fasthttp.HostClient
)

func createRedisClient() {
  redisdb = redis.NewClient(&redis.Options{
    Addr:         ":6379",
    DialTimeout:  10 * time.Second,
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    PoolSize:     10,
    PoolTimeout:  30 * time.Second,
  })
  preFilter = redisdb.Subscribe(PreFilterChannelName)
  // Wait for confirmation that subscription is created before publishing anything.
  _, err = preFilter.Receive()
  if err != nil {
    panic(err)
  }

  respond = redisdb.Subscribe(RespondChannelName)
  // Wait for confirmation that subscription is created before publishing anything.
  _, err = respond.Receive()
  if err != nil {
    panic(err)
  }
  // Perpare a client, which fetches webpages via HTTP proxy listening
  // on the localhost:8080.
  httpCliet = &fasthttp.HostClient{
    Addr: "localhost:8080",
  }
}

func authorize(rawRequest string) (Request, error) {
    var newRequest Request

    if err := newRequest.UnmarshalBinary([]byte(rawRequest)); err != nil {
      fmt.Printf("Unable to unmarshal data into the Request struct due to: %s \n", err)
      return
    }

    statusCode, body, err := c.Get(nil, "http://google.com/foo/bar")

}

func pollMessages() {
  ch := preFilter.Channel()
  err = redisdb.Set(requestID, string(cJson), 300000000 * time.Nanosecond).Err()
  if err != nil {
    log.Printf("error: %v", err)
  }

  // Consume messages.
  for msg := range ch {
    log.Println("Processing Message: %v, %v", msg.Channel, msg.Payload)
    // Read request object.
    val, err = redisdb.Get(msg.Payload).Result()
    if err != nil {
      log.Printf("error: %v", err)
    }

    request, err = authorize(val)

    err = redisdb.Publish(RespondChannelName, requestID).Err()
    if err != nil {
      log.Printf("error: %v", err)
    }
  }
}

func main() {
  createRedisClient()
  pollMessages(c)
}

