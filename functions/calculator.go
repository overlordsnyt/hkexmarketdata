package functions

import (
	"math/big"
	"strings"
)

func CommaStringNumberTransToBigInt(numstr *string) *big.Float {
	*numstr = strings.ReplaceAll(*numstr, ",", "")
	bigfloat, _ := new(big.Float).SetString(*numstr)
	return bigfloat
}

//初始化一个一亿的bigfloat待除
var YIYI, _ = new(big.Float).SetString("100000000")
var ZERO = new(big.Float)

func CalculatePureIncomeDevideYi(buystr, sellstr *string) *float64 {
	buy := CommaStringNumberTransToBigInt(buystr)
	sell := CommaStringNumberTransToBigInt(sellstr)
	buy = buy.Sub(buy, sell)
	buy = buy.Quo(buy, YIYI)
	f64, _ := buy.Float64()
	return &f64
}

func IsNegative(buy *float64) bool {
	return *buy < 0
}
