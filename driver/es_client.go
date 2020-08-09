package driver

import (
	"github.com/sunil206b/store_utils_go/logger"
	"src/github.com/olivere/elastic"
	"time"
)

func GetEsClient() *elastic.Client {
	log := logger.GetLogger()
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log),
		elastic.SetInfoLog(log),
	)
	if err != nil {
		panic(err)
	}
	return client
}