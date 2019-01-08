package main

import (
    "encoding/json"
    "fmt"
    "github.com/go-redis/redis"
    "github.com/labstack/echo"
    "log"
    "net/http"
    "time"
)

const (
    PreFilterChannelName string = "pre-filter"
    RespondChannelName   string = "respond"
    MaxMessages = 100
)

var (
    redisdb   *redis.Client
    respond   *redis.PubSub
)

type response struct {
    Status    int
    Value     string
}

func init() {
    redisdb = redis.NewClient(&redis.Options{
        Addr:         "redis:6379",
        DialTimeout:  10 * time.Second,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        PoolSize:     10,
        PoolTimeout:  30 * time.Second,
    })
    respond = redisdb.PSubscribe(fmt.Sprintf("%v:*", RespondChannelName))
    // Wait for confirmation that subscription is created before publishing anything.
    _, err := respond.Receive()
    if err != nil {
        panic(err)
    }
}

func publish(id string, pubbed chan bool) {
    log.Printf("Publishing: %v", fmt.Sprintf("%v:%v", PreFilterChannelName, id))
    err := redisdb.Publish(fmt.Sprintf("%v:%v", PreFilterChannelName, id), id).Err()
    if err != nil {
      log.Printf("Got Error Publish: %v", err.Error())
      pubbed <- false
    }
    pubbed <- true
}

func subscribe(id string, resp chan response) {
    log.Println("My id looks good here:::", id)
    for i := 0; i < MaxMessages; i++ {
        msgi, err := respond.ReceiveTimeout(30000000 * time.Nanosecond)
        if err != nil {
            log.Printf("err: %v", err)
            // resp <- response{http.StatusInternalServerError, "Failed read from redis"}
        }
        switch msg := msgi.(type) {
        case *redis.Message:
            if id == msg.Payload {
                val, err := redisdb.Get(id).Result()
                log.Printf("val: %v, err: %v", val, err)
                if err != nil {
                    resp <- response{http.StatusInternalServerError, "Failed read from redis"}
                    close(resp)
                    break
                }
                resp <- response{http.StatusOK, val}
                close(resp)
                break
            }
            log.Println("My id does not look correct anymore:::", id)
            log.Printf("Not Found! id: %v, payload: %v", id, msg.Payload)
            // resp <- response{http.StatusInternalServerError, "Not Found!"}
        default:
            log.Printf("Not a Message: %v", msg)
        }
    }
}

func queue_request(id string) response {
    req_channel := make(chan bool)
    resp_channel := make(chan response)
    go subscribe(id, resp_channel)
    go publish(id, req_channel)

    req, resp := <-req_channel, <-resp_channel
    if req == false {
        return response{http.StatusInternalServerError, "Failed to publish"}
    }

    return resp
}

func EdgeController(c echo.Context) error {
    var (
        cJson   []byte
        err     error
    )

    log.Printf("context: %v", c)
    cJson, err = json.Marshal(c)
    if err != nil {
        log.Printf("Got Error Marshal: %v, %v", cJson, err.Error())
        return c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err.Error()))
    }

    requestID := c.Response().Header().Get(echo.HeaderXRequestID)
    // Store the request object on key requestID.
    // no ttl
    // err = redisdb.Set(requestID, string(cJson), 0).Err()
    // ttl 300000000 * time.Nanosecond == 300 milliseconds
    log.Printf("Setting key on redis: %v, %v, ", requestID, string(cJson))
    err = redisdb.Set(requestID, string(cJson), 300000000 * time.Nanosecond).Err()
    if err != nil {
        log.Printf("Got Error Set: %v", err.Error())
        return c.String(http.StatusInternalServerError, fmt.Sprintf("error: %v", err.Error()))
    }

    resp := queue_request(requestID)
    return c.String(resp.Status, resp.Value)
}
