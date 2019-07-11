package excel

import (
	"github.com/tealeg/xlsx"
	"hkexgo/calculator"
	_type "hkexgo/type"
)

var titleStyle, headerStyle, cellBaseStyle *xlsx.Style

//单元格变化颜色
const oddFgColor, redColor, greenColor, borderColor = "00DDE3F3", "00DB473F", "0000B050", "00BFBFBF"
const NUM_FORMAT = "0.00_ " //浮点数excel显示格式，保留小数位，尾空格

func init() {
	//全边框单元格样式
	allBorderStyle := xlsx.NewStyle()
	allBorderStyle.Border.BottomColor = borderColor
	allBorderStyle.Border.Top = "thin"
	allBorderStyle.Border.TopColor = borderColor
	allBorderStyle.Border.Left = "thin"
	allBorderStyle.Border.LeftColor = borderColor
	allBorderStyle.Border.Right = "thin"
	allBorderStyle.Border.RightColor = borderColor
	//大标题样式
	tmp := *allBorderStyle
	titleStyle = &tmp
	titleStyle.Alignment.Horizontal = "center"
	titleStyle.Alignment.Vertical = "center"
	titleStyle.Fill.FgColor = "00678DEF"
	titleStyle.Fill.PatternType = "solid"
	titleStyle.Font.Name = "Heiti SC Medium"
	titleStyle.Font.Size = 14
	titleStyle.Font.Color = "FFFFFFFF"
	titleStyle.Font.Bold = true
	//表头样式
	tmp2 := *allBorderStyle
	headerStyle = &tmp2
	headerStyle.Font.Size = 12
	headerStyle.Font.Name = "Heiti SC Medium"
	headerStyle.Font.Bold = true
	headerStyle.Alignment.Horizontal = "center"
	headerStyle.Alignment.Vertical = "center"
	headerStyle.Alignment.WrapText = true
	headerStyle.Border.Bottom = "thin"
	//数据单元格基础样式
	cellBaseStyle = xlsx.NewStyle()
	cellBaseStyle.Alignment.Vertical = "center"
	cellBaseStyle.Alignment.Horizontal = "center"
	cellBaseStyle.Fill.PatternType = "solid"
	cellBaseStyle.Font.Size = 11
	cellBaseStyle.Font.Name = "Heiti SC Light"
	cellBaseStyle.Font.Family = 1
	cellBaseStyle.Border.Left = "thin"
	cellBaseStyle.Border.LeftColor = borderColor
	cellBaseStyle.Border.Right = "thin"
	cellBaseStyle.Border.RightColor = borderColor
}

func GenerateTitle(sheet *xlsx.Sheet, row, col, horizonMergeNum int, title string) {
	cell := sheet.Cell(row, col)
	cell.Merge(horizonMergeNum, 0)
	cell.SetStyle(titleStyle)
	cell.SetString(title)
}

func ChangeColWidth(sheet *xlsx.Sheet, colNum int, width float64) {
	col := sheet.Col(colNum)
	col.Width = width
}

func GenerateHeader(sheet *xlsx.Sheet, stRow, stCol int, headers *[]string) {
	col := stCol
	for _, str := range *headers {
		cell := sheet.Cell(stRow, col)
		cell.SetString(str)
		cell.SetStyle(headerStyle)
		col++
	}
}

func FillHKEXDataLine(sheet *xlsx.Sheet, row, col int, stk *_type.Stock, nowStyle *xlsx.Style) {
	//按位置获取要填充的单元格
	rankCell := sheet.Cell(row, col)
	col++
	stockCodeCell := sheet.Cell(row, col)
	col++
	stockNameCell := sheet.Cell(row, col)
	col++
	pureBuyCell := sheet.Cell(row, col)
	col++
	lastTradePureBuyCell := sheet.Cell(row, col)
	//填充数据
	rankCell.SetInt(stk.Rank)
	stockCodeCell.SetString(stk.StockCode)
	stockNameCell.SetString(stk.StockName)
	pureBuyCell.SetFloatWithFormat(stk.TodayIncome, NUM_FORMAT)
	if stk.LastTradeDayIncome != 0 {
		lastTradePureBuyCell.SetFloatWithFormat(stk.LastTradeDayIncome, NUM_FORMAT)
	}
	rankCell.SetStyle(nowStyle)
	stockCodeCell.SetStyle(nowStyle)
	stockNameCell.SetStyle(nowStyle)
	//根据正负收益改变净买入数字颜色
	pbStyle, ltStyle := *nowStyle, *nowStyle
	pbStyle.Font.Color = greenColor
	ltStyle.Font.Color = greenColor
	if calculator.IsPositive(&stk.TodayIncome) {
		pbStyle.Font.Color = redColor
	}
	if calculator.IsPositive(&stk.LastTradeDayIncome) {
		ltStyle.Font.Color = redColor
	}
	pureBuyCell.SetStyle(&pbStyle)
	lastTradePureBuyCell.SetStyle(&ltStyle)
}

func FillHKEXData(sheet *xlsx.Sheet, stRow, stCol int, stktb *_type.StockTable) {
	for i, stk := range *stktb {
		//单双行样式
		nowStyle := *cellBaseStyle
		if i%2 == 0 {
			nowStyle.Fill.FgColor = oddFgColor
		}
		FillHKEXDataLine(sheet, stRow, stCol, &stk, &nowStyle)
		stRow++
	}
	PaintBottomBorder(sheet, stRow-1, stCol, 5)
}

func PaintBottomBorder(sheet *xlsx.Sheet, stRow, stCol, cellNum int) {
	for i := 0; i < cellNum; i++ {
		cell := sheet.Cell(stRow, stCol+i)
		cell.GetStyle().Border.Bottom = "thin"
		cell.GetStyle().Border.BottomColor = borderColor
	}
}
