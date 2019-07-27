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
	"sync"
	"time"
)

var config *configuration.Configuration

func main() {
	var assignDate, lastTradeDate string
	fmt.Print("assign date: ")
	fmt.Scanln(&assignDate)
	fmt.Print("last trade date: ")
	fmt.Scanln(&lastTradeDate)

	waitGroup := new(sync.WaitGroup)
	var assignDateTop10, lastTradeDateTop10 _type.Hkex

	go func() {
		assignDateTop10, _ = dealer.GetHKEXJson(assignDate, waitGroup)
	}()

	go func() {
		lastTradeDateTop10, _ = dealer.GetHKEXJson(lastTradeDate, waitGroup)
	}()

	//必须sleep才能接到WaitGroup done的同步信号 why???
	time.Sleep(200 * time.Nanosecond)
	waitGroup.Wait()

	if assignDateTop10 == nil {
		msg := fmt.Sprintf("assgin date not scratched any data.")
		enterClose(&msg)
		return
	}

	var hkTableSearchMap *map[string]map[string]_type.Table
	var lastTradeCodeIncomeSearchMap *map[string]map[string]float64
	go func() {
		waitGroup.Add(1)
		defer waitGroup.Done()
		hkTableSearchMap = dealer.ScrachAssignDateTop10JsonToMarketNameSearchMap(&assignDateTop10)
	}()
	go func() {
		waitGroup.Add(1)
		defer waitGroup.Done()
		lastTradeCodeIncomeSearchMap = dealer.GenerateStockCodePureIncomeSearchMapFromLastTradeDateJson(&lastTradeDateTop10)
	}()

	time.Sleep(200 * time.Nanosecond)
	waitGroup.Wait()

	config = configuration.LoadConfiguration()

	hkMergedTwoDaysTable := dealer.MergeLastTradeIncomeToAssignDateTable(hkTableSearchMap, lastTradeCodeIncomeSearchMap)

	hkSortedTable := dealer.SortAllMarketTable(hkMergedTwoDaysTable)

	dateMD := dealer.TransferMonthDay(&assignDate)

	filename := fmt.Sprint(assignDate, ".xlsx")

	generateXLSX(hkSortedTable, &filename, dateMD)
	insertPNG(&filename)

	msg := fmt.Sprintf("\nsratch success!\nsaved as '%v'\n", filename)

	enterClose(&msg)
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
	sheet.SetColWidth(0, 0, width.RankWidth+0.7)
	sheet.SetColWidth(1, 3, width.TodayWidth+0.7)
	sheet.SetColWidth(4, 4, width.LastTradeDayWidth+0.7)
	sheet.SetColWidth(5, 5, width.SeparateWidth+0.7)
	sheet.SetColWidth(6, 6, width.RankWidth+0.7)
	sheet.SetColWidth(7, 9, width.TodayWidth+0.7)
	sheet.SetColWidth(10, 10, width.LastTradeDayWidth+0.7)

	//设置行高
	height := config.Height
	sheet.Cell(0, 0).Row.SetHeight(height.TitleHeight)
	sheet.Cell(1, 0).Row.SetHeight(height.HeaderHeight)
	for i := 2; i < sheet.MaxRow; i++ {
		sheet.Rows[i].SetHeight(height.ContentHeight)
	}

	//保存excel文件
	xlsxfile.Save(*filename)
}

func insertPNG(fileName *string) {
	waterMark := config.WaterMark
	f, _ := excelize.OpenFile(*fileName)
	err := f.AddPicture(sheetName, waterMark.Position, waterMark.File,
		`{"x_offset": 10, "y_offset": 10, "print_obj": true, "lock_aspect_ratio": true, "locked": true, "positioning": "oneCell"}`)
	if err != nil {
		fmt.Println(err.Error())
	}
	f.Save()
}

func enterClose(message *string) {
	fmt.Println(*message)
	fmt.Println("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
