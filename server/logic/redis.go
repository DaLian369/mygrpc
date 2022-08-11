package logic

import (
	"fmt"
	"log"
)

func (r *RedisManager) genUUIDKey(id int64) string {
	return fmt.Sprintf("exchange_%d", id)
}

// 也可以改成map存储
func (r *RedisManager) getExchangeKey(id int64, uuid string) string {
	return fmt.Sprintf("exchange_%d_%s", id, uuid)
}

func (r *RedisManager) DelExchangeKey(id int64, uuid string) (err error) {
	// 改成lua脚本保证失败状态才会解开，成功后当作缓存使用不删除
	// return r.Del(r.getExchangeKey(id, uuid))

	lua := `
if redis.call("get", KEYS[1]) == ARGV[1] then
	redis.call("del", KEYS[1])
else 
	return 0
end
	`
	_, err = r.Eval(lua, 1, r.getExchangeKey(id, uuid), EXCHANGE_STATUS_WAIT)
	return err
}

func (r *RedisManager) CheckExchangeKey(id int64, uuid string) (lock bool) {
	key := r.getExchangeKey(id, uuid)
	err := r.Set(key, EXCHANGE_STATUS_WAIT, "EX", 3, "NX") // 这里的作用是互斥锁
	if err != nil {
		log.Printf("CheckExchangeKey err: %v", err)
		return
	}
	lock = err == nil
	return
}

func (r *RedisManager) SuccExchangeKey(id int64, uuid string) (err error) {
	key := r.getExchangeKey(id, uuid)
	err = r.Set(key, EXCHANGE_STATUS_SUCC, "EX", 60*100) // 这里的作用是缓存
	if err != nil {
		log.Printf("SuccExchangeKey err: %v", err)
		return
	}
	// uuid转账成功后与用户解除绑定
	key = r.genUUIDKey(id)
	err = r.SafeDel(key, uuid)
	if err != nil {
		log.Printf("del uuid failed: %v", err)
		err = nil
	}
	return
}

func (r *RedisManager) SafeDel(key string, value interface{}) (err error) {
	lua := `
if redis.call("get", KEYS[1]) == ARGV[1] then
	redis.call("del", KEYS[1])
else 
	return 0
end
	`
	_, err = r.Eval(lua, 1, key, value)
	return
}
