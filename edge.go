package main

import (
  "net/http"
  "github.com/labstack/echo"
  "log"
  "github.com/go-redis/redis"
  "encoding/json"
  "time"
  "fmt"
  "math/rand"
)

const (
  PreFilterChannelName string = "pre-filter"
  RespondChannelName string = "respond"
  Letters string = "0123456789ABCDEF"
)

var (
  redisdb   *redis.Client
  preFilter *redis.PubSub
  respond   *redis.PubSub
)

type response struct {
  PubSub  *redis.PubSub
  RequestID   string
  Status      int
  Value       string
  Context     echo.Context
}

func init() {
  rand.Seed(time.Now().UnixNano())
  redisdb = redis.NewClient(&redis.Options{
    Addr:         ":6379",
    DialTimeout:  10 * time.Second,
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
    PoolSize:     10,
    PoolTimeout:  30 * time.Second,
  })

  respond = redisdb.Subscribe(RespondChannelName)
  // Wait for confirmation that subscription is created before publishing anything.
  _, err = respond.Receive()
  if err != nil {
    panic(err)
  }

}

func randomBytes(n int) ([]byte, error) {
  bytes := make([]byte, n)
  _, err := rand.Read(bytes)
  if err != nil {
    return bytes, err
  }

  return bytes, nil
}

func randomString(bytes []byte) string {
  for i, b := range bytes {
    bytes[i] = Letters[b%16]
  }

  return string(bytes)
}

func listen(r *response) error {
  var (
    err error
    val string
    msgi interface{}
  )

  for {
    log.Printf("Beginning to loop")
    msgi, err = r.PubSub.ReceiveTimeout(50000000 * time.Nanosecond)
    if err != nil {
      return err
    }
    msg := msgi.(*redis.Message)
    log.Printf("looking for: %v, Have: %v", r.RequestID, msg.Payload)
    if r.RequestID == msg.Payload {
      log.Printf("Found matching request: %v, %v", msg.Channel, msg.Payload)
      // Read request object.
      val, err = redisdb.Get(r.RequestID).Result()
      log.Printf("val: %v, err: %v", val, err)
      if err != nil {
        r.Status = http.StatusInternalServerError
        r.Value = fmt.Sprintf("error: %v", err)
        return err
      } else {
        r.Status = http.StatusOK
        r.Value = val
      }
      break
    } else {
      log.Printf("Request Not FOUND!!!: looking for: %v, Have: %v", r.RequestID, msg.Payload)
      continue
    }
  }
  return r.Context.String(r.Status, r.Value)
}

func EdgeController(c echo.Context) error {
  var (
    bytes     []byte
    requestID string
    err       error
    cJson     []byte
  )

  log.Printf("context: %v", c)
  cJson, err = json.Marshal(c)
  if err != nil {
    return c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
  }

  bytes, err = randomBytes(32)
  if err != nil {
    return c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
  }
  requestID = randomString(bytes)
  // Store the request object on key requestID.
  // no ttl
  // err = redisdb.Set(requestID, string(cJson), 0).Err()
  // ttl 300000000 * time.Nanosecond == 300 milliseconds
  err = redisdb.Set(requestID, string(cJson), 300000000 * time.Nanosecond).Err()
  if err != nil {
    return c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
  }

  // Publish a message.
  log.Printf("Publish Error: %v, %v", PreFilterChannelName, requestID)
  err = redisdb.Publish(PreFilterChannelName, requestID).Err()
  if err != nil {
    return c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err))
  }

  r := &response{PubSub: respond, RequestID: requestID, Context: c}
  return listen(r)
}
