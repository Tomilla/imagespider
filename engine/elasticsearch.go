package engine

import (
	"github.com/wuxiangzhou2010/imagespider/config"
	"gopkg.in/olivere/elastic.v5"
)

func NewConnection() *elastic.Client {

	client, err := elastic.NewClient(elastic.SetURL(config.C.GetEngineElasticUrl()), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	return client
}
