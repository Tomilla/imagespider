package common

type Engine struct {
    WorkerCount int    `json:"workerCount"`
    ElasticUrl  string `json:"elasticUrl"`
}
