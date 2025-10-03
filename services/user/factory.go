package user

type UserService struct {
	UserStore UserStore // Implementation here
}

func NewUserService(userStore UserStore) *UserService {
	return &UserService{
		UserStore: userStore,
	}
}
