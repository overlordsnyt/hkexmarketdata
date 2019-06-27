package _type

import (
	"fmt"
	"hkexgo/functions"
	"strconv"
	"strings"
)

type Stock struct {
	Rank               int
	StockCode          string
	StockName          string
	TodayIncome        float64
	LastTradeDayIncome float64
}

type StockTable []Stock

func NewStockTable(adtb *[][]string) *StockTable {
	stktb := make(StockTable, len(*adtb))
	for i, v := range *adtb {
		stk := new(Stock)
		stk.Rank, _ = strconv.Atoi(v[0])
		stk.StockCode = fmt.Sprintf("%06s", v[1])
		stk.StockName = strings.TrimRight(v[2], "　")
		stk.TodayIncome = *functions.CalculatePureIncomeDevideYi(&v[3], &v[4])
		stktb[i] = *stk
	}
	return &stktb
}

func (stktb *StockTable) SetLastTradeDayIncome(index *int, income *float64) {
	(*stktb)[*index].LastTradeDayIncome = *income
}

func (stk *Stock) SetLastTradeDayIncome(income *float64) {
	stk.LastTradeDayIncome = *income
}
