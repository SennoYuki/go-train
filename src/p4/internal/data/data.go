package data

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet()
