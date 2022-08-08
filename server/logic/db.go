package logic

import (
	"database/sql"
	"fmt"
	"mygrpc/proto"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DefaultDB *ProcessSt

type ProcessSt struct {
	*sql.DB
}

func Init(cfg *proto.MysqlSt) (err error) {
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
	DefaultDB = &ProcessSt{
		db,
	}
	return
}

func Exchange(from, fromMoney, to, toMoney int64, uuid string) (err error) {
	// 开启事务
	tx, err := DefaultDB.Begin()
	if err != nil {
		fmt.Println("db tx err:", err)
		return
	}
	// 检查唯一key
	sql := fmt.Sprintf("insert into unique_key values('%s')", uuid)
	_, err = tx.Exec(sql)
	if err != nil {
		fmt.Println("unique_key err", err)
		tx.Rollback()
		return
	}
	// 插入订单
	sql = fmt.Sprintf("insert into `order` (`from`, from_money, `to`, to_money) values(%d, %d, %d, %d)", from, fromMoney, to, toMoney)
	_, err = tx.Exec(sql)
	if err != nil {
		fmt.Println("order err", err)
		tx.Rollback()
		return
	}
	// 扣除
	sql = fmt.Sprintf("insert into account(id, money) values(%d, %d) on duplicate key update money=money+%d", from, fromMoney, fromMoney)
	_, err = tx.Exec(sql)
	if err != nil {
		fmt.Println("account from err", err)
		tx.Rollback()
		return
	}
	// 增加
	sql = fmt.Sprintf("insert into account(id, money) values(%d, %d) on duplicate key update money=money+%d", to, toMoney, toMoney)
	_, err = tx.Exec(sql)
	if err != nil {
		fmt.Println("account to err", err)
		tx.Rollback()
		return
	}
	// 查询转账后余额
	sql = fmt.Sprintf("select money from account where id=%d", from)
	r := tx.QueryRow(sql)
	if r == nil {
		fmt.Println("account query err", err)
		tx.Rollback()
		return
	}
	// 映射查询结果
	var money int64
	err = r.Scan(&money)
	if err != nil {
		fmt.Println("account scan err", err)
		tx.Rollback()
		return
	}
	if money < 0 {
		err = fmt.Errorf("money not enough")
		fmt.Println("money err", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}
