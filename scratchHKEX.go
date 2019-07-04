package main

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/tealeg/xlsx"
	"hkexgo/configuration"
	"hkexgo/dealer"
	"hkexgo/excel"
	"hkexgo/type"
	_ "image/png"
	"os"
)

var config *configuration.Configuration

func main() {
	var assignDate, lastTradeDate string
	fmt.Print("assign date: ")
	fmt.Scanln(&assignDate)
	fmt.Print("last trade date: ")
	fmt.Scanln(&lastTradeDate)

	assignDateTop10, _ := dealer.GetHKEXJson(assignDate)
	if assignDateTop10 == nil {
		return
	}

	lastTradeDateTop10, _ := dealer.GetHKEXJson(lastTradeDate)
	if lastTradeDateTop10 == nil {
		return
	}

	config = configuration.LoadConfiguration()

	hkTableSearchMap := dealer.ScrachAssignDateTop10JsonToMarketNameSearchMap(&assignDateTop10)

	lastTradeCodeIncomeSearchMap := dealer.GenerateStockCodePureIncomeSearchMapFromLastTradeDateJson(&lastTradeDateTop10)

	hkMergedTwoDaysTable := dealer.MergeLastTradeIncomeToAssignDateTable(hkTableSearchMap, lastTradeCodeIncomeSearchMap)

	hkSortedTable := dealer.SortAllMarketTable(hkMergedTwoDaysTable)

	dateMD := dealer.TransferMonthDay(&assignDate)

	filename := fmt.Sprint(assignDate, ".xlsx")

	generateXLSX(hkSortedTable, &filename, dateMD)
	insertPNG(&filename)

	fmt.Printf("\nsratch success!\nsaved as '%v'\n", filename)

	enterClose()
}

const SSEN, SZSEN, sheetName = "SSE Northbound", "SZSE Northbound", "沪深港通"

var headers = &[]string{"排名", "股票代码", "股票名称", "净买入（亿元）", "前一交易日净买入额（亿元）"}

func generateXLSX(hkTable *map[string]*_type.StockTable, filename *string, dateStr *string) {
	xlsxfile := xlsx.NewFile()
	sheet, _ := xlsxfile.AddSheet(sheetName)

	excel.GenerateTitle(sheet, 0, 0, 4, "沪股通（"+*dateStr+"）")
	excel.GenerateTitle(sheet, 0, 6, 4, "深股通（"+*dateStr+"）")
	excel.ChangeColWidth(sheet, 5, 2)
	excel.GenerateHeader(sheet, 1, 0, headers)
	excel.GenerateHeader(sheet, 1, 6, headers)
	excel.FillHKEXData(sheet, 2, 0, (*hkTable)[SSEN])
	excel.FillHKEXData(sheet, 2, 6, (*hkTable)[SZSEN])

	//设置列宽
	width := config.Width
	sheet.SetColWidth(0, 0, width.RankWidth)
	sheet.SetColWidth(1, 3, width.TodayWidth)
	sheet.SetColWidth(4, 4, width.LastTradeDayWidth)
	sheet.SetColWidth(6, 6, width.RankWidth)
	sheet.SetColWidth(7, 9, width.TodayWidth)
	sheet.SetColWidth(10, 10, width.LastTradeDayWidth)

	//保存excel文件
	xlsxfile.Save(*filename)
}

func insertPNG(fileName *string) {
	f, _ := excelize.OpenFile(*fileName)
	err := f.AddPicture(sheetName, "J8", "watermark.png",
		`{"x_offset": 10, "y_offset": 10, "print_obj": true, "lock_aspect_ratio": true, "locked": true, "positioning": "oneCell"}`)
	if err != nil {
		fmt.Println(err.Error())
	}
	f.Save()
}

func enterClose() {
	fmt.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
