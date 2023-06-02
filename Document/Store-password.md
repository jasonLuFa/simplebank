# Store password 
## Bcrypt
### 過程
#### 儲存過程
1. 藉由 bcrypt funtion 加密 ( cost 和 salt ) 成 hash string
2. 將 hash string 存到 DB

#### 驗證過程
1. 將 hash string 從 DB 找出來
2. 將 request 過來的密碼使用 bcrypt 的比較方法和 DB 的 hash string 做比較 

### 介紹
- `$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy`
	- `alg` -> `2a`  : The hash algorithm identifier (Bcrypt)
	- `cost` -> `10` : Cost factor ( $2^{10}$ = 1,024 rounds of key expansion)，數字越大安全性越高，但會需要越多間，須隨科技機器運算能力進步而增加
	- `salt` -> `N9qo8uLOickgx2ZMRZoMye` : 16-byte (128-bit) salt, base64 encoded to 22 characters，隨機字元到密碼中 ( 所以即使密碼一樣，hash 出來的值也會不一樣，可以避免 rainbow table attack )
	- `hash` -> `IjZAgcfl7p92ldGxad68LJZdL17lhWy` : 24-byte (192-bit) hash, base64 encoded to 31 characters

### example
- 使用 package

```go
import (
"fmt"
"golang.org/x/crypto/bcrypt"
)
```
- hashPassword

```go
func HashPassword(password string) (string, error) {
hashedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
if err != nil {
return "", fmt.Errorf("failed to hash password: %w", err)
}
return string(hashedPassword),nil
}
``` 
- checkPassword

```go
func CheckPassword(password string, hashedPassword string) error{
return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
```

## Other implementation
- argon2
