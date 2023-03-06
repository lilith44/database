package database

import "xorm.io/builder"

type Pager interface {
	OrderBy() string
	Cond() builder.Cond
	Limit() (int, int)
}
