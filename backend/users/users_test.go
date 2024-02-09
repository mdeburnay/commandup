package users

import (
	"main/models"
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

func (mock *MockUserRepository) GetUser(id int) (models.User, error) {
	args := mock.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

func TestGetUser(t *testing.T) {
	// Set mock values
	// Create an instance of our mock repository
	mockUserRepo := new(MockUserRepository)

	// Setup expectations
	mockUser := models.User{ID: 12345678}
	mockUserRepo.On("GetUser", 12345678).Return(mockUser, nil)

	// Call the function we want to test
	result, err := mockUserRepo.GetUser(12345678)

	// Assert expectations
	mockUserRepo.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, mockUser, result)
}
