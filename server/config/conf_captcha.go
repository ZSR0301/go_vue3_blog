package config

type Captcha struct {
	Height   int     `json:"height" yaml:"height"`
	Width    int     `json:"width" yaml:"width"`
	Length   int     `json:"length" yaml:"lenght"`
	MaxSkew  float64 `json:"max_skew" yaml:"max_skew"`
	DotCount int     `json:"dot_count" yaml:"dot_count"`
}
