package service

import "errors"

var (
	//Order errors
	ErrOrderAlreadyExists             = errors.New("order already exists")
	ErrOrderDoesNotExistToAnotherUser = errors.New("order does not exist to another user")
	ErrInvalidOrder                   = errors.New("invalid order")
	ErrOrderNotFound                  = errors.New("order not found")
	//

	//Balance errors
	ErrBalanceInvalidOrderNumber         = errors.New("invalid order number")
	ErrBalanceInvalidNegativeAmount      = errors.New("negative amount")
	ErrBalanceInvalidInsufficientBalance = errors.New("insufficient balance")
	//

	//User errors
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserInvalidPassword = errors.New("invalid password")
	ErrUserFailedRegister  = errors.New("failed to register user")
	//
)
