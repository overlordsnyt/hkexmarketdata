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

	assignDate := "2019-06-17"
	assignDateTop10, url := GetHKEXJson(assignDate)
	if assignDateTop10 == nil {
		return
	}

	lastTradeDate := "2019-06-14"
	lastTradeDateTop10, url := GetHKEXJson(lastTradeDate)
	if lastTradeDateTop10 == nil {
		return
	}

	hkTableSearchMap := ScrachAssignDateTop10JsonToMarketNameSearchMap(&assignDateTop10)

	lastTradeCodeIncomeSearchMap := *GenerateStockCodePureIncomeSearchMapFromLastTradeDateJson(&lastTradeDateTop10)

	//TODO 把上一交易日相应净买入根据股票代码跟在指定交易日的股票信息后
	//TODO 根据指定交易日的净买入从高到低排序该市场的股票

	income := lastTradeCodeIncomeSearchMap["SSE Northbound"]["600519"]
	println(income)

	GenerateXLSX(hkTableSearchMap, url)
}

func GetHKEXJson(assignDate string) (hkex.Hkex, string) {
	date, _ := time.Parse("2006-01-02", assignDate)
	formatDate := date.Format("20060102")
	url := fmt.Sprintf("https://sc.hkex.com.hk/TuniS/www.hkex.com.hk/chi/csm/DailyStat/data_tab_daily_%sc.js", formatDate)

	resp, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("访问港交所网站拿取%v数据时出现错误：%v", date, err.Error())
		fmt.Println(err)
		return nil, ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jsonBytes := body[10:]

	var assignDateTop10 hkex.Hkex
	json.Unmarshal(jsonBytes, &assignDateTop10)
	return assignDateTop10, url
}

func ScrachAssignDateTop10JsonToMarketNameSearchMap(assignDateTop10 *hkex.Hkex) *map[string]map[string]hkex.Table {
	hkTableSearchMap := make(map[string]map[string]hkex.Table)
	for _, v := range *assignDateTop10 {
		tableMap := make(map[string]hkex.Table)

		for _, vv := range v.Content {
			tableMap[vv.Table.Classname] = vv.Table
		}

		hkTableSearchMap[v.Market] = tableMap
	}
	return &hkTableSearchMap
}

func GenerateStockCodePureIncomeSearchMapFromLastTradeDateJson(lastTradeTop10 *hkex.Hkex) *map[string]map[string]float64 {
	marketStockCodePureIncome := make(map[string]map[string]float64)
	for _, market := range *lastTradeTop10 {
		stockCodePureIncome := make(map[string]float64)

		for _, oneStockInfo := range market.Content[1].Table.Tr {
			stockInfoArr := oneStockInfo.Td[0]
			pureIncome, _ := CalculatePureIncomeDevideYi(&stockInfoArr[3], &stockInfoArr[4])
			stockCodePureIncome[stockInfoArr[1]] = *pureIncome
		}
		marketStockCodePureIncome[market.Market] = stockCodePureIncome
	}
	return &marketStockCodePureIncome
}

const SSEN, SZSEN = "SSE Northbound", "SZSE Northbound"
const VALID_TABLE = "top10Table"
const NUM_FORMAT = "0.0000_ " //尾空格，坑死人

func GenerateXLSX(hkTableSearchMap *map[string]map[string]hkex.Table, url string) {
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
	greenColor := "0000B050"
	//数据单元格的基础样式
	dataBaseStyle := new(xlsx.Style)
	dataBaseStyle.Alignment.Vertical = "center"
	dataBaseStyle.Alignment.Horizontal = "center"
	dataBaseStyle.Fill.PatternType = "solid"
	dataBaseStyle.Font.Size = 11
	dataBaseStyle.Font.Name = "Heiti SC Light"
	dataBaseStyle.Font.Family = 1

	//获取沪股通、深股通两表
	ssTab := (*hkTableSearchMap)[SSEN][VALID_TABLE]
	szTab := (*hkTableSearchMap)[SZSEN][VALID_TABLE]
	//新增行、单元格，并往其中塞获取到的、处理过的数据
	for i := 0; i < 10; i++ {
		//获取市场表的单条数据
		ssItem := ssTab.Tr[i].Td[0]
		szItem := szTab.Tr[i].Td[0]

		//在新行增加相应数量的单元格，并把带填单元格命名
		tableValueRow := sheet.AddRow()
		tableValueRow.SetHeight(32)
		ssRankVal := tableValueRow.AddCell()
		ssStockCodeVal := tableValueRow.AddCell()
		ssStockNameVal := tableValueRow.AddCell()
		ssPureIncomeVal := tableValueRow.AddCell()
		ssYesterdayPureBuyVal := tableValueRow.AddCell()
		tableValueRow.AddCell()
		szRankVal := tableValueRow.AddCell()
		szStockCodeVal := tableValueRow.AddCell()
		szStockNameVal := tableValueRow.AddCell()
		szPureIncomeVal := tableValueRow.AddCell()
		szYesterdayPureBuyVal := tableValueRow.AddCell()

		//根据单双行设置单元格底色
		nowStyle := *dataBaseStyle
		if i%2 == 0 {
			nowStyle.Fill.FgColor = oddFgColor
		} else {
			nowStyle.Fill.FgColor = "00FFFFFF"
		}
		//给单元格填色
		ssRankVal.SetStyle(&nowStyle)
		ssStockCodeVal.SetStyle(&nowStyle)
		ssStockNameVal.SetStyle(&nowStyle)
		sspiStyle := nowStyle
		ssPureIncomeVal.SetStyle(&sspiStyle)
		ssldStyle := nowStyle
		ssYesterdayPureBuyVal.SetStyle(&ssldStyle)
		szRankVal.SetStyle(&nowStyle)
		szStockCodeVal.SetStyle(&nowStyle)
		szStockNameVal.SetStyle(&nowStyle)
		szpiStyle := nowStyle
		szPureIncomeVal.SetStyle(&szpiStyle)
		szldStyle := nowStyle
		szYesterdayPureBuyVal.SetStyle(&szldStyle)

		//从单条数据中取出属于特定单元格的数据，并进行计算、上色
		ssRank, _ := strconv.Atoi(ssItem[0])
		ssRankVal.SetInt(ssRank)
		ssStockCodeVal.SetString(ssItem[1])
		ssStockNameVal.SetString(strings.TrimRight(ssItem[2], "　"))
		ssPureIncome, ssNeg := CalculatePureIncomeDevideYi(&ssItem[3], &ssItem[4])
		ssPureIncomeVal.SetFloatWithFormat(*ssPureIncome, NUM_FORMAT)
		if *ssNeg {
			ssPureIncomeVal.GetStyle().Font.Color = greenColor
		} else {
			ssPureIncomeVal.GetStyle().Font.Color = redColor
		}
		ssYesterdayPureBuyVal.SetString("waiting...")
		szRank, _ := strconv.Atoi(ssItem[0])
		szRankVal.SetInt(szRank)
		szStockCodeVal.SetString(fmt.Sprintf("%06s", szItem[1]))
		szStockNameVal.SetString(strings.TrimRight(szItem[2], "　"))
		szPureIncome, szNeg := CalculatePureIncomeDevideYi(&szItem[3], &szItem[4])
		szPureIncomeVal.SetFloatWithFormat(*szPureIncome, NUM_FORMAT)
		if *szNeg {
			szPureIncomeVal.GetStyle().Font.Color = greenColor
		} else {
			szPureIncomeVal.GetStyle().Font.Color = redColor
		}
		szYesterdayPureBuyVal.SetString("waiting...")
	}

	//设置列宽
	sheet.SetColWidth(0, 0, 6.1)
	sheet.SetColWidth(1, 3, 11.7)
	sheet.SetColWidth(4, 4, 20.7)
	sheet.SetColWidth(6, 6, 6.1)
	sheet.SetColWidth(7, 9, 11.7)
	sheet.SetColWidth(10, 10, 20.7)

	//保存excel文件
	excel.Save("golanghkexcel.xlsx")
}

func CommaStringNumberTransToBigInt(numstr *string) *big.Float {
	*numstr = strings.ReplaceAll(*numstr, ",", "")
	bigfloat, _ := new(big.Float).SetString(*numstr)
	return bigfloat
}

//初始化一个一亿的bigfloat待除
var YIYI, _ = new(big.Float).SetString("100000000")
var ZERO = new(big.Float)

func CalculatePureIncomeDevideYi(buystr, sellstr *string) (*float64, *bool) {
	buy := CommaStringNumberTransToBigInt(buystr)
	sell := CommaStringNumberTransToBigInt(sellstr)
	rawIncome := buy.Sub(buy, sell)
	yiIncome := rawIncome.Quo(rawIncome, YIYI)
	f64, _ := yiIncome.Float64()
	neg := yiIncome.Cmp(ZERO) == -1
	return &f64, &neg
}
