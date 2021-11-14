package profile

import "Cube-back/elasticsearch"

func UserEsSearch(keyWord, page string) (int, interface{}) {
	return elasticsearch.Client.Search("user", keyWord, page, []string{"name"})
}
