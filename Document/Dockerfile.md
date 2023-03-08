# [Dockerfile](https://docs.docker.com/engine/reference/builder/)

- build image from dockerfile : `docker build -t <image_name>:<tag> .`
- inspect container : `docker container inspect <container name>`
- 運行容器並輸入指定環境變數 : `docker run --name <container name> -p 8080:8080 -e <env var key=env var value> <image name>`
  - `--network <network name>` : 加到指定 network，否則會到預設的 bridge network
    - 使用 bridge network 只能用 ip:port 做連線，ip 可從 `docker network inspect bridge` 去看指定 container 的 ip ( 會是 127.x.x.x )
    - 使用自訂 network 可以用 container name 作為 ip:port
  - `-e <env var key=env var value> <image name>` : 輸入 container 所需的環境變數 ( ex : -e POSTGRES_PASSWORD=secret )
- `docker network ls` : list all docker network
- `docker network inspect <network name>` : inspect specific docker network
- `docker network create <network_name>` : create a custom network ( 原 bridge 的 network 無法用 container name 連線只能用 IP )
- `docker network connect <network name> <container name>` : 將 container 加到指定 network

## ✏️ use multistage to reduce the size of docker image size

- create a image about simplebank server

  - origin method ( methdo 1 )

    ```dockerfile
    FROM golang:1.19-alpine3.16
    # create /app directory in image
    WORKDIR /app
    # first dot means copy everything from current folder
    # second dot means current directory( /app ) inside the image
    COPY . .
    COPY app.env .
    # build out app to single binary executible file
    RUN go build -o main main.go
    # It functions as a type of documentation between the person who builds the image and the person who runs the container, about which ports are intended to be published
    EXPOSE 8080
    # default command to run when the container starts
    CMD ["/app/main"]
    ```

  - after useing multistage ( method 2 )

    ```dockerfile
    # build stage
    FROM golang:1.19-alpine3.16 AS builder
    WORKDIR /app
    COPY . .
    RUN go build -o main main.go

    # Run stage
    FROM alpine:3.16
    WORKDIR /app
    # the dot mean /app in alphine:3.13
    COPY --from=builder /app/main .
    COPY app.env .

    EXPOSE 8080
    CMD ["/app/main"]
    ```

- Different : use multistage to build image will smaller 30 times than the origin method
- Reason : 方法一它包含此專案中所有 golang 所需的 package， 但實際上我們需要的只有 build 出來後的執行檔，所以方法二中只存在 build 出來的執行黨
