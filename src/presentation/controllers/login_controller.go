package controller

import (
	"reflect"

	auth_usecase "github.com/Victor-Fiamoncini/auth_clean_architecture/src/domain/usecases/auth_usecase"
	"github.com/Victor-Fiamoncini/auth_clean_architecture/src/presentation/contracts"
	"github.com/Victor-Fiamoncini/auth_clean_architecture/src/presentation/contracts/payloads"
	"github.com/Victor-Fiamoncini/auth_clean_architecture/src/presentation/http"
	shared_custom_errors "github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/errors"
	"github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/types"
	"github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/validators"
)

// LoginController struct
type LoginController struct {
	AuthUseCase    auth_usecase.IAuthUseCase
	EmailValidator validators.IEmailValidator
}

// NewLoginController func
func NewLoginController(
	authUseCase auth_usecase.IAuthUseCase,
	emailValidator validators.IEmailValidator,
) contracts.IController {
	return &LoginController{
		AuthUseCase:    authUseCase,
		EmailValidator: emailValidator,
	}
}

// Handle LoginController method
func (lc *LoginController) Handle(httpRequest contracts.IRequest) contracts.IResponse {
	httpResponse := http.NewResponse()

	if httpRequest == nil || !reflect.ValueOf(httpRequest.GetBody()).IsValid() {
		return httpResponse.ServerError()
	}

	parsedBody := httpRequest.GetBody().(*payloads.LoginPayload)

	email := parsedBody.Email
	password := parsedBody.Password

	if email == "" {
		return httpResponse.BadRequest(shared_custom_errors.NewMissingParamError("email"))
	}

	lc.EmailValidator.SetEmail(email)

	if !lc.EmailValidator.Run() {
		return httpResponse.BadRequest(shared_custom_errors.NewInvalidParamError("email"))
	}

	if password == "" {
		return httpResponse.BadRequest(shared_custom_errors.NewMissingParamError("password"))
	}

	lc.AuthUseCase.SetEmail(email)
	lc.AuthUseCase.SetPassword(password)

	accessToken, err := lc.AuthUseCase.Auth()

	if err != nil || accessToken == "" {
		return httpResponse.Unauthorized()
	}

	responseBody := make(types.Map)
	responseBody["AccessToken"] = accessToken

	return httpResponse.Success(responseBody)
}
