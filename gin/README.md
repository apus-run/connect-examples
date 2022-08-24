## 构建 server 和 client
```shell
make build
```

## 启动 server
```shell
$ ./bin/gin-server
```

## 启动 client
```shell
$ ./bin/client
```

## 访问服务

run client with curl command

```sh
curl --header "Content-Type: application/json" \
    --data '{"name": "foobar"}' \
        http://localhost:8800/user.v1.UserService/Say
```

health check

```sh
curl --header "Content-Type: application/json" \
    --data '{"service": "user.v1.UserService"}' \
    http://localhost:8800/grpc.health.v1.Health/Check
```

run client with [grpcurl](https://github.com/fullstorydev/grpcurl) command

```sh
grpcurl \
  -plaintext \
  -d '{"name": "foobar"}' \
  localhost:8800 \
  user.v1.UserService/Say
```

health check

```sh
grpcurl \
  -plaintext \
  -d '{"service": "user.v1.UserService"}' \
  localhost:8800 \
  grpc.health.v1.Health/Check
```


## 登录服务
```sh
curl --header "Content-Type: application/json" \
    --data '{"username": "moocss", "password": "123456"}' \
        http://localhost:8800/user.v1.UserService/Login
```

## 机器人服务
```sh
grpcurl \
-plaintext \
-d '{"name": "foobar"}' \
localhost:8800 \
user.v1.UserService/Introduce
```
