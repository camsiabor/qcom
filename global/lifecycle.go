package global

type Module interface {
	Terminate(g *G) error
}
