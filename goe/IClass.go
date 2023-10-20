package goe

type IClass interface {
	Build(goe *Goe)
	Name() string
}
