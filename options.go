package gommander

type Option struct {
	name  string
	help  string
	short string
	long  string
	args  []*Argument
}
