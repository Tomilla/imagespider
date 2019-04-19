package config

import "github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/model"

type Config struct {
	OneFolder bool
	Image     ImageConfig
}

func NewConfig() *Config {
	return &Config{}
}

var AllConf = NewConfig()

func (c *Config) GetImageChan() chan model.Topic {
	return AllConf.Image.ImageChan

}
