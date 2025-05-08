package validator_test

import (
	"api-service/internal/domain/restaurant"
	"api-service/internal/infrastructure/persistence"
	iValidator "api-service/internal/infrastructure/validator"
	"api-service/tests/internal/infrastructure/db"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var customValidator *iValidator.CustomValidator

func TestMain(m *testing.M) {
	testDB = db.SetupTestDB()

	userRepo := persistence.NewUserRepository(testDB)
	restaurantRepo := persistence.NewRestaurantRepository(testDB)
	customValidator = iValidator.NewCustomValidator(userRepo, restaurantRepo)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestUniqueRestaurantSlug_Valid(t *testing.T) {
	field := &mockFieldLevel{value: "test-restaurant"}
	status := customValidator.UniqueRestaurantSlug(field)

	assert.True(t, status)
}

func TestUniqueRestaurantSlug_Invalid(t *testing.T) {
	existingRestaurant := restaurant.Restaurant{Slug: "existing-restaurant"}
	testDB.Create(&existingRestaurant)

	field := &mockFieldLevel{value: "existing-restaurant"}
	status := customValidator.UniqueRestaurantSlug(field)

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
