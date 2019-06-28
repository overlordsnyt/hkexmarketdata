package _type

import (
	"fmt"
	"hkexgo/calculator"
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
		stk.TodayIncome = *calculator.CalculatePureIncomeDevideYi(&v[3], &v[4])
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

func (stktb *StockTable) Len() int {
	return len(*stktb)
}

func (stktb *StockTable) Less(i, j int) bool {
	return (*stktb)[i].TodayIncome < (*stktb)[j].TodayIncome
}

func (stktb *StockTable) Swap(i, j int) {
	(*stktb)[i], (*stktb)[j] = (*stktb)[j], (*stktb)[i]
	(*stktb)[i].Rank, (*stktb)[j].Rank = (*stktb)[j].Rank, (*stktb)[i].Rank
}
