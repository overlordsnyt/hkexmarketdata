package configuration

type Configuration struct {
	Width  Width  `comment:"排名列 | 今日股票代码、名称、净买入 | 上一交易日净买入 宽度"`
	Height Height `comment:"标题行高 | 表头行高 | 内容行高"`
}

type Width struct {
	RankWidth         float64
	TodayWidth        float64
	LastTradeDayWidth float64
	SeparateWidth     float64 `comment:"两表分隔中行宽度"`
}

type Height struct {
	TitleHeight   float64
	HeaderHeight  float64
	ContentHeight float64
}

func NewConfiguration() *Configuration {
	defaultCfg := new(Configuration)
	width := new(Width)
	width.RankWidth = 5.4
	width.TodayWidth = 11
	width.LastTradeDayWidth = 20
	width.SeparateWidth = 1.4
	defaultCfg.Width = *width
	height := new(Height)
	height.TitleHeight = 25
	height.HeaderHeight = 32
	height.ContentHeight = 32
	defaultCfg.Height = *height
	return defaultCfg
}
