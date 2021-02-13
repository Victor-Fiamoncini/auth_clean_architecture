package authusecase_test

import (
	"testing"

	"github.com/Victor-Fiamoncini/auth_clean_architecture/src/domain/entities"
	usecases "github.com/Victor-Fiamoncini/auth_clean_architecture/src/domain/usecases/auth_usecase"
	luber "github.com/Victor-Fiamoncini/auth_clean_architecture/src/infra/repositories/load_user_by_email_repository"
	luber_mocks "github.com/Victor-Fiamoncini/auth_clean_architecture/src/infra/repositories/load_user_by_email_repository/mocks"
	uatr "github.com/Victor-Fiamoncini/auth_clean_architecture/src/infra/repositories/update_access_token_repository"
	uatr_mocks "github.com/Victor-Fiamoncini/auth_clean_architecture/src/infra/repositories/update_access_token_repository/mocks"
	encrypter "github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/helpers/encrypter"
	encrypter_mocks "github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/helpers/encrypter/mocks"
	token_generator "github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/helpers/token_generator"
	token_generator_mocks "github.com/Victor-Fiamoncini/auth_clean_architecture/src/shared/helpers/token_generator/mocks"
	"github.com/stretchr/testify/assert"
)

func makeSut() (usecases.IAuthUseCase, luber.ILoadUserByEmailRepository, encrypter.IEncrypter, token_generator.ITokenGenerator, uatr.IUpdateAccessTokenRepository) {
	loadUserByEmailRepositorySpy := luber_mocks.NewLoadUserByEmailRepositorySpy()
	encrypterSpy := encrypter_mocks.NewEncrypterSpy()
	tokenGeneratorSpy := token_generator_mocks.NewTokenGeneratorSpy()
	updateAccessTokenRepositorySpy := uatr_mocks.NewUpdateAccessTokenRepositorySpy()

	encrypterSpy.SetIsValid(true)
	tokenGeneratorSpy.SetAccessToken("any_token")

	user := entities.NewUser()
	user.SetID("any_id")
	user.SetPassword("hashed_password")

	loadUserByEmailRepositorySpy.SetUser(user)

	authUseCase := usecases.NewAuthUseCase(loadUserByEmailRepositorySpy, encrypterSpy, tokenGeneratorSpy, updateAccessTokenRepositorySpy)

	return authUseCase, loadUserByEmailRepositorySpy, encrypterSpy, tokenGeneratorSpy, updateAccessTokenRepositorySpy
}

func TestShouldCallLoadUserByEmailRepositoryWithCorrectEmail(t *testing.T) {
	sut, loadUserByEmailRepositorySpy, _, _, _ := makeSut()

	sut.SetEmail("any_email@mail.com")
	sut.SetPassword("any_password")

	_, err := sut.Auth()

	assert.Equal(t, "any_email@mail.com", loadUserByEmailRepositorySpy.GetEmail())
	assert.Nil(t, err)
}

func TestShouldCallEncrypterWithCorrectValues(t *testing.T) {
	sut, loadUserByEmailRepositorySpy, encrypterSpy, _, _ := makeSut()

	sut.SetEmail("valid_email@mail.com")
	sut.SetPassword("any_password")

	_, err := sut.Auth()

	assert.Equal(t, "any_password", encrypterSpy.GetPassword())
	assert.Equal(t, loadUserByEmailRepositorySpy.GetUser().GetPassword(), encrypterSpy.GetHashedPassword())
	assert.Nil(t, err)
}

func TestShouldCallTokenGeneratorWithCorrectUserID(t *testing.T) {
	sut, loadUserByEmailRepositorySpy, _, tokenGeneratorSpy, _ := makeSut()

	sut.SetEmail("valid_email@mail.com")
	sut.SetPassword("valid_password")

	_, err := sut.Auth()

	assert.Equal(t, loadUserByEmailRepositorySpy.GetUser().GetID(), tokenGeneratorSpy.GetUserID())
	assert.Nil(t, err)
}

func TestShouldReturnAnAccessTokenIfCorrectCredentialsAreProvided(t *testing.T) {
	sut, _, _, tokenGeneratorSpy, _ := makeSut()

	sut.SetEmail("valid_email@mail.com")
	sut.SetPassword("valid_password")

	accessToken, err := sut.Auth()

	assert.Equal(t, tokenGeneratorSpy.GetAccessToken(), accessToken)
	assert.NotNil(t, accessToken)
	assert.NotEmpty(t, accessToken)
	assert.Nil(t, err)
}

func TestShouldCallUpdateAccessTokenRepositoryWithCorrectValues(t *testing.T) {
	sut, loadUserByEmailRepositorySpy, _, tokenGeneratorSpy, updateAccessTokenRepositorySpy := makeSut()

	sut.SetEmail("valid_email@mail.com")
	sut.SetPassword("valid_password")

	_, err := sut.Auth()

	assert.Equal(t, updateAccessTokenRepositorySpy.GetUserID(), loadUserByEmailRepositorySpy.GetUser().GetID())
	assert.Equal(t, updateAccessTokenRepositorySpy.GetAccessToken(), tokenGeneratorSpy.GetAccessToken())
	assert.Nil(t, err)
}
