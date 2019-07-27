package dealer

import (
	"encoding/json"
	"fmt"
	"hkexgo/calculator"
	_type "hkexgo/type"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

const VALID_TABLE = "top10Table"

func GetHKEXJson(assignDate string, waitGroup *sync.WaitGroup) (_type.Hkex, string) {
	waitGroup.Add(1)
	defer waitGroup.Done()
	date, _ := time.Parse("2006-01-02", assignDate)
	formatDate := date.Format("20060102")
	url := fmt.Sprintf("https://sc.hkex.com.hk/TuniS/www.hkex.com.hk/chi/csm/DailyStat/data_tab_daily_%sc.js", formatDate)

	resp, err := http.Get(url)
	if err != nil {
		err = fmt.Errorf("访问港交所网站拿取%v数据时出现错误：%v", assignDate, err.Error())
		fmt.Println(err)
		return nil, ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jsonBytes := body[10:]

	var assignDateTop10 _type.Hkex
	json.Unmarshal(jsonBytes, &assignDateTop10)
	return assignDateTop10, url
}

func ScrachAssignDateTop10JsonToMarketNameSearchMap(assignDateTop10 *_type.Hkex) *map[string]map[string]_type.Table {
	hkTableSearchMap := make(map[string]map[string]_type.Table)
	for _, v := range *assignDateTop10 {
		tableMap := make(map[string]_type.Table)

		for _, vv := range v.Content {
			tableMap[vv.Table.Classname] = vv.Table
		}

		hkTableSearchMap[v.Market] = tableMap
	}
	return &hkTableSearchMap
}

func GenerateStockCodePureIncomeSearchMapFromLastTradeDateJson(lastTradeTop10 *_type.Hkex) *map[string]map[string]float64 {
	marketStockCodePureIncome := make(map[string]map[string]float64)
	for _, market := range *lastTradeTop10 {
		stockCodePureIncome := make(map[string]float64)

		for _, oneStockInfo := range market.Content[1].Table.Tr {
			stockInfoArr := oneStockInfo.Td[0]
			pureIncome := *calculator.CalculatePureIncomeDevideYi(&stockInfoArr[3], &stockInfoArr[4])
			stockCodePureIncome[stockInfoArr[1]] = pureIncome
		}
		marketStockCodePureIncome[market.Market] = stockCodePureIncome
	}
	return &marketStockCodePureIncome
}

func MergeLastTradeIncomeToAssignDateTable(adtb *map[string]map[string]_type.Table, ltcim *map[string]map[string]float64) *map[string]*_type.StockTable {
	hkStrArrTable := make(map[string][][]string)
	for k, v := range *adtb {
		tr := v[VALID_TABLE].Tr
		for _, trv := range tr {
			trstrarr := &trv.Td[0]
			hkStrArrTable[k] = append(hkStrArrTable[k], *trstrarr)
		}
	}
	hkStockTable := make(map[string]*_type.StockTable)
	for k, v := range hkStrArrTable {
		hkStockTable[k] = _type.NewStockTable(&v)
		for i, stk := range *hkStockTable[k] {
			code := v[i][1]
			codeMatchedLTDI := (*ltcim)[k][code]
			if &codeMatchedLTDI != nil {
				stk.SetLastTradeDayIncome(&codeMatchedLTDI)
				hkStockTable[k].SetLastTradeDayIncome(&i, &codeMatchedLTDI)
			}
		}
	}
	return &hkStockTable
}

func SortAllMarketTable(rawTable *map[string]*_type.StockTable) *map[string]*_type.StockTable {
	for _, v := range *rawTable {
		sort.Sort(sort.Reverse(v))
	}
	return rawTable
}
