package validator

type MessageValidator interface {
	ValidateMessage(userID string) error
}
