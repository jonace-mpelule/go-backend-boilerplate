package health

type Module struct {
	handler *Handler
}

func NewModule() *Module {
	return &Module{
		handler: NewHandler(),
	}
}
