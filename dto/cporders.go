package dto

import (
	"cryptobroker/reportservice/common"
	"time"
)

type CpOrder struct {
	OrderId               string `gorm:"PRIMARY_KEY"`
	Symbol                string `gorm:"INDEX:scan_index"`
	Side                  string
	Type                  string
	Amount                string
	Price                 string
	State                 string
	CreatedAt             int64 `gorm:"INDEX:scan_index"`
	FilledAmount          string
	FilledAccumulativeQty string
	Fee                   string
	UpdatedAt             int64
}

func (cpOrder *CpOrder) Save() error {
	if err := db.Save(cpOrder).Error; err != nil {
		return dbError2CodeError(err)
	}

	return nil
}

func (cpOrder *CpOrder) BetweenTime(startAt, endAt int64) ([]*CpOrder, error) {
	var orders []*CpOrder
	if err := db.Where("created_at > ? and created_at < ?").Find(&orders).Error; err != nil {
		return nil, dbError2CodeError(err)
	}
	return nil, nil
}

func (cpOrder *CpOrder) BySymbolBefore25H() ([]*CpOrder, error) {
	if cpOrder.Symbol == "" {
		panic("bad symbol paramter")
	}
	var orders []*CpOrder
	if err := db.Where(
		"created_at > ? and symbol=?",
		common.Millisec(time.Now().Add(-time.Hour*25)),
		cpOrder.Symbol,
	).Find(&orders).Error; err != nil {
		return nil, dbError2CodeError(err)
	}
	return orders, nil
}
