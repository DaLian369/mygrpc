# mygrpc
使用grpc实现简单转账功能，另外引入配置，go pprof，mysql，事务，缓存，连接池，消息队列，日志，监控，websocket等。

## 思路
首先模拟客户端点击进入转账页面获取一个令牌uuid，再调用转账Exchange接口时将令牌uuid传过来，数据库插入唯一key，保证幂等性。

【问题一】最开始是每次进来随机生成一个令牌uuid，那么这个用户重复进来可以获取多个令牌进行转账使用。

为了保证客户端多次请求只能获取相同未消费的令牌uuid，所以引入redis缓存，每次请求检查是否已有令牌，如果有返回，没有的话加`set ex nx`互斥设置。

【问题二】转账接口Exchange里没有使用缓存对请求的令牌进行互斥校验，导致高并发情况下大量请求落到数据库，数据库压力大，并且造成死锁错误。

解决办法是通过在转账接口的入口处对令牌`set ex nx`设置互斥，保证n个请求只有一个能落到db，完成业务流程，大大减小了数据库压力。

【问题三】在某个令牌使用后，后续重复使用这个令牌会进入db，开启事务，在插入唯一key时才会报错返回。

解决办法是使用redis缓存将使用过的token缓存起来，后续进来先判断有没有使用过。



## 进度
[f] 配置文件

[f] go pprof

[f] mysql、事务

[f] 缓存、连接池

[f] 消息队列 kafka

[w] websocket

[f] 日志

[f] 监控 prometheus

[f] 测试

## server
启动`go run server/main.go -c conf/default.yaml`

## client
运行`go run client/main.go`

## proto
运行`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/my.proto`

## mysql
创建表sql`/server/sql/sql.sql`

## redis
批量删除前缀key方法：

keys+xargs批量删除:
`redis-cli -h 127.0.0.1 -p 6379 keys "ex*" | xargs redis-cli -h 127.0.0.1 -p 6379 del`

scan+xargs批量删除:
`redis-cli -h 127.0.0.1 -p 6379 --scan --pattern "ex*" | xargs redis-cli -h 127.0.0.1 -p 6379 del`

坑：

1、`conn.Do()`方法的返回结果需要使用`redis.String()`或其他方法包装后返回。

2、互斥锁解开使用lua脚本`if redis.call("get", KEYS[1]) == ARGV[1] then redis.call("del", KEYS[1]) else return 0 end`

3、redis lua脚本里的 KEYS和ARGV必须大写，否则报错识别不出来keys。
```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
	redis.call("del", KEYS[1])
else 
	return 0
end
```

## go pprof
url `http://127.0.0.1:9000/debug/pprof/`

## prometheus
url `http://127.0.0.1:9000/metrics`

## 测试
```sh
cd server/logic
go test -v
```

