package main

import (
	"github.com/tealeg/xlsx"
	"hkexgo/dealer"
	"hkexgo/excel"
	"hkexgo/type"
)

func main() {

	assignDate := "2019-06-17"
	assignDateTop10, _ := dealer.GetHKEXJson(assignDate)
	if assignDateTop10 == nil {
		return
	}

	lastTradeDate := "2019-06-14"
	lastTradeDateTop10, _ := dealer.GetHKEXJson(lastTradeDate)
	if lastTradeDateTop10 == nil {
		return
	}

	hkTableSearchMap := dealer.ScrachAssignDateTop10JsonToMarketNameSearchMap(&assignDateTop10)

	lastTradeCodeIncomeSearchMap := dealer.GenerateStockCodePureIncomeSearchMapFromLastTradeDateJson(&lastTradeDateTop10)

	hkMergedTwoDaysTable := dealer.MergeLastTradeIncomeToAssignDateTable(hkTableSearchMap, lastTradeCodeIncomeSearchMap)

	hkSortedTable := dealer.SortAllMarketTable(hkMergedTwoDaysTable)

	//fmt.Println(hkSortedTable)

	GenerateXLSX(hkSortedTable)
}

//TODO 输出带基本样式的数据excel

const SSEN, SZSEN = "SSE Northbound", "SZSE Northbound"

var headers = &[]string{"排名", "股票代码", "股票名称", "净买入（亿元）", "前一交易日净买入额（亿元）"}

//const NUM_FORMAT = "0.0000_ " //尾空格，坑死人

func GenerateXLSX(hkTable *map[string]*_type.StockTable) {
	xlsxfile := xlsx.NewFile()
	sheet, _ := xlsxfile.AddSheet("沪深港通")

	excel.GenerateTitle(sheet, 0, 0, 4, "沪股通")
	excel.GenerateTitle(sheet, 0, 6, 4, "深股通")
	excel.ChangeColWidth(sheet, 5, 2)
	excel.GenerateHeader(sheet, 1, 0, headers)
	excel.GenerateHeader(sheet, 1, 6, headers)
	excel.FillHKEXData(sheet, 2, 0, (*hkTable)[SSEN])
	excel.FillHKEXData(sheet, 2, 6, (*hkTable)[SZSEN])

	//设置列宽
	sheet.SetColWidth(0, 0, 6.1)
	sheet.SetColWidth(1, 3, 11.7)
	sheet.SetColWidth(4, 4, 20.7)
	sheet.SetColWidth(6, 6, 6.1)
	sheet.SetColWidth(7, 9, 11.7)
	sheet.SetColWidth(10, 10, 20.7)

	//保存excel文件
	xlsxfile.Save("golanghkexcel.xlsx")
}
