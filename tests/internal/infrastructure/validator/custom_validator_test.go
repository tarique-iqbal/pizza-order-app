package validation_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	"pizza-order-api/internal/domain/user"
	"pizza-order-api/internal/infrastructure/persistence"
	infrastructureValidator "pizza-order-api/internal/infrastructure/validator"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var userRepo user.UserRepository
var customValidator *infrastructureValidator.CustomValidator

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect test database:", err)
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatal("Failed to migrate test database:", err)
	}

	testDB = db
	userRepo = persistence.NewUserRepository(testDB)
	customValidator = infrastructureValidator.NewCustomValidator(userRepo)

	code := m.Run()
	os.Exit(code)
}

func TestUniqueEmail_Valid(t *testing.T) {
	field := &mockFieldLevel{value: "newuser@example.com"}
	status := customValidator.UniqueEmail(field)

	assert.True(t, status)
}

func TestUniqueEmail_Invalid(t *testing.T) {
	existingUser := user.User{Email: "existing@example.com"}
	testDB.Create(&existingUser)

	field := &mockFieldLevel{value: "existing@example.com"}
	status := customValidator.UniqueEmail(field)

	assert.False(t, status)
}

type mockFieldLevel struct {
	value string
}

func (m *mockFieldLevel) Top() reflect.Value      { return reflect.Value{} }
func (m *mockFieldLevel) Parent() reflect.Value   { return reflect.Value{} }
func (m *mockFieldLevel) Field() reflect.Value    { return reflect.ValueOf(m.value) }
func (m *mockFieldLevel) StructFieldName() string { return "" }
func (m *mockFieldLevel) Param() string           { return "" }
func (m *mockFieldLevel) GetTag() string          { return "" }
func (m *mockFieldLevel) FieldName() string       { return "email" }
func (m *mockFieldLevel) GetStructFieldOK() (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}
func (m *mockFieldLevel) GetStructFieldOK2() (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}
func (m *mockFieldLevel) ExtractType(field reflect.Value) (reflect.Value, reflect.Kind, bool) {
	return field, field.Kind(), false
}
func (m *mockFieldLevel) GetStructFieldOKAdvanced(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool) {
	return reflect.Value{}, reflect.Invalid, false
}
func (m *mockFieldLevel) GetStructFieldOKAdvanced2(val reflect.Value, namespace string) (reflect.Value, reflect.Kind, bool, bool) {
	return reflect.Value{}, reflect.Invalid, false, false
}
