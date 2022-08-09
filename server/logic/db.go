package logic

import (
	"database/sql"
	"fmt"
	"log"
	"mygrpc/proto"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var DefaultProc *ProcessSt

type ProcessSt struct {
	db *sql.DB

	cache *RedisManager
}

func Init(cfg *proto.ConfigSt) (err error) {
	DefaultProc = &ProcessSt{}
	err = InitDB(&cfg.Mysql, DefaultProc)
	if err != nil {
		fmt.Println("init db err", err)
		return
	}
	err = InitCache(&cfg.Redis, DefaultProc)
	if err != nil {
		fmt.Println("init cache err", err)
		return
	}
	return
}

func InitDB(cfg *proto.MysqlSt, proc *ProcessSt) (err error) {
	connstr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cfg.Username, cfg.Password, cfg.Hostport, cfg.Database)
	// 创建连接
	db, err := sql.Open("mysql", connstr)
	if err != nil {
		fmt.Println("db init failed err:", err)
		return err
	}
	// 检查连接
	err = db.Ping()
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(cfg.Poolsize)        // 最大连接数
	db.SetMaxIdleConns(cfg.Idlesize)        // 最大空闲数
	db.SetConnMaxIdleTime(time.Minute * 10) // 连接最大存活时间
	proc.db = db
	return
}

func (p *ProcessSt) GenUUID(id int64) (uid string, err error) {
	// 查询是否有
	c := p.cache
	key := c.genUUIDKey(id)
	bs, err := c.GetData(key)
	if err != nil {
		return "", err
	}
	if len(bs) > 0 {
		uid = string(bs)
		return
	}
	// 如果没有设置令牌
	uid = uuid.New().String()
	err = c.Set(key, uid, "ex", 10*60, "nx")
	if err != nil {
		err = nil
		// 设置失败再查一次
		bs, err = c.GetData(key)
		if err != nil {
			return "", err
		}
		uid = string(bs)
	}
	if uid == "" {
		err = fmt.Errorf("gen uuid failed")
	}
	return
}

func (p *ProcessSt) Exchange(from, fromMoney, to, toMoney int64, uuid string) (err error) {
	// 加锁，校验uuid是否完成，在这里加入令牌状态解决并发进入数据库时的死锁错误，防止大量请求落到db
	succ := p.cache.CheckExchangeKey(from, uuid)
	if !succ {
		err = fmt.Errorf("exchange repeated")
		log.Printf("lock: %v", err)
		return
	}
	// 流程结束后，判断key状态进行解锁/缓存
	defer func() {
		p.cache.DelExchangeKey(from, uuid)
	}()

	// 开启事务
	db := p.db
	tx, err := db.Begin()
	if err != nil {
		log.Printf("db tx err: %v", err)
		return
	}
	// 检查唯一key
	// TODO 并发情况下会出现错误：Error 1213: Deadlock found when trying to get lock; try restarting transaction
	sql := fmt.Sprintf("insert into unique_key values('%s')", uuid)
	_, err = tx.Exec(sql)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			p.cache.SuccExchangeKey(from, uuid)
		}
		log.Printf("unique_key err: %v", err)
		tx.Rollback()
		return
	}
	// 插入订单
	sql = fmt.Sprintf("insert into `order` (`from`, from_money, `to`, to_money) values(%d, %d, %d, %d)", from, fromMoney, to, toMoney)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("order err: %v", err)
		tx.Rollback()
		return
	}
	// 扣除
	sql = fmt.Sprintf("insert into account(id, money) values(%d, %d) on duplicate key update money=money+%d", from, fromMoney, fromMoney)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("account from err: %v", err)
		tx.Rollback()
		return
	}
	// 增加
	sql = fmt.Sprintf("insert into account(id, money) values(%d, %d) on duplicate key update money=money+%d", to, toMoney, toMoney)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("account to err: %v", err)
		tx.Rollback()
		return
	}
	// 查询转账后余额
	sql = fmt.Sprintf("select money from account where id=%d", from)
	r := tx.QueryRow(sql)
	if r == nil {
		log.Printf("account query err: %v", err)
		tx.Rollback()
		return
	}
	// 映射查询结果
	var money int64
	err = r.Scan(&money)
	if err != nil {
		log.Printf("account scan err: %v", err)
		tx.Rollback()
		return
	}
	if money < 0 {
		err = fmt.Errorf("money not enough")
		log.Printf("money err: %v", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	p.cache.SuccExchangeKey(from, uuid)
	return
}
