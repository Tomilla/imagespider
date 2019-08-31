package elasticsearch

import (
    "context"
    "log"

    elastic "github.com/elastic/go-elasticsearch/v7"

    "github.com/Tomilla/imagespider/common"
    "github.com/Tomilla/imagespider/common/model"
    "github.com/Tomilla/imagespider/util"
)

type ela struct {
    c         *elastic.Client
    topicChan chan model.Topic
}

func NewConnection() *elastic.Client {
    endpoint := common.C.GetEngineElasticUrl()
    if endpoint == "" {
        return nil
    }
    cfg := elastic.Config{
        Addresses: []string{
            common.C.GetEngineElasticUrl(),
        },
    }
    client, err := elastic.NewClient(cfg)
    if err != nil {
        panic(err)
    }
    return client
}

func New(topicChan chan model.Topic) *ela {
    client := NewConnection()
    return &ela{c: client, topicChan: topicChan}
}

func (e *ela) saveElasticSearch() {
    for {
        topic, ok := <-e.topicChan
        if !ok {
            panic("saveElasticSearch not ok")
        }
        _, err := e.c.Index().
            Index("t66y").
            Type("topics").Id(util.Hash(topic.Url)). // hash string to get unique id
            BodyJson(topic).Do(context.Background())
        if err != nil {
            panic(err)
        }
        log.Printf("[saveElasticSearch] %+v %+v\n", topic.Name, topic.Url)
    }
}

func init() {
    ch := make(chan model.Topic)
    common.C.SetElasticChan(ch)
    e := New(ch)
    go e.saveElasticSearch()
}
