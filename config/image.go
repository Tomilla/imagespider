package config

import (
	"github.com/wuxiangzhou2010/daily_learning/go/spider_proj/crawler_t66y/model"
	"os"
)

type ImageConfig struct {
	ImageChan   chan model.Topic
	Path        string
	UniqFolder  bool
	WorkerCount int
}

func NewImageConfig() *ImageConfig {
	//current working directory as the default image download folder
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return &ImageConfig{
		UniqFolder:  true,
		Path:        wd + "fruit",
		ImageChan:   make(chan model.Topic),
		WorkerCount: 20,
	}
}
