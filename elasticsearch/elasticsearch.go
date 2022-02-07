package elasticsearch

import (
	"Cube-back/log"
	"Cube-back/models/common/configure"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"strconv"
	"strings"
)

type EsClient struct {
	client *elasticsearch7.Client
}

var Client *EsClient

type Conf struct {
	ElasticsearchIp   string
	ElasticsearchPort string
}

func (es *EsClient) Create(index, content string, id int) {
	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: strconv.Itoa(id),
		Body:       strings.NewReader(content),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting response: %s", err))
	}
	defer res.Body.Close()
}

func (es *EsClient) Delete(index, DocumentID string) {
	req := esapi.DeleteRequest{
		Index:      index,
		DocumentID: DocumentID,
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting response: %s", err))
	}
	defer res.Body.Close()
}

func (es *EsClient) SearchAll(index string) (int, []interface{}) {
	var length int
	var maps []interface{}
	var r map[string]interface{}
	var buf bytes.Buffer
	query := map[string]interface{}{
		"size": 10000,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"index": map[string]string{"order": "desc"}},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Error(fmt.Sprintf("Error encoding query: %s", err))
	}
	res, err := es.client.Search(
		es.client.Search.WithIndex(index),
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithBody(&buf),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting response: %s", err))
		return length, maps
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error(fmt.Sprintf("Error parsing the response body: %s", err))
		return length, maps
	}
	if _, ok := r["status"]; ok && r["status"].(float64) != 200 {
		return length, maps
	}
	length = int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	maps = r["hits"].(map[string]interface{})["hits"].([]interface{})
	defer res.Body.Close()
	return length, maps
}

func (es *EsClient) Search(index, keyWord, page string, keyBox []string) (int, []interface{}) {
	var length int
	var maps []interface{}
	from, _ := strconv.Atoi(page)
	var r map[string]interface{}
	var buf bytes.Buffer
	query := map[string]interface{}{
		"from": (from - 1) * 10,
		"size": 10,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{},
			},
		},
		"highlight": map[string]interface{}{
			"pre_tags":  "<b class='key' style='color:red'>",
			"post_tags": "</b>",
			"fields":    map[string]interface{}{},
		},
		"sort": []map[string]interface{}{
			map[string]interface{}{"index": map[string]string{"order": "desc"}},
		},
	}
	for _, item := range keyBox {
		i := map[string]interface{}{
			"match": map[string]string{
				item: keyWord,
			},
		}
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"] = append(query["query"].(map[string]interface{})["bool"].(map[string]interface{})["should"].([]map[string]interface{}), i)
		query["highlight"].(map[string]interface{})["fields"].(map[string]interface{})[item] = map[string]interface{}{}
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Error(fmt.Sprintf("Error encoding query: %s", err))
	}
	res, err := es.client.Search(
		es.client.Search.WithIndex(index),
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithBody(&buf),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting response: %s", err))
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Error(fmt.Sprintf("Error parsing the response body: %s", err))
	}
	length = int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))
	maps = r["hits"].(map[string]interface{})["hits"].([]interface{})
	defer res.Body.Close()
	return length, maps
}

func (es *EsClient) Build(index string) {
	req := esapi.IndexRequest{
		Index: index,
	}
	res, err := req.Do(context.Background(), es.client)
	if err != nil {
		log.Error(fmt.Sprintf("Error getting response: %s", err))
	}
	defer res.Body.Close()
}

func init() {
	Client = new(EsClient)
	conf := new(Conf)
	configure.Get(&conf)
	cfg := elasticsearch7.Config{
		Addresses: []string{
			"http://" + conf.ElasticsearchIp + ":" + conf.ElasticsearchPort,
		},
	}
	client, err := elasticsearch7.NewClient(cfg)
	if err != nil {
		log.Error(fmt.Sprintf("Error creating the client: %s", err))
	}
	Client.client = client
}
