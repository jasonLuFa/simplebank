# <h1 align="center">:bank: Simplebank </h1>

## 💪 Purpose

- This project is to familiar with some development tools and tech using golang
  - [gin](https://github.com/gin-gonic/gin) web framework
  - [viper](https://github.com/spf13/viper) to read config
  - [goMock](https://github.com/golang/mock) to test API with CRUD
  - [validator](https://github.com/go-playground/validator) to custom the validation to request data
  - [strechter/testify](https://github.com/stretchr/testify) to write a unit test

## 📋 Content

- [DataBase](https://github.com/jasonLuFa/simplebank#-database)
  - [postgres](https://github.com/jasonLuFa/simplebank#%EF%B8%8F-create-a-postgres-instance)
  - [Concept of DB transaction](https://github.com/jasonLuFa/simplebank#%EF%B8%8F-concept-of-db-transaction)
  - [transaction isolation level](https://github.com/jasonLuFa/simplebank#%EF%B8%8F-transaction-isolation-level)
- [Makefile](https://github.com/jasonLuFa/simplebank#%EF%B8%8F-makefile)
- [golang migrate](https://github.com/jasonLuFa/simplebank#-golang-migrate)
- [Sqlc](https://github.com/jasonLuFa/simplebank#-sqlc)
- [Github Action](https://github.com/jasonLuFa/simplebank#-github-action)
- [PASETO vs JWT](https://github.com/jasonLuFa/simplebank#-paseto-vs-jwt)
- [Dockerfile](https://github.com/jasonLuFa/simplebank#dockerfile)
- [Docker compose](https://github.com/jasonLuFa/simplebank#docker-compose)

## 💾 DataBase

### ✏️ Create a Postgres instance

![ERD](https://github.com/jasonLuFa/simplebank/blob/master/simlebank_dbdiagram.png)

- run postgres in docker container
  1. pull postgres image : `docker pull <image_name>`
     - ex : `docker pull postgres:12-alpine`
  2. start a postgres instance : `docker run --name <container_name> -e POSTGRES_USER=<user_name> -e POSTGRES_PASSWORD=<user_password> -p <host_posts>:<container_ports_in_docker_network> -d <image_name>`
     - ex : `docker run --name postgres12 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -p 5433:5432 -d postgres:12-alpine`
     - `-e` : 設定 container 的環境變數
     - `-d` : 將 contanier 在背景運行( datach mode )
     - `-p <host_posts>:<container_ports_in_docker_network>` : mapping port 通常 docker container 是運行在不同的虛擬網路( virtual network )，而非我們的本機網路( local network )，所以我們須建立兩者間的連線橋梁( bridge )
       - 5432 is the port inside the docker network
       - 5433 is the port you use on your host machine to connect to the port inside the docker network
  3. `docker ps` : list the current running container
     - `-a` : list all container stop or running
  4. exec psql in container : `docker exec -it <container_name> <command>`
     - ex : `docker exec -it postgres12 psql -U admin`
       - psql : command in docker to access the Postgres console
       - -U {username} : connect to the {user name} user
       - NOTE : The PostgreSQL image sets up trust authentication locally so you may notice a password is not required when connecting from localhost (inside the same container). However, a password will be required if connecting from a different host/container
  5. log container info : `docker logs <container_name/container_ID>`
     - ex : `docker logs postgres12`

### ✏️ Concept of DB Transaction

- 為什麼我們需要 transaction ?

  1.  即使系統壞掉了，也可提供可信賴( reliable )和一致性( consistency )的處理
  2.  在併發( concurrently )狀況，不同的程式接觸到 DB 也能彼此獨立

- 目的 : 滿足 ACID
  1.  **A**tomicity : 所有 transaction 操作成功則執行，否則回滾( rollback )所有操作，DB 資訊維持一樣
  2.  **C**onsistency : 所有交易必須符合設定的限制
      - [Constraint](https://www.postgresql.org/docs/current/ddl-constraints.html) ( 條件約束 ) : 欄位上的限制 ( ex : UNIQUE , amount > 0 ... )
      - [Cascade](https://www.postgresql.org/docs/8.2/ddl-constraints.html) ( 級聯 ) : 當 PK 裡的值被刪除或更新，別張表中相對應其值的 FK 是否也 被刪除或更新
      - [Trigger](https://www.postgresql.org/docs/current/sql-createtrigger.html) ( 觸發程序 ) : 在 INSERT、UPDATE、DELETE 等事件發生時，會觸發執行的程式
  3.  **I**solation : 併發中的所有 transaction 彼此獨立( 根據 Isolation level )，互不影響
  4.  Durability : 所有被寫入的 transaction 必須持久化的儲存( persistent storage )

#### Deadlock

---

- 我們為什麼要 lock 資料表 ?

  - 因為多個 transaction _Update_ DB 資料前，如果沒 lock 會先 _Read_ 到相同的資料 ( 例 : balance =100 )，此時 transaction1 _Update_ 後( 例 : balance = 90 )，transaction2 _Read_ 到的資料會是舊的 ( 例 : balance =100 )
  - 所以我們想 lock 住資料表，一次只給一個 transaction 做 _Update_

- deadlock 的發生情況 : 當兩個 transaction 都在互等對方時

#### Lock 種類

---

- `shared lock` ( read lock ) : 查詢更新的資料，例如 : SELECT
  - 當此頁面有 shared lock，其他 concurrency 只可取得 shared lock 或 update lock
- `update lock` : 更新前把資料改成 Update Lock
  - 當此頁面有 update lock，其他 concurrency 只能取得 shared lock
- `exclusive lock` ( write lock ) : 確定要更新當下改成 exclusive lock.，例如 : INSERT、UPDATE、DELETE

  - 當此頁面有 exclusive lock，其他 concurrency 皆無法取的

- lock 兼容矩陣

  |                | shared lock | exclusive lock |
  | -------------- | ----------- | -------------- |
  | shared lock    | ✅          | ❌             |
  | exclusive lock | ❌          | ❌             |

  - _Update_ 一筆資料，會經過的 lock
    - shared lock -> update lock -> exclusive lock

### ✏️ Transaction Isolation Level

- 顯示當前 Isolation Levels
  - MySQL : `select @@transaction_isolation;`
  - Postgres : `show transaction isolation level;`
- 修改當前 Isolation Levels
  - MySQL : 修改當前 session 的 Transactions Levels ( 如不加 session 則為 global ) `set session transaction isolation level <isolation level name>`
  - Postgres : 只能在每個 transaction 去修改其 transaction level `set transaction isolation level <isolation level name>`

#### Read Phenomena

---

##### dirty read

- 讀取到不正確的資訊 ( 未 commit 的資料 )( 可能是未 commit 的資料，但後來被 rollback 了，造成資訊錯誤 )

##### non-Repeatable read

- 在相同 transaction 中，相同 query 得到不同欄位值的結果( 因為可能其他 committed transaction 修改到該欄位 )

##### phantom read

- 在相同 transaction 中，相同 query 得到不同筆數( set of rows )的結果( 因為可能其他 committed transaction 修改到資料表的資料筆數 )

##### serialization anomaly

- 多個 transactions 同時進行並 commit 時，結果跟將其 transaction 隨機順序依序 commit 的結果不同 ( the result of successfully committing a group of transactions is inconsistent with all possible orderings of running those transactions one at a time. )

#### 4 Standard Isolation levels ( low -> High )

---

##### Read uncommitted

- MySQL : 會讀取到未 commit 的資料
- Postgres : 此 isolaion level 跟 Read committed level 一樣 ( 所以也可以所 postgres 只有下面三種 isolation levels )

##### read committed

- 只會讀取已 commit 的資料 ( MySQL 和 Postgres 都一樣 )

##### repeatable read

- 在此 transaction 中，讀到的資訊都不變，不論其他 transaction 有沒有新的 commit
  - MySQL :
    - 例 : transaction_A 的 table_A 有一欄位值為 80，只要 transaction_A 未 commit 前，讀取到的都是 80 ; 當 transaction_B 修改 table_A 該欄位值為 70 並 commit，此時 transaction_A 讀取到仍是 80，但如果 transaction_A 想將該欄位值 -10，此時 commit 後會是正確的 60，雖然結果是正確的，但是過程中讀取的資訊是錯誤的
  - Postgres :
    - 例 : 讀取資訊都不變，但如果有其他 transaction 已經修改並 commit Table_A 欄位 A，則其他 transaction 無法再進行修改 Table_A 欄位 A
    - 例 : 如同時進行兩個 transaction，但 transaction A 先新增一筆資料，為當前 Table_A 欄位 A 為所有加總( ex : 30+40+50 )並 commit 後欄位 A 會有四筆資料 30, 40, 50, 120 ; 此時 transaction B 在 transaction A commit 完做相同事( 會得到四筆資料 30, 40, 50, 120 )， 但是 commit 後會得到 30, 40, 50, 120, 120，此為錯誤答案( 正確應該為 30, 40, 50, 120, 240 )，此及 serialization anomaly

##### serializable :

- MySQL : 會將所有 SELECT query 轉換成 SELECT FOR SHARE，並當有多個 transactions 運行時，只允許讀取資料不能更新或刪除，如果其中一個 transaction 嘗試更新或刪除，則會被阻擋( block )，等到其他 transaction commit 才執行，或是等到 timeout 直接中斷 ; 如兩個 transaction 都執行更新或刪除，則會產生 deadlock 並直接中斷
- Postgres : 可解決 Repeatable Read 範例 2 的問題

- 結論 :

| MSQL                  | read uncommitted | read committed | repeatable read | serializable |
| --------------------- | ---------------- | -------------- | --------------- | ------------ |
| dirty read            | ✅               | ❌             | ❌              | ❌           |
| non-repeatable read   | ✅               | ✅             | ❌              | ❌           |
| phantom read          | ✅               | ✅             | ❌              | ❌           |
| serialization anomaly | ✅               | ✅             | ✅              | ❌           |

- 有四種 isolation levels
- 使用 locking Mechanism
- 預設 isolation level 為 : repeated read

---

| Postgres              | read uncommitted | read committed | repeatable read | serializable |
| --------------------- | ---------------- | -------------- | --------------- | ------------ |
| dirty read            | ❌               | ❌             | ❌              | ❌           |
| non-repeatable read   | ✅               | ✅             | ❌              | ❌           |
| phantom read          | ✅               | ✅             | ❌              | ❌           |
| serialization anomaly | ✅               | ✅             | ✅              | ❌           |

- 有三種 isolation levels
- 使用 dependencies detection
- 預設 isolation level 為 : read committed
- 當 isolation level 越高，就會有越多的 error, timeout, deadlock 要處理

# ⌨️ Makefile

- 使用 Makefile 來管理我們程式的構建，減少了大量輸入、拼寫錯誤，簡化構建項目的難度。真實線上環境配合 CI/CD 更佳
- 將一長串的的 command 縮短成自訂的 command
- ex : 建一個 Makefile 檔案

  ```
  postgres:
    docker run --name postgres12 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -p 5433:5432 -d postgres:12-alpine
  test:
    go test -v -cover ./...
  ```

  - 之後 command 只要下 `make postgres`，則就會執行 `docker run --name postgres12 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -p 5433:5432 -d postgres:12-alpine` 這一長串的指令了

## 💻 [golang migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### [常用指令](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#usage)

- `create [-ext E] [-dir D] [-seq] [-digits N] [-format] NAME`
- `migrate up {N}`
- `migrate down {N}`

### 範例

1. `mkdir -p db/migration`
2. `migrate create -ext sql -dir db/migration -seq init_schema`
   - 會在 db/migration 下自動生成兩個檔案 000001_init_schema.down.sql 和 000001_init_schema.sql
3. 使用 `migrate -path db/migration -database "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable" -verbose up` 可執行 000001_init_schema.up.sql 此檔案裡的所有 sql 指令，`migrate -path db/migration -database "postgresql://admin:admin@localhost:5433/simple_bank?sslmode=disable" -verbose down` 則執行 000001_init_schema.down.sql 此檔案裡的所有 sql 指令 ( 這兩個指令也可寫在 Makefile 裡 )

## 💻 [Sqlc](https://github.com/kyleconroy/sqlc)

- Writing the CRUD
- Why to use sqlc instead of GORM, SQLX ?
  1. learning how to write sql (GORM must learn to write queries using gorm's funnction )
  1. very fast & easy to use ( GORM run slowly on high load )
  1. Automatic code generation
  1. Catch SQL query errors before generating codes ( SQLX won't occur error util runtime )

## 💻 [Github Action](https://docs.github.com/en/actions)

- 可藉由 GitHub Action 去自動測試我們所有的 testing，並設定觸發時間( 例如 : pull request, merge to master ... )
- 以下名詞定義 ( 顆粒度 : Workflow > Job = Runner > step > Action ):
  - Workflow :
    1.  是一個自動的程序
    2.  由多個 job 所組成
    3.  可藉由 Event, scheduled, manually 來觸發
    4.  使用 .yaml file 來撰寫
  - Runner :
    1. 用來運行 job 的 server
    2. 一次只會運行一個 job
    3. 會將結果回傳 github
  - Job :
    1. 一系列的步驟運行在 runner 中
    2. 種類
       - normal jobs 平行運行
       - dependent jobs 依序運行
  - Step :
    1. 在 job 中依序運行
    2. 包含多個 Action
  - Action :
    1. 是獨立的指令
    2. 在 Step 中依序運行
    3. 可以被重複使用 ( 所以可以使用別人撰寫好了 Github Action )

## 💻 [PASETO](github.com/o1egl/paseto) vs [JWT](https://jwt.io/)

### ✏️ Token-based Authentication

- access token normally has a lifetime duration before it gets expired
-

### ✏️ Symmetric key Encryption

- 只有一個 secret key 用於加密( encryption )和解密( dencryption )
- 加密過程非常快
- 適合處理大型資料

### ✏️ Asymetric key Encryption

- 有 public key 和 private key
  - public key 用於加密 ( 每個人都可以獲得 )
  - private key 用於解密 ( 必須保護好 )
- 加密過程很慢
- 適合處理小型資料

### ✏️ JWT

- 所有資料存在 JWT 只有 base64 編碼，並非加密，所以避免存取重要資訊
- token 可分為 Header, payload, verify singature
  - Header : signing algorithm
  - payload : 關於登入者的資訊、包括 token 過期時間 ... 你想存的資訊
  - verify singature : 數位簽章( digital signature )

### ✏️ PASETO( platform agnostic security tokens )

- 跟 JWT 類似，但是比 JWT 更安全及容易實作

## [Dockerfile](https://docs.docker.com/engine/reference/builder/)

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

### ✏️ use multistage to reduce the size of docker image size

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

## [Docker compose](https://docs.docker.com/compose/)

- Start multiple services at once and control their start-up orders
- `docker compose up` : 啟動或重啟 docker-compose.yaml 下定義的所有 services
  - postgres 和 simplebank API service 都運行在 docker network ( 內網 )
- `docker compose down` : 移除所有已存在的 containers 和 network
