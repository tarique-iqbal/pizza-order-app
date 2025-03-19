package container_test

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"pizza-order-api/container"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testContainer *container.Container

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	cmd := exec.Command("go", "env", "GOMOD")
	output, _ := cmd.Output()
	moduleRoot := filepath.Dir(string(output))

	if err := godotenv.Load(moduleRoot + "/.env.test"); err != nil {
		log.Fatalf("Error loading .env.test file: %v", err)
	}

	var err error
	testContainer, err = container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize test container: %v", err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestNewContainer_ShouldInitializeCorrectly(t *testing.T) {
	assert.NotNil(t, testContainer, "Container should not be nil")
	assert.NotNil(t, testContainer.DB, "Database connection should be initialized")
	assert.NotNil(t, testContainer.UserHandler, "UserHandler should be initialized")
}
