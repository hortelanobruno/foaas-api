package validator

import "fmt"

type MessageValidatorImpl struct{}

func NewMessageValidatorImpl() *MessageValidatorImpl {
	return &MessageValidatorImpl{}
}

func (*MessageValidatorImpl) ValidateMessage(userID string) error {
	if userID == "" {
		return fmt.Errorf("userID can't be empty")
	}

	return nil
}
