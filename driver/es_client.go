package driver

import (
	"src/github.com/olivere/elastic"
	"time"
)

func GetEsClient() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetGzip(true),
		/*elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),*/
	)
	if err != nil {
		panic(err)
	}
	return client
}