package config

import "github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/model"

var StartPages = []string{
	//"http://t66y.com/thread0806.php?fid=8",  // 新时代
	//"http://t66y.com/thread0806.php?fid=16", //达盖尔
	"http://t66y.com/thread0806.php?fid=21", //下载区
}

type Config struct {
	Image ImageConfig
}

func NewConfig() *Config {
	return &Config{}
}

var AllConf = NewConfig()

func (c *Config) GetImageChan() chan model.Topic {
	return AllConf.Image.ImageChan

}
