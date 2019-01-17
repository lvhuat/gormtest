package main

import (
	"test/gormtest/dto"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

type CpOrderCursor struct {
	Symbol      string `gorm:"PRIMARY_KEY"`
	Value       string
	SyncVersion int64
	UpdatedAt   int64
}

func (cpOrderCursor *CpOrderCursor) Load() error {
	if err := db.Where("symbol = ?", cpOrderCursor.Symbol).Find(cpOrderCursor).Error; err != nil {
		return err
	}

	return nil
}

func (cpOrderCursor *CpOrderCursor) AllOrderSyncCursor() ([]*CpOrderCursor, error) {
	var cursors []*CpOrderCursor
	if err := db.Find(&cursors).Error; err != nil {
		return nil, err
	}

	return cursors, nil
}

func NewCpOrderCursor(symbol string) *CpOrderCursor {
	return &CpOrderCursor{
		Symbol: symbol,
	}
}

func (cpOrderCursor *CpOrderCursor) Save() error {
	if err := db.Save(cpOrderCursor).Error; err != nil {
		return err
	}

	return nil
}

func (cpOrderCursor *CpOrderCursor) AquireUpdateKey() (bool, error) {
	retDB := db.Model(cpOrderCursor).
		Where("updated_at=?", cpOrderCursor.UpdatedAt).
		Update("updated_at", time.Now().UnixNano()/int64(time.Millisecond))
	if err := retDB.Error; err != nil {
		return false, err
	}

	return retDB.RowsAffected == 1, nil
}

func test1() {
	d, err := gorm.Open("mysql", "root:lw123456@tcp(127.0.0.1:3306)/cryptobroker_test")
	if err != nil {
		panic(err)
	}
	db = d
	defer db.Close()

	db.DB().SetMaxIdleConns(3)
	db.DB().SetMaxOpenConns(50)

	cursor := NewCpOrderCursor("BTC_USDT")
	db.Set("gorm:table_options", "CHARSET=utf8").AutoMigrate(&CpOrderCursor{})

	if err := cursor.Save(); err != nil {
		panic(err)
	}
	_, err = cursor.AquireUpdateKey()
	if err != nil {
		panic(err)
	}
}

func test2() {
	if err := dto.Init("root:lw123456@tcp(127.0.0.1:3306)/cryptobroker_test"); err != nil {
		panic(err)
	}
	cursor := dto.NewCpOrderCursor("BTC_USDT")
	if err := cursor.Save(); err != nil {
		panic(err)
	}
	_, err := cursor.AquireUpdateKey()
	if err != nil {
		panic(err)
	}
}

func main() {
	test1()
	//test2()
	// Output:
	// panic: DB_ERROR,Error 1265: Data truncated for column 'updated_at' at row 1
	//
	// goroutine 1 [running]:
	// main.test2()
	//         D:/src/test/gormtest/main.go:97 +0x10e
	// main.main()
	//         D:/src/test/gormtest/main.go:103 +0x27
}
