package fixtures

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"restaurant-service/internal/domain/restaurant"
	"restaurant-service/tests/testutil"
)

func LoadRestaurantFixtures(t *testing.T, db *gorm.DB) error {
	checklist := restaurant.NewChecklist()
	checklist.Complete(restaurant.ChecklistBasic)

	restaurants := []restaurant.Restaurant{
		{
			Name:      "Pizza Paradise",
			VATNumber: "DE987687654",
			Checklist: checklist,
			CreatedAt: time.Now().UTC(),
		},
		{
			Name:         "Anatolische Küche",
			VATNumber:    "DE987321321",
			Slug:         testutil.StringPtr("anatolische-kueche"),
			Email:        testutil.StringPtr("kontakt@anatolisch.de"),
			Phone:        testutil.StringPtr("+49 40 76543210"),
			DeliveryType: restaurant.DeliveryOwn,
			DeliveryKm:   testutil.Int16Ptr(7),
			Specialties:  datatypes.JSON([]byte(`["italian"]`)),
			Checklist: restaurant.Checklist{
				restaurant.ChecklistBasic:    true,
				restaurant.ChecklistContract: true,
				restaurant.ChecklistAddress:  true,
				restaurant.ChecklistDelivery: true,
				restaurant.ChecklistPayment:  true,
			},
			CreatedAt: time.Now().UTC(),
		},
	}

	for _, r := range restaurants {
		r.ID = testutil.MustNewID()
		r.OwnerID = testutil.MustNewID()

		err := db.Create(&r).Error
		require.NoError(t, err)
	}

	return nil
}
