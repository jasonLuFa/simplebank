# [Docker compose](https://docs.docker.com/compose/)

- Start multiple services at once and control their start-up orders
- `docker compose up` : 啟動或重啟 docker-compose.yaml 下定義的所有 services
  - postgres 和 simplebank API service 都運行在 docker network ( 內網 )
- `docker compose down` : 移除所有已存在的 containers 和 network
