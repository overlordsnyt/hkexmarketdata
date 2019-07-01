package excel

import "github.com/tealeg/xlsx"

var titleStyle, headerStyle *xlsx.Style

func init() {
	titleStyle = xlsx.NewStyle()
	titleStyle.Alignment.Horizontal = "center"
	titleStyle.Alignment.Vertical = "center"
	titleStyle.Fill.FgColor = "00678DEF"
	titleStyle.Fill.PatternType = "solid"
	titleStyle.Font.Name = "Heiti SC Medium"
	titleStyle.Font.Size = 14
	titleStyle.Font.Color = "FFFFFFFF"
	titleStyle.Font.Bold = true

	headerStyle := xlsx.NewStyle()
	headerStyle.Font.Size = 12
	headerStyle.Font.Name = "Heiti SC Medium"
	headerStyle.Font.Bold = true
	headerStyle.Alignment.Horizontal = "center"
	headerStyle.Alignment.Vertical = "center"
	headerStyle.Alignment.WrapText = true
}

func GenerateTitle(sheet *xlsx.Sheet, row, col, horizonMergeNum *int, title *string) {
	cell := sheet.Cell(*row, *col)
	cell.Merge(*horizonMergeNum, 0)
	cell.SetStyle(titleStyle)
	cell.SetString(*title)
}

func ChangeColWidth(sheet *xlsx.Sheet, colNum *int, width *float64) {
	col := sheet.Col(*colNum)
	col.Width = *width
}

func GenerateHeader(sheet *xlsx.Sheet, rowNum, colNum *int, headers *[]string) {
	col := *colNum
	for _, str := range *headers {
		cell := sheet.Cell(*rowNum, col)
		cell.SetString(str)
		cell.SetStyle(headerStyle)
		col++
	}
}
