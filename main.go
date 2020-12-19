package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"io/ioutil"
	"log"
	"reflect"
)

const (
	booksFile = "./books.json"
	mappingFile = "./mapping.json"
)

type Book struct {
	Author string `json:"author"`
	Title string `json:"title"`
	Description string `json:"description"`
	PubDate string `json:"pub_date"`
	Category string `json:"category"`
}

func main()  {
	ctx := context.Background()
	es, err := elastic.NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	defer es.Stop()

	exists, err := es.IndexExists("books").Do(ctx)

	if !exists {
		fmt.Println("---------Create Index----------")
		byte, err := ioutil.ReadFile(mappingFile)
		if err != nil {
			log.Fatalln(err)
		}

		cInd, err := es.CreateIndex("books").BodyString(string(byte)).Do(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		if !cInd.Acknowledged {
			log.Fatalln("client not ack")
		}
	}

	byte, err := ioutil.ReadFile(booksFile)
	if err != nil {
		log.Fatalln(err)
	}

	var books []Book
	err = json.Unmarshal(byte, &books)
	if err != nil {
		log.Fatalln(err)
	}
	// 新增Book文档
	//var wg sync.WaitGroup
	//for i := 0; i < len(books); i++ {
	//	wg.Add(1)
	//	go func(index int) {
	//		defer wg.Done()
	//		pi, err := es.Index().Index("books").Id(strconv.Itoa(index + 1)).BodyJson(books[index]).Do(ctx)
	//		if err != nil {
	//			log.Fatalln(err)
	//		}
	//		fmt.Printf("Indexed tweet %s to index %s, type %s\n", pi.Id, pi.Index, pi.Type)
	//	}(i)
	//}
	//wg.Wait()

	// ES查询
	// query := elastic.NewMatchQuery()
	query := elastic.NewMatchQuery("description", "美白")
	query = query.Analyzer("ik_max_word")

	hl := elastic.NewHighlight()
	hl = hl.Field("description")
	hl = hl.PreTags("<em>").PostTags("</em>")

	result, err := es.Search().Index("books").Highlight(hl).Query(query).Sort("pub_date", true).From(0).Size(10).Pretty(true).Do(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	for _, hit := range result.Hits.Hits {
		fmt.Println(hit.Highlight["description"])
	}
	var b Book
	for _, item := range result.Each(reflect.TypeOf(b)) {
		if book, ok := item.(Book); ok {
			fmt.Printf("Book info: title [%s] author [%s] pubDate [%s]\r\n", book.Title, book.Author, book.PubDate)
		}
	}
}
