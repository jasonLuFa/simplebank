# 指令建立相對應 go 指令
- 範例
	```
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	  proto/*.proto
	```
	- `--go_out` 和 `--go_opt=paths=source_relative` : protobuf message 生成 go code 到指定位置
	- `--go-grpc_out` 和 `--go-grpc_opt=paths=source_relative` :  protobuf service 生成 go code 到指定位置
	- `proto/*.proto` : 基於 proto 資料夾底下的所有 `.proto` 檔去生成 go code
