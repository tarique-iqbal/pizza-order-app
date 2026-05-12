package restaurant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"restaurant-service/internal/domain/restaurant"
)

func TestNewChecklist(t *testing.T) {
	checklist := restaurant.NewChecklist()

	assert.False(t, checklist[restaurant.ChecklistBasic])
	assert.False(t, checklist[restaurant.ChecklistContract])
	assert.False(t, checklist[restaurant.ChecklistAddress])
	assert.False(t, checklist[restaurant.ChecklistDelivery])
	assert.False(t, checklist[restaurant.ChecklistPayment])
}

func TestChecklist_Complete(t *testing.T) {
	checklist := restaurant.NewChecklist()

	err := checklist.Complete(restaurant.ChecklistBasic)

	assert.NoError(t, err)
	assert.True(t, checklist[restaurant.ChecklistBasic])
}

func TestChecklist_Complete_InvalidItem(t *testing.T) {
	checklist := restaurant.NewChecklist()

	err := checklist.Complete(restaurant.ChecklistItem("invalid"))

	assert.Error(t, err)
	assert.Equal(t, "invalid checklist item: invalid", err.Error())
}

func TestChecklist_Reopen(t *testing.T) {
	checklist := restaurant.NewChecklist()

	err := checklist.Complete(restaurant.ChecklistBasic)
	assert.NoError(t, err)

	assert.True(t, checklist[restaurant.ChecklistBasic])

	err = checklist.Reopen(restaurant.ChecklistBasic)

	assert.NoError(t, err)
	assert.False(t, checklist[restaurant.ChecklistBasic])
}

func TestChecklist_Reopen_InvalidItem(t *testing.T) {
	checklist := restaurant.NewChecklist()

	err := checklist.Reopen(restaurant.ChecklistItem("invalid"))

	assert.Error(t, err)
	assert.Equal(t, "invalid checklist item: invalid", err.Error())
}

func TestChecklist_IsCompleted_ReturnsFalse(t *testing.T) {
	checklist := restaurant.NewChecklist()

	err := checklist.Complete(restaurant.ChecklistBasic)
	assert.NoError(t, err)

	assert.False(t, checklist.IsCompleted())
}

func TestChecklist_IsCompleted_ReturnsTrue(t *testing.T) {
	checklist := restaurant.NewChecklist()

	items := []restaurant.ChecklistItem{
		restaurant.ChecklistBasic,
		restaurant.ChecklistContract,
		restaurant.ChecklistAddress,
		restaurant.ChecklistDelivery,
		restaurant.ChecklistPayment,
	}

	for _, item := range items {
		err := checklist.Complete(item)
		assert.NoError(t, err)
	}

	assert.True(t, checklist.IsCompleted())
}
