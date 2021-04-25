package dao

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var db *sql.DB

func InitDb() (err error) {
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/train?charset=utf8mb4&timeout=5s")
	return
}

// GetUserNum 查询用户数量
// 查询满足某条件数量的 sql 一般不需要 wrap 和return error, 因为这是符合预期的
// 不过实际执行count函数也不会返回 sql.ErrNoRows
func GetUserNum(arg ...interface{}) (num int, err error) {
	query := "select count(1) from user"
	//query := "select 1 from user where 0 limit 1"
	if err = db.QueryRow(query).Scan(&num); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.Wrapf(err, "query error||sql=%s||arg=%v", query, arg)
	}
	return
}

// GetUserIdentityNumber 查询用户身份证号码
// 假设这里 QueryRow 返回一个 sql.ErrNoRows, 则number一定为一个空字符串, 明显是非法值,
// 为了在调用方可以只根据 error 而不是返回值处理, 或者说保证 error!=nil 那么这个查询是ok的, 所以在这里包装并返回 error
func GetUserIdentityNumber(userId int64) (number string, err error) {
	query := "select identity_number from user where user_id=?"
	//query := "select 1 from user where 0 limit 1"
	if err = db.QueryRow(query, userId).Scan(&number); err != nil {
		err = errors.Wrap(err, "query error")
	}
	err = errors.Wrapf(err, "query error||sql=%s||arg=%v", query, userId)
	return
}

func GetSomething() (something string, err error) {
	something, err = GetSomethingA()
	if err != nil {
		err = errors.WithMessage(err, "handle GetSomethingA")
	}
	return
}

// GetSomethingA 假设空字符串是 GetSomethingA 的一个非法返回值, 那么必须 wrap 并返回 error
func GetSomethingA() (something string, err error) {
	query := "select 1 from user where 0 limit 1"
	if err = db.QueryRow(query).Scan(&something); err != nil {
		err = errors.Wrapf(err, "query error||sql=%s||arg=%v", query, interface{}(nil))
	}
	return
}

// GetSomethingB sql.ErrNoRows 符合预期时不向上返回 error
func GetSomethingB() (something string, err error) {
	query := "select 1 from user where 0 limit 1"
	if err = db.QueryRow(query).Scan(&something); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		err = errors.Wrapf(err, "query error||sql=%s||arg=%v", query, interface{}(nil))
	}
	return
}
