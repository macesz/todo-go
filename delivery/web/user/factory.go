package user

// UserHandlers groups HTTP handler functions.
// Like a Java controller class or JS route handler object.
type UserHandlers struct {
	Service UserService
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(service UserService) *UserHandlers {
	return &UserHandlers{
		Service: service,
	}
}
