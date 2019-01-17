package dto

import (
	"cryptobroker/reportservice/common"
	"time"
)

type CpOrderCursor struct {
	Symbol      string `gorm:"PRIMARY_KEY"`
	Value       string
	SyncVersion int64
	UpdatedAt   int64
}

func (cpOrderCursor *CpOrderCursor) Load() error {
	if err := db.Where("symbol = ?", cpOrderCursor.Symbol).Find(cpOrderCursor).Error; err != nil {
		return dbError2CodeError(err)
	}

	return nil
}

func (cpOrderCursor *CpOrderCursor) AllOrderSyncCursor() ([]*CpOrderCursor, error) {
	var cursors []*CpOrderCursor
	if err := db.Find(&cursors).Error; err != nil {
		return nil, dbError2CodeError(err)
	}

	return cursors, nil
}

func NewCpOrderCursor(symbol string) *CpOrderCursor {
	return &CpOrderCursor{
		Symbol: symbol,
	}
}

func (cpOrderCursor *CpOrderCursor) Save() error {
	cpOrderCursor.UpdatedAt = millisec(time.Now())
	if err := db.Save(cpOrderCursor).Error; err != nil {
		return dbError2CodeError(err)
	}

	return nil
}

func (cpOrderCursor *CpOrderCursor) AquireUpdateKey() (bool, error) {
	retDB := db.Model(cpOrderCursor).
		Where("updated_at=?", cpOrderCursor.UpdatedAt).
		Update("updated_at", common.Millisec(time.Now()))
	if err := retDB.Error; err != nil {
		return false, dbError2CodeError(err)
	}

	return retDB.RowsAffected == 1, nil
}
