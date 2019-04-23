package elasticsearch

import (
	"context"
	"log"

	"github.com/wuxiangzhou2010/imagespider/config"
	"github.com/wuxiangzhou2010/imagespider/model"
	"github.com/wuxiangzhou2010/imagespider/util"
)

type ela struct {
	c         *elastic.Client
	topicChan chan model.Topic
}

func New(topicChan chan model.Topic) *ela {
	endpoint := config.C.GetEngineElasticUrl()
	if endpoint == "" {
		return nil
	}
	client, err := elastic.NewClient(elastic.SetURL(config.C.GetEngineElasticUrl()), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
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
	config.C.SetElasticChan(ch)
	e := New(ch)
	go e.saveElasticSearch()
}
