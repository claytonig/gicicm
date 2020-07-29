package common

// common errors
const (
	AccountAlreadyExistsError = "account already exists"
	InternalServerError       = "internal server error"
	BadRequestError           = "invalid input format"
	PasswordValidationError   = "invalid password, should have more than 8 characters, atleast 1 symbol, 1 uppercase character and a number"
	EmailValidationError      = "invalid email"
	AccountNotFoundError      = "account does not exist"
	InvalidCredentialsError   = "invalid credentials"
	UnAuthorizedError         = "not permitted to perform this operation"
)
