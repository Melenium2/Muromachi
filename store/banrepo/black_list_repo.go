package banrepo

type BlackList interface {
	AddBlock()
	CheckBlock()
}
