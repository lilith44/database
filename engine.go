package database

import (
	"time"

	"github.com/lilith44/easy"
	xl "github.com/lilith44/xorm-logger"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type Database struct {
	engine    *xorm.Engine
	snowflake *easy.Snowflake
}

type JoinFunc func(engine *xorm.Engine) *xorm.Session

func New(c Config, logger *xl.ZapLogger, snowflake *easy.Snowflake) (*Database, error) {
	engine, err := xorm.NewEngine(c.Type, c.DSN())
	if err != nil {
		return nil, err
	}

	engine.ShowSQL(c.ShowSQL)
	engine.SetMaxOpenConns(c.Connection.MaxOpen)
	engine.SetMaxIdleConns(c.Connection.MaxIdle)
	engine.SetConnMaxLifetime(c.Connection.MaxLifetime)
	engine.SetTZLocation(time.Local)
	engine.SetLogger(logger)

	return &Database{engine: engine, snowflake: snowflake}, nil
}

// Engine 返回engine对象
func (d *Database) Engine() *xorm.Engine {
	return d.engine
}

// Transaction 事务操作
func (d *Database) Transaction(f func(*xorm.Session) (any, error)) (any, error) {
	return d.engine.Transaction(f)
}

// Insert 插入多条数据
func (d *Database) Insert(session *xorm.Session, beans ...easy.Model) (int64, error) {
	for _, bean := range beans {
		if bean.UseSnowflakeIdAsDefault() && bean.PK() == 0 {
			bean.SetPK(d.snowflake.NextId())
		}
	}

	return session.Insert(easy.ToAnySlice(beans)...)
}

// Delete 删除数据，bean的非零值字段将作为等于条件
func (d *Database) Delete(session *xorm.Session, beans ...easy.Model) (int64, error) {
	return session.Delete(easy.ToAnySlice(beans)...)
}

// DeleteByCond 根据cond删除数据，bean的非零值字段也会作为等于条件
func (d *Database) DeleteByCond(session *xorm.Session, bean easy.Model, cond builder.Cond) (int64, error) {
	return session.Where(cond).Delete(bean)
}

// UpdateById 根据主键更新数据，bean的非零值字段将进行更新
func (d *Database) UpdateById(session *xorm.Session, bean easy.Model, cols ...string) (int64, error) {
	return session.ID(bean.PK()).Cols(cols...).Update(bean)
}

// UpdateByCond 根据cond更新数据，bean的非零值字段将进行更新
func (d *Database) UpdateByCond(session *xorm.Session, bean easy.Model, cond builder.Cond, cols ...string) (int64, error) {
	return session.Where(cond).Cols(cols...).Update(bean)
}

// IncrById 根据主键进行字段加减，value为空切片时表示+1
func (d *Database) IncrById(session *xorm.Session, bean easy.Model, column string, value ...any) (int64, error) {
	return session.ID(bean.PK()).Incr(column, value...).Update(bean)
}

// IncrByCond 根据cond进行字段加减，value为空切片时表示+1
func (d *Database) IncrByCond(
	session *xorm.Session,
	bean easy.Model,
	cond builder.Cond,
	column string, value ...any,
) (int64, error) {
	return session.Where(cond).Incr(column, value...).Update(bean)
}

// Get 查询，bean的非零值字段会作为等于条件
func (d *Database) Get(bean easy.Model, cols ...string) (bool, error) {
	return d.engine.Cols(cols...).Get(bean)
}

// GetByCond 根据cond查询，bean的非零值字段也会作为等于条件
func (d *Database) GetByCond(bean easy.Model, cond builder.Cond, cols ...string) (bool, error) {
	return d.engine.Where(cond).Cols(cols...).Get(bean)
}

// Exist 查询是否存在，bean的非零值字段会作为等于条件。与Get相比，不会填充bean的值
func (d *Database) Exist(bean easy.Model, cols ...string) (bool, error) {
	return d.engine.Cols(cols...).Exist(bean)
}

// ExistByCond 根据cond查询是否存在，bean的非零值字段也会作为等于条件。与GetByCond相比，不会填充bean的值
func (d *Database) ExistByCond(bean easy.Model, cond builder.Cond, cols ...string) (bool, error) {
	return d.engine.Where(cond).Cols(cols...).Exist(bean)
}

// Count 计数，bean的非零值字段也会作为等于条件
func (d *Database) Count(bean easy.Model) (int64, error) {
	return d.engine.Count(bean)
}

// CountByCond 根据cond进行计数，bean的非零值字段也会作为等于条件
func (d *Database) CountByCond(bean easy.Model, cond builder.Cond) (int64, error) {
	return d.engine.Where(cond).Count(bean)
}

// Find 列表查询，bean的非零值字段会作为等于条件
func (d *Database) Find(beansPtr any, bean easy.Model, join ...JoinFunc) error {
	if len(join) != 0 {
		return join[0](d.engine).Find(beansPtr, bean)
	}
	return d.engine.Find(beansPtr, bean)
}

// FindByCond 根据cond进行列表查询，bean的非零值字段会作为等于条件
func (d *Database) FindByCond(beansPtr any, bean easy.Model, cond builder.Cond, join ...JoinFunc) error {
	if len(join) != 0 {
		return join[0](d.engine).Where(cond).Find(beansPtr, bean)
	}
	return d.engine.Where(cond).Find(beansPtr, bean)
}

// FindByPaging 分页查询
func (d *Database) FindByPaging(beansPtr any, paging Pager, join ...JoinFunc) (int64, error) {
	if len(join) != 0 {
		return join[0](d.engine).Where(paging.Cond()).OrderBy(paging.OrderBy()).Limit(paging.Limit()).FindAndCount(beansPtr)
	}
	return d.engine.Where(paging.Cond()).OrderBy(paging.OrderBy()).Limit(paging.Limit()).FindAndCount(beansPtr)
}
