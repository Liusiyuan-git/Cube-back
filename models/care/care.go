package care

type Care struct {
	Id    int
	Care  string `orm:"index"`
	Cared string `orm:"index"`
}
