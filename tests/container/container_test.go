package container_test

import (
	"log"
	"os"
	"testing"

	"pizza-order-api/container"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testContainer *container.Container

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	if err := godotenv.Load("../../.env.test"); err != nil {
		log.Fatalf("Error loading .env.test file: %v", err)
	}

	var err error
	testContainer, err = container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize test container: %v", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestNewContainer_ShouldInitializeCorrectly(t *testing.T) {
	assert.NotNil(t, testContainer, "Container should not be nil")
	assert.NotNil(t, testContainer.DB, "Database connection should be initialized")
	assert.NotNil(t, testContainer.UserHandler, "UserHandler should be initialized")
}
