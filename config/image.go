package config

import "github.com/Tomilla/imagespider/model"

type ImageConfig struct {
    Path        string `json:"path"` // 路径
    UniqFolder  bool   `json:"isUniqFolder"`
    WorkerCount int    `json:"workerCount"`
    ImageChan   chan model.Topic
    HungryChan  chan bool
}

// type ImageConfig struct {
// 	Path        string `json:"path"`// 路径
// 	UniqFolder  bool`json:"isUniqFolder"`
// 	WorkerCount int	`json:"workerCount"`
// }
