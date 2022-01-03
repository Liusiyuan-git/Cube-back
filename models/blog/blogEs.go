package blog

import "Cube-back/elasticsearch"

func blogEsSearch(keyWord, page string) (int, interface{}) {
	return elasticsearch.Client.Search("blog", keyWord, page, []string{"name", "title", "text", "label_type"})
}
