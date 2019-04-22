package engine

import (
	"context"
	"log"

	"github.com/wuxiangzhou2010/imagespider/config"
	"github.com/wuxiangzhou2010/imagespider/model"
	"github.com/wuxiangzhou2010/imagespider/util"
	"gopkg.in/olivere/elastic.v5"
)

func NewConnection() *elastic.Client {
	endpoint := config.C.GetEngineElasticUrl()
	if endpoint == "" {
		return nil
	}
	client, err := elastic.NewClient(elastic.SetURL(config.C.GetEngineElasticUrl()), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	return client
}

func (e *ConcurrentEngine) saveElasticSearch(topic model.Topic) {

	_, err := e.elastic.Index().
		Index("t66y").
		Type("topics").Id(util.Hash(topic.Url)). // hash string to get unique id
		BodyJson(topic).Do(context.Background())
	if err != nil {
		panic(err)
	}
	log.Printf("[elasticsearch] %+v\n", topic.Name)

}
