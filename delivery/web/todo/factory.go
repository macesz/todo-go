package todo

// TodoHandlers groups HTTP handler functions.
// Like a Java controller class or JS route handler object.
type TodoHandlers struct {
	Service TodoService
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(service TodoService) *TodoHandlers {
	return &TodoHandlers{
		Service: service,
	}
}
