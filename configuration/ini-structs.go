package configuration

type Configuration struct {
	Width Width `comment:"排名列 | 今日股票代码、名称、净买入 | 上一交易日净买入 宽度"`
}

type Width struct {
	RankWidth         float64
	TodayWidth        float64
	LastTradeDayWidth float64
}

func NewConfiguration() *Configuration {
	defaultCfg := new(Configuration)
	width := new(Width)
	width.RankWidth = 6.1
	width.TodayWidth = 11.7
	width.LastTradeDayWidth = 20.7
	defaultCfg.Width = *width
	return defaultCfg
}
