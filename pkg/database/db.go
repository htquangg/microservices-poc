package database

import (
	"context"
	"database/sql"
	"time"

	"xorm.io/builder"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

type DB struct {
	ctx context.Context
	cfg *Config
	e   *xorm.Engine
}

func New(ctx context.Context, cfg *Config) (*DB, error) {
	xormEngine, err := NewXORMEngine(cfg)
	if err != nil {
		return nil, err
	}

	xormEngine.SetMapper(names.GonicMapper{})
	xormEngine.ShowSQL(cfg.LogSQL)
	xormEngine.SetMaxOpenConns(cfg.MaxOpenConns)
	xormEngine.SetMaxIdleConns(cfg.MaxIdleConns)
	xormEngine.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime))
	xormEngine.SetTZDatabase(time.UTC)
	xormEngine.SetDefaultContext(ctx)

	return &DB{
		e:   xormEngine,
		cfg: cfg,
		ctx: &Context{
			Context: ctx,
			e:       xormEngine,
		},
	}, nil
}

// Engine will get a db Engine from this context or return an Engine restricted to this context
func (db *DB) Engine(ctx context.Context) Engine {
	if e := db.engine(ctx); e != nil {
		return e
	}
	return db.e.Context(ctx)
}

// engine will get a db Engine from this context or return nil
func (db *DB) engine(ctx context.Context) Engine {
	if engined, ok := ctx.(Engined); ok {
		return engined.Engine()
	}
	enginedInterface := db.ctx.Value(enginedContextKey)
	if enginedInterface != nil {
		return enginedInterface.(Engined).Engine()
	}
	return nil
}

// TxContext represents a transaction Context,
// it will reuse the existing transaction in the parent context or create a new one.
func (db *DB) TxContext(parentCtx context.Context) (*Context, Committer, error) {
	if sess, ok := db.inTransaction(parentCtx); ok {
		return newContext(parentCtx, sess, true), &halfCommitter{committer: sess}, nil
	}

	sess := db.e.NewSession()
	if err := sess.Begin(); err != nil {
		return nil, nil, err
	}
	return newContext(db.ctx, sess, true), sess, nil
}

// WithTx represents executing database operations on a transaction, if the transaction exist,
// this function will reuse it otherwise will create a new one and close it when finished.
func (db *DB) WithTx(parentCtx context.Context, f func(ctx context.Context) error) error {
	if sess, ok := db.inTransaction(parentCtx); ok {
		err := f(newContext(parentCtx, sess, true))
		if err != nil {
			// rollback immediately, in case the caller ignores returned error and tries to commit the transaction.
			_ = sess.Close()
		}
		return err
	}
	return db.txWithNoCheck(parentCtx, f)
}

func (db *DB) txWithNoCheck(parentCtx context.Context, f func(ctx context.Context) error) error {
	sess := db.e.NewSession()
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}

	if err := f(newContext(parentCtx, sess, true)); err != nil {
		return err
	}

	return sess.Commit()
}

func (db *DB) Insert(ctx context.Context, beans ...any) error {
	_, err := db.Engine(ctx).Insert(beans...)
	return err
}

func (db *DB) Exec(ctx context.Context, sqlAndArgs ...any) (sql.Result, error) {
	return db.Engine(ctx).Exec(sqlAndArgs...)
}

func (db *DB) GetByBean(ctx context.Context, bean any) (bool, error) {
	return db.Engine(ctx).Get(bean)
}

func (db *DB) DeleteByBean(ctx context.Context, bean any) (int64, error) {
	return db.Engine(ctx).Delete(bean)
}

func (db *DB) DeleteByID(ctx context.Context, id int64, bean any) (int64, error) {
	return db.Engine(ctx).ID(id).NoAutoCondition().NoAutoTime().Delete(bean)
}

func (db *DB) FindIDs(ctx context.Context, tableName, idCol string, cond builder.Cond) ([]int64, error) {
	ids := make([]int64, 0, 10)
	if err := db.Engine(ctx).Table(tableName).
		Cols(idCol).
		Where(cond).
		Find(&ids); err != nil {
		return nil, err
	}
	return ids, nil
}

func (db *DB) DecrByIDs(ctx context.Context, ids []int64, decrCol string, bean any) error {
	_, err := db.Engine(ctx).Decr(decrCol).In("id", ids).NoAutoCondition().NoAutoTime().Update(bean)
	return err
}

func (db *DB) DeleteBeans(ctx context.Context, beans ...any) (err error) {
	e := db.Engine(ctx)
	for i := range beans {
		if _, err = e.Delete(beans[i]); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) TruncateBeans(ctx context.Context, beans ...any) (err error) {
	e := db.Engine(ctx)
	for i := range beans {
		if _, err = e.Truncate(beans[i]); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) CountByBean(ctx context.Context, bean any) (int64, error) {
	return db.Engine(ctx).Count(bean)
}

func (db *DB) InTransaction(ctx context.Context) bool {
	_, ok := db.inTransaction(ctx)
	return ok
}

func (db *DB) inTransaction(ctx context.Context) (*xorm.Session, bool) {
	e := db.engine(ctx)
	if e == nil {
		return nil, false
	}

	switch t := e.(type) {
	case *xorm.Engine:
		return nil, false
	case *xorm.Session:
		if t.IsInTx() {
			return t, true
		}
		return nil, false
	default:
		return nil, false
	}
}
