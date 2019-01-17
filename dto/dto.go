package dto

import (
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lworkltd/kits/service/restful/code"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithField("pkg", "dto")

	initOnce sync.Once
	inited   bool
	db       *gorm.DB
)

// Init 初始化
func Init(url string) error {
	if inited {
		panic("Don't init twice")
	}
	var err error
	initOnce.Do(func() {
		err = initMysql(url)
	})

	return err
}

func initMysql(url string) error {
	d, err := gorm.Open("mysql", url)
	if err != nil {
		return err
	}
	db = d
	db.DB().SetMaxIdleConns(3)
	db.DB().SetMaxOpenConns(50)

	initTables(&CpOrder{}, &CpOrderCursor{})

	log.Print("Init MySQL")

	return nil
}

func dbError2CodeError(err error) code.Error {
	if err.Error() == "record not found" {
		return code.NewMcodef("NOT_FOUND", err.Error())
	}

	if strings.Index(err.Error(), "Error 1062") != -1 {
		return code.NewMcodef("DUPLICATED", err.Error())
	}

	return code.NewMcodef("DB_ERROR", err.Error())
}

// IsNotFound 判断是否是一个未查到目标
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	cerr, ok := err.(code.Error)
	if !ok {
		return false
	}

	if cerr.Mcode() == "NOT_FOUND" {
		return true
	}

	return false
}

func IsDuplicated(err error) bool {
	if err == nil {
		return false
	}
	cerr, ok := err.(code.Error)
	if !ok {
		return false
	}

	if cerr.Mcode() == "DUPLICATED" {
		return true
	}

	return false
}

func isDBNotFound(dbError error) bool {
	return dbError.Error() == "record not found"
}

func initTables(tables ...interface{}) {
	for _, table := range tables {
		if db.HasTable(table) {
			continue
		}
		db.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(table)
	}
}

// TranscationItem 事务操作
type TranscationItem interface {
	Save(*gorm.DB) error
}

// DoTranscations 将savers作为一个事务整体处理
func DoTranscations(savers []TranscationItem) error {
	// 单个不参与事务
	if len(savers) == 1 {
		return savers[0].Save(nil)
	}

	tx := db.Begin()
	for _, saver := range savers {
		if err := saver.Save(tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// 转换为毫秒
func millisec(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
