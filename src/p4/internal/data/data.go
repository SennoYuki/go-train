package data

import (
	"database/sql"
	"sync"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewData)

type Data struct {
	db  *sql.DB
	rdb *redis.Client
}

// TODO:: 1,配置 2,初始化多个client
func NewData() (data *Data, err error) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		data.db, err = sql.Open("mysql", "")
		wg.Done()
	}()
	data.rdb = redis.NewClient(nil)
	wg.Done()
	if err != nil {
		data = nil
	}
	return
}
