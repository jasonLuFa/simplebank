# 💻 [PASETO](github.com/o1egl/paseto) vs [JWT](https://jwt.io/)

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

## 💻 JWT

- 所有資料存在 JWT 只有 base64 編碼，並非加密，所以避免存取重要資訊
- token 可分為 Header, payload, verify singature
  - Header : signing algorithm
  - payload : 關於登入者的資訊、包括 token 過期時間 ... 你想存的資訊
  - verify singature : 數位簽章( digital signature )

## 💻 PASETO( platform agnostic security tokens )

- 跟 JWT 類似，但是比 JWT 更安全及容易實作
