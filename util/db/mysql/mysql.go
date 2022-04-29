package mysql

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Repo interface {
	GetDb() *gorm.DB // 获取 DB 链接
	DbClose() error  // 关闭 DB 链接
}

type dbRepo struct {
	db *gorm.DB
}

func (d *dbRepo) GetDb() *gorm.DB {
	return d.db
}

func (d *dbRepo) DbClose() error {
	sqlDb, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDb.Close()
}

type Mysql struct {
	Host            string
	Port            int
	User            string
	Pass            string
	DbName          string
	MaxIdleConn     int           // 设置最大连接数，用于设置闲置的连接数
	MaxOpenConn     int           // 设置连接池，用于设置最大打开的连接数，默认值为0表示不限制
	MaxLifetimeConn time.Duration // 设置最大连接超时
}

func (m Mysql) Client() (Repo, error) {
	db, err := m.dbConnect()
	if err != nil {
		return &dbRepo{db: db}, err
	}
	return &dbRepo{db: db}, nil
}

func (m Mysql) dbConnect() (*gorm.DB, error) {
	var (
		err   error
		db    *gorm.DB
		sqlDb *sql.DB
		DSN   string
	)
	DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", m.User, m.Pass,
		m.Host, m.Port, m.DbName)
	obj := mysql.New(mysql.Config{
		DSN:                       DSN,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	})
	db, err = gorm.Open(obj, &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("[db connection failed] Database name: %s", m.DbName))
	}
	sqlDb, err = db.DB()
	if err != nil {
		return nil, err
	}
	sqlDb.SetMaxOpenConns(m.MaxOpenConn)
	sqlDb.SetMaxIdleConns(m.MaxIdleConn)
	sqlDb.SetConnMaxLifetime(time.Minute * m.MaxLifetimeConn)
	return db, nil
}
