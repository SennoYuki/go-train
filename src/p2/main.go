package main

import (
	"fmt"

	"go-train/src/p2/dao"
)

/*
第二章：异常处理
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

分解成两个问题
1. 碰到 sql.ErrNoRows 是否需要返回错误?
	这个视具体业务场景定,
	比如有一张表可能用户不具有此属性就不写入, 而函数的定义也将空值视为合法的返回值, 那么直接干掉这个 error
	如果函数内部确定这里应该有这行数据, 但是查询返回sql.ErrNoRows, 那么应该返回 error
	在空值为合法值的情况下必须返回 error, 这样上层才能确定这个函数执行是否出错
	在 函数内不确定 sql.ErrNoRows 是否合法的情况下应该返回 error 交给上层去处理
2. 返回 sql.ErrNoRows 是否需要 wrap?
	yes, 这样可以帮助排查问题
*/
func main() {
	dao.
}
