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
  - see [sqlc configuration document](https://docs.sqlc.dev/en/stable/reference/config.html) in [v1](https://docs.sqlc.dev/en/stable/reference/config.html#version-2) and [v2](https://docs.sqlc.dev/en/stable/reference/config.html#version-1)
  - we use version1 configuration in this project 
  
    ```yaml
    version: "1"
    packages:
      - name: "db"                     # name of the go package that will be generated
        path: "./db/sqlc"              # place to store the generated golang code files
        queries: "./db/query/"         # tell sqlc where to look for the SQL query files
        schema: "./db/migration/"      # place to containing the database schema or migration files
        engine: "postgresql"           # what database engine
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true           # generated Querier interface for the generated package, it might be useful to mock db for testing 
        emit_exact_table_names: false  # if set false, sqlc attempts to singularize plural table names ( ex: accounts table -> Account struct )
        emit_empty_slices: true        # If true, slices returned by :many queries will be empty instead of nil
    ```
  - here is the version2 configurations which are equivalent in terms of their functionality.
  
    ```yaml
    version: "2"
    sql:
      - engine: "postgresql"
        queries: "./db/query/"
        schema: "./db/migration/"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_table_aliases: false
        emit_empty_slices: true
        gen:
          go:
            package: "db"
            dir: "./db/sqlc"
    ```
- `sqlc generate` : automatically generate the code to `./db/sqlc` ,based on `./db/query` folder(base on setting of configurations)
  - `./db/query/account.sql`
    - Note : comment on top of query will instruct sqlc how to generate the Golang function signature for this query
    - ex : name of the function will be CreateAccount, and it should return 1 single Account object
    
      ```sql
      -- name: CreateAccount :one
      INSERT INTO accounts (
        owner,
        balance,
        currency
      ) VALUES (
        $1, $2, $3
      ) RETURNING *;
      ```
