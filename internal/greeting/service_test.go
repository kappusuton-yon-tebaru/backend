package greeting_test

import (
	"testing"

	"github.com/kappusuton-yon-tebaru/backend/internal/greeting"
	"github.com/stretchr/testify/assert"
)

type mockService struct{}

func newMock() greeting.Service {
	return &mockService{}
}

func (svc *mockService) GetGreeting() string {
	return "Hello"
}

func TestGetGreeting(t *testing.T) {
	svc := newMock()

	assert.Equal(t, svc.GetGreeting(), "Hello", "Greeting message not matched")
}
