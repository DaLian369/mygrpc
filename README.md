# mygrpc
使用grpc实现简单转账功能，另外引入配置，go pprof，mysql，事务，缓存，连接池，消息队列，日志，监控，websocket等。

## server
启动`go run server/main.go -c conf/default.yaml`

## client
运行`go run client/main.go`

## mysql
创建表sql`/server/sql/sql.sql`