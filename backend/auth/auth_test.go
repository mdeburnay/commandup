package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type MockAuthRepo struct {
	mock.Mock
}

func (mock *MockAuthRepo) Login(email, password string) (bool, error) {
	args := mock.Called(email, password)
	return args.Bool(0), args.Error(1)
}

func TestLogin(t *testing.T) {
	email := "example@email.com"
	password := "myPassword"

	mockAuthRepo := new(MockAuthRepo)

	mockAuthRepo.On("Login", email, password).Return(true, nil)

	success, err := mockAuthRepo.Login(email, password)

	assert.NoError(t, err, "Error was not expected during login")
	assert.True(t, success, "Login should be successful with correct credentials")
	mockAuthRepo.AssertExpectations(t)
}
