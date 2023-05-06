# 介紹 
- 是一種支援同步( synchronous ) 及 非同步 asynchronous ) 通信方式
- gRPC : g (代表不同意思，[參考](https://github.com/grpc/grpc/blob/master/doc/g_stands_for.md)) Remote Procedure Calls 
	- 簡單來說是一個協議，允許執行( 例如呼叫 function )另一台電腦的程式( 可以是不同語言 )
	- client 和 server 使用 stubs 去做溝通
- 使用 protobuf 來定義服務和消息，從而生成強類型的客戶端和服務端代碼，減少了手動編寫代碼的錯誤風險

# Protocol
```protobuf
syntax = "proto3";

package pb;

option go_package = "<module_name_of_go_mod_init_plus_root_folder_name>";

message <NameOfTheMessage>{
  <data-type> <name_of_field_1> = tag_1;
  <data-type> <name_of_field_2> = tag_2;
  ...
  <data-type> <name_of_field_N> = tag_N;
}

service <NameOfTheService>{
	rpc <FunctionName>(<message-type>) return (<message-type>){}
	rpc <FunctionName>(steam <message-type>) return (<message-type>){}
	rpc <FunctionName>(<message-type>) return (steam <message-type>){}
	rpc <FunctionName>(steam <message-type>) return (steam <message-type>){}
}
```

# 優點
- 高效率( 低延遲 )
	- 因為 HTTP 協定關係，可以連線多工( multiplexing )，意思是能從單一 TCP Connect 中，同時傳輸處理多個 Request & Response
	- 且使用高效的二進制序列化格式 ( 相比 json 格式容量更小 )
	- 支持全局超時設置和區域性超時設置，這可以控制客戶端和服務器之間的等待時間
- 高可用 :
	- 支援多種語言，可以跨平臺使用
	- 提供了許多可靠性保障機制，例如消息確認、重試、斷點續傳等

# 缺點
- 不容易除錯
- 不容易實作

# 種類
## 1. Unary
- 跟 rest 一樣，一個 request 只回傳一個 response
-  protobuf
	```protobuf
	service GreetService{
		rpc Greet(GreetRequest) returns (GreetResponse){};
	}
	```
## 2. Server streaming
- 一個 request 可以有多個 response
-  protobuf
	```protobuf
	service GreetService{
		rpc Greet(GreetRequest) returns (stream GreetResponse){};
	}
	```

## 3. Client streaming
- 多個 request 並回傳一個 response
-  protobuf
	```protobuf
	service GreetService{
		rpc Greet(stream GreetRequest) returns (GreetResponse){};
	}
	```
## 4. Bi directional streaming
- 多個 request 和多個 reponse
-  protobuf
	```protobuf
	service GreetService{
		rpc Greet(stream GreetRequest) returns (stream GreetResponse){};
	}
	```
