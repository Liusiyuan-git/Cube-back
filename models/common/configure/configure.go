package configure

import (
	"Cube-back/log"
	"encoding/json"
	"os"
)

func Get(conf interface{}) {
	file, err := os.Open("conf/conf.json")
	if err != nil {
		log.Error(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(conf)
	if err != nil {
		log.Error(err)
	}
	defer file.Close()
}

func init() {
}
