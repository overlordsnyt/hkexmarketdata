package configuration

type Configuration struct {
	*Width     `comment:"排名列 | 今日股票代码、名称、净买入 | 上一交易日净买入 宽度"`
	*Height    `comment:"标题行高 | 表头行高 | 内容行高"`
	*WaterMark `comment:"水印设置"`
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

type WaterMark struct {
	Position string `comment:"水印单元格位置：excel竖横单元格坐标"`
	File     string `comment:"水印图片文件名，只支持png图片"`
}

func NewConfiguration() *Configuration {
	defaultCfg := new(Configuration)
	defaultCfg.Width = new(Width)
	defaultCfg.Height = new(Height)
	defaultCfg.WaterMark = new(WaterMark)
	defaultCfg.Width.RankWidth = 5.4
	defaultCfg.Width.TodayWidth = 11
	defaultCfg.Width.LastTradeDayWidth = 20
	defaultCfg.Width.SeparateWidth = 1.4
	defaultCfg.Height.TitleHeight = 25
	defaultCfg.Height.HeaderHeight = 32
	defaultCfg.Height.ContentHeight = 32
	defaultCfg.WaterMark.Position = "J8"
	defaultCfg.WaterMark.File = "watermark.png"
	return defaultCfg
}
