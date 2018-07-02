package elastic

import (
	"github.com/olivere/elastic"
	"time"
	"context"
)

type Eser interface {
	Init() error
	Put(area string,message []byte) error
	Index() error
}

type es struct {
	client *elastic.Client
	addr string
	ctx context.Context
	index string
}

//addr http://192.168.2.10:9201
func New(addr,index string) Eser{
	return &es{
		addr:addr,
		index:index,
		ctx:context.Background(),
	}
}


func (e *es) Init() error {
	client, err := elastic.NewClient(
		elastic.SetURL(e.addr),
		elastic.SetSniff(false),    //TODO 代表什么意思
		elastic.SetHealthcheckInterval(10*time.Second))
		//elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		//elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))

	if err != nil {
		return err
		// Handle error
	}

	e.client = client
	return nil
}



/*
// Add a document to the index @ liuliqiang.info
 tweet := Tweet{User: "olivere", Message: "Take Five"}
 _, err = client.Index().
     Index("twitter").
     Type("tweet").
     Id("1").
     BodyJson(tweet).
     Do()
 if err != nil {
     // Handle error
     panic(err)
 }

*/


func (e *es) Put(area string,message []byte) error {
	_, err := e.client.Index().
		Index(area).
		Type("doc").
		BodyJson(string(message)).
		Do(e.ctx)

	return err
}



// Index 索引不存在就创建
func (e *es) Index() error {
	exists, err := e.client.IndexExists(e.index).Do(e.ctx)
	if err != nil {
		return err
	}

	if !exists {
		// Create a new index.
		createIndex, err := e.client.CreateIndex(e.index).Do(e.ctx)
		if err != nil {
			return err
			// Handle error
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	return nil
}