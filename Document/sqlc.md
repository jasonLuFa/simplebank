# 💻 Query DB in golang ([sqlc](https://github.com/kyleconroy/sqlc))
## Four way to write the CRUD
### :one: [database/sql](https://pkg.go.dev/database/sql) packge : low-level standard library 
- pros
  - run fast & straighforward
- cons
  - write code with sql syntax, and manaual mapping SQL fields to variables
  - Easy to make mistakes, not caught util runtime
### :two: [GORM](https://gorm.io/docs/)
- pros
  - CRUD functions already implemented very short production code
- cons
  - Must learn to write queries using gorm's function
  - run slowly on high load
### :three: [sqlx](https://github.com/jmoiron/sqlx)
- pros
  - fast & easy to use
- cons
  - Fields mapping via query text & struct tags
  - Failure won't occur until runtime

### 4️⃣ [sqlc](https://github.com/kyleconroy/sqlc)
- pros
  - very fast & easy to use
  - Write SQL queries, then we can use sqlc command automatically generated code which also using database/sql paackage 
  - Catch SQl query errors before generating codes
- cons
  - Version1 Only Support MySQL, Postgres (Version2 supported more)

## Conclusion : Why to use sqlc instead of GORM, SQLX ?
  1. learning how to write sql (GORM must learn to write queries using gorm's funnction )
  2. very fast & easy to use ( GORM run slowly on high load )
  3. Automatic code generation
  4. Catch SQL query errors before generating codes ( SQLX won't occur error util runtime )

## Using [sqlc](https://github.com/kyleconroy/sqlc)
- `sqlc help`
- `sqlc init` : init configuration file( sqlc.yaml )
  - see [sqlc configuration document in v1/v2](https://docs.sqlc.dev/en/stable/reference/config.html)  
  - we use version1 configuration in this project
