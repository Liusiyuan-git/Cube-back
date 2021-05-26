package picture

import "io/ioutil"

func GetPicture(path string) ([]byte, error) {
	picture, err := ioutil.ReadFile(path)
	return picture, err
}
