package core

type Logic interface {
	Process(Model) (Model, error)
}
