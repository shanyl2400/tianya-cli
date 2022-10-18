package repository

type State interface {
	Save()
	Load()
}
