package database

type Handler interface {
	Write()
	Remove()
	Query()
	Load()
}
