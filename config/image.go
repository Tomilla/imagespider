package config

import "github.com/wuxiangzhou2010/luandun/go/spider_proj/crawler_t66y/model"

type ImageConfig struct {
	Path        string `json:"path"` // 路径
	UniqFolder  bool   `json:"isUniqFolder"`
	WorkerCount int    `json:"workerCount"`
	ImageChan   chan model.Topic
}

//type ImageConfig struct {
//	Path        string `json:"path"`// 路径
//	UniqFolder  bool`json:"isUniqFolder"`
//	WorkerCount int	`json:"workerCount"`
//}
