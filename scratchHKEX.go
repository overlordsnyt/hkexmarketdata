package main

import (
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"hkexgo/hkex"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	assignDate := "2019-06-11"
	assignDateTop10, url := GetHKEXJson(assignDate)
	//var jsonarr []interface{} = assignDateTop10.
	if assignDateTop10 == nil {
		return
	}
	fmt.Println(assignDateTop10)

	hkSearchMap := make(map[string]hkex.HkexElement)
	hkTableSearchMap := make(map[string]map[string]hkex.Table)
	for _, v := range assignDateTop10 {
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

	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("访问港交所网站拿取%v数据时出现错误：%v", date, err.Error())
		return nil, ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jsonBytes := body[10:]

	var assignDateTop10 hkex.Hkex
	json.Unmarshal(jsonBytes, &assignDateTop10)
	return assignDateTop10, url
}

const SSEN, SZSEN = "SSE Northbound", "SZSE Northbound"
const VALID_TABLE = "top10Table"
const YIYI int64 = 100000000

func GenerateXLSX(hkTableSearchMap map[string]map[string]hkex.Table, url string) {
	//创建excel文件
	excel := xlsx.NewFile()
	sheet, _ := excel.AddSheet("沪深港通")

	{ //首行，市场标题头设置
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
		//titleStyle.Fill.BgColor="00678DEF"
		titleStyle.Fill.FgColor = "00678DEF"
		titleStyle.Fill.PatternType = "solid"
		titleStyle.Font.Name = "Heiti SC Medium"
		titleStyle.Font.Size = 14
		titleStyle.Font.Color = "FFFFFFFF"
		titleStyle.Font.Bold = true
		title1.Merge(4, 0)
		title1.SetString("沪股通")
		title1.SetStyle(titleStyle)
		title2.Merge(4, 0)
		title2.SetString("深股通")
		title2.SetStyle(titleStyle)
	}
	//设置两市分割列样式
	middleCol := sheet.Col(5)
	middleCol.Width = 2.0

	{
		//市场具体项表头创建
		tableHeaderRow := sheet.AddRow()
		ssRank := tableHeaderRow.AddCell()
		ssStockCode := tableHeaderRow.AddCell()
		ssStockName := tableHeaderRow.AddCell()
		ssPureBuy := tableHeaderRow.AddCell()
		ssYesterdayPureBuy := tableHeaderRow.AddCell()
		tableHeaderRow.AddCell()
		szRank := tableHeaderRow.AddCell()
		szStockCode := tableHeaderRow.AddCell()
		szStockName := tableHeaderRow.AddCell()
		szPureBuy := tableHeaderRow.AddCell()
		szYesterdayPureBuy := tableHeaderRow.AddCell()
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
		//市场表头样式
		tableHeaderStyle := xlsx.NewStyle()
		tableHeaderStyle.Font.Size = 12
		tableHeaderStyle.Font.Name = "Heiti SC Medium"
		tableHeaderStyle.Font.Bold = true
		tableHeaderStyle.Alignment.Horizontal = "center"
		tableHeaderStyle.Alignment.Vertical = "center"
		tableHeaderStyle.Alignment.WrapText = true
		ssRank.SetStyle(tableHeaderStyle)
		ssStockCode.SetStyle(tableHeaderStyle)
		ssStockName.SetStyle(tableHeaderStyle)
		ssPureBuy.SetStyle(tableHeaderStyle)
		ssYesterdayPureBuy.SetStyle(tableHeaderStyle)
		szRank.SetStyle(tableHeaderStyle)
		szStockCode.SetStyle(tableHeaderStyle)
		szStockName.SetStyle(tableHeaderStyle)
		szPureBuy.SetStyle(tableHeaderStyle)
		szYesterdayPureBuy.SetStyle(tableHeaderStyle)
	}

	//表内容变化所用颜色
	oddFgColor := "00DDE3F3"
	redColor := "00DB473F"
	greenColor := "#00B050"

	//获取沪股通、深股通两表
	ssTab := hkTableSearchMap[SSEN][VALID_TABLE]
	szTab := hkTableSearchMap[SZSEN][VALID_TABLE]
	//新增行、单元格，并往其中塞获取到的、处理过的数据
	for i := 0; i < 10; i++ {
		ssItem := ssTab.Tr[i].Td[0]
		szItem := szTab.Tr[i].Td[0]

		tableValueRow := sheet.AddRow()
		ssRankVal := tableValueRow.AddCell()
		ssStockCodeVal := tableValueRow.AddCell()
		ssStockNameVal := tableValueRow.AddCell()
		ssPureBuyVal := tableValueRow.AddCell()
		ssYesterdayPureBuyVal := tableValueRow.AddCell()
		tableValueRow.AddCell()
		szRankVal := tableValueRow.AddCell()
		szStockCodeVal := tableValueRow.AddCell()
		szStockNameVal := tableValueRow.AddCell()
		szPureBuyVal := tableValueRow.AddCell()
		szYesterdayPureBuyVal := tableValueRow.AddCell()

		ssRankVal.SetString(ssItem[0])
		ssStockCodeVal.SetString(ssItem[1])
		ssStockNameVal.SetString(strings.TrimRight(ssItem[2], " "))
		ssBuy := CommaStringNumberTransToBigInt(ssItem[3])
		ssSell := CommaStringNumberTransToBigInt(ssItem[4])
		//ssPureIncome:=big.Float.(new(big.Int).SetInt64(0).Sub(ssBuy,ssSell).Int64(),YIYI).
	}

	//保存excel文件
	excel.Save("golanghkexcel.xlsx")

	//ssTable:=hkTableSearchMap[SSEN][VALID_TABLE]
}

func CommaStringNumberTransToBigInt(numstr string) *big.Int {
	numstr = strings.ReplaceAll(numstr, ",", "")
	bigint, _ := new(big.Int).SetString(numstr, 10)
	return bigint
}
