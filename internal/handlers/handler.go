package handlers

type (
	Handler interface{}
)

func NewHandler() Handler {
	return &handler{}
}

type handler struct{}
