package user

type RegisterUserOptions struct {
	ChatID    int64
	FirstName string
	LastName  string
	Username  string
	Language  string
	Source    *string
}
