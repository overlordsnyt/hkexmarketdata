package main

import (
	json2 "encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"hkexgo/hkex"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	assignDate := "2019-06-11"
	json, url := GetHKEXJson(assignDate)
	//var jsonarr []interface{} = json.
	fmt.Println(json)

	hkSearchMap := make(map[string]hkex.HkexElement)
	hkTableSearchMap := make(map[string]map[string]hkex.Table)
	for _, v := range json {
		hkSearchMap[v.Market] = v
		tableMap := make(map[string]hkex.Table)

		for _, vv := range v.Content {
			tableMap[vv.Table.Classname] = vv.Table
		}

		hkTableSearchMap[v.Market] = tableMap
	}

	/*excel:=xlsx.NewFile()
	sheet,_:=excel.AddSheet("hkex1")
	row:=sheet.AddRow()
	cell:=row.AddCell()
	cell.Value="first excel!"
	excel.Save("golanghkexcel.xlsx")*/

	GenerateXLSX(hkTableSearchMap, url)
}

func GetHKEXJson(assignDate string) (hkex.Hkex, string) {
	date, _ := time.Parse("2006-01-02", assignDate)
	formatDate := date.Format("20060102")
	url := fmt.Sprintf("https://sc.hkex.com.hk/TuniS/www.hkex.com.hk/chi/csm/DailyStat/data_tab_daily_%sc.js", formatDate)

	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jsonBytes := body[10:]

	var json hkex.Hkex
	json2.Unmarshal(jsonBytes, &json)
	return json, url
}

const SSEN, SZSEN = "SSE Northbound", "SZSE Northbound"
const VALID_TABLE = "top10Table"

func GenerateXLSX(hkTableSearchMap map[string]map[string]hkex.Table, url string) {
	//创建excel文件
	excel := xlsx.NewFile()
	sheet, _ := excel.AddSheet("沪深港通")

	//首行，市场标题头设置
	firstRow := sheet.AddRow()
	title1 := firstRow.AddCell()
	firstRow.AddCell()
	firstRow.AddCell()
	firstRow.AddCell()
	firstRow.AddCell()
	firstRow.AddCell()
	title2 := firstRow.AddCell()
	//市场标题头填值、样式
	titleStyle := xlsx.NewStyle()
	titleStyle.Alignment.Horizontal = "center"
	titleStyle.Alignment.Vertical = "center"
	title1.Merge(4, 0)
	title1.SetString("沪股通")
	title1.SetStyle(titleStyle)
	title2.Merge(4, 0)
	title2.SetString("深股通")
	title2.SetStyle(titleStyle)

	//设置两市分割列样式
	middleCol := sheet.Col(5)
	middleCol.Width = 2.0

	//市场具体项表头创建
	tableHeadRow := sheet.AddRow()
	ssRank := tableHeadRow.AddCell()
	ssStockCode := tableHeadRow.AddCell()
	ssStockName := tableHeadRow.AddCell()
	ssPureBuy := tableHeadRow.AddCell()
	ssYesterdayPureBuy := tableHeadRow.AddCell()
	tableHeadRow.AddCell()
	szRank := tableHeadRow.AddCell()
	szStockCode := tableHeadRow.AddCell()
	szStockName := tableHeadRow.AddCell()
	szPureBuy := tableHeadRow.AddCell()
	szYesterdayPureBuy := tableHeadRow.AddCell()
	//市场表头填值
	ssRank.SetString("排名")
	ssStockCode.SetString("股票代码")
	ssStockName.SetString("股票名称")
	ssPureBuy.SetString("净买入（亿元）")
	ssYesterdayPureBuy.SetString("前一交易日净买入额（亿元）")
	szRank.SetString("排名")
	szStockCode.SetString("股票代码")
	szStockName.SetString("股票名称")
	szPureBuy.SetString("净买入（亿元）")
	szYesterdayPureBuy.SetString("前一交易日净买入额（亿元）")

	//保存excel文件
	excel.Save("golanghkexcel.xlsx")

	//ssTable:=hkTableSearchMap[SSEN][VALID_TABLE]
}
