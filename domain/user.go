package domain

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string
}

// Custom errors for user validation, need to develop further...., its just a start
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	// Add more checks (e.g., password strength)
	return nil
}
