package _type

type Hkex []HkexElement

type HkexElement struct {
	ID         int64     `json:"id"`
	Date       string    `json:"date"`
	Market     string    `json:"market"`
	TradingDay int64     `json:"tradingDay"`
	Content    []Content `json:"content"`
}

type Content struct {
	Style int64 `json:"style"`
	Table Table `json:"table"`
}

type Table struct {
	Classname string     `json:"classname"`
	Schema    [][]string `json:"schema"`
	Tr        []Tr       `json:"tr"`
}

type Tr struct {
	Td [][]string `json:"td"`
}
