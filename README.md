## Overview

### Run using docker

`docker-compose up`

### test using ab or hey
`ab -n 2 http://localhost:9001/`
`hey -n 50 http://localhost:9001/`

### with concurrency
`ab -c 10 -n 50 http://localhost:9001/`
`hey -c 10 -n 50 http://localhost:9001/`




https://github.com/go-redis/redis/blob/0064936c5b77c874f8e1f089c955ec8bef3818c4/sentinel.go#L353