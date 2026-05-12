package restaurant

import "fmt"

type ChecklistItem string

const (
	ChecklistBasic    ChecklistItem = "basic"
	ChecklistContract ChecklistItem = "contract"
	ChecklistAddress  ChecklistItem = "address"
	ChecklistDelivery ChecklistItem = "delivery"
	ChecklistPayment  ChecklistItem = "payment"
)

type Checklist map[ChecklistItem]bool

func NewChecklist() Checklist {
	return Checklist{
		ChecklistBasic:    false,
		ChecklistContract: false,
		ChecklistAddress:  false,
		ChecklistDelivery: false,
		ChecklistPayment:  false,
	}
}

func (c Checklist) Complete(item ChecklistItem) error {
	if !item.isValid() {
		return fmt.Errorf("invalid checklist item: %s", item)
	}

	c[item] = true

	return nil
}

func (c Checklist) Reopen(item ChecklistItem) error {
	if !item.isValid() {
		return fmt.Errorf("invalid checklist item: %s", item)
	}

	c[item] = false

	return nil
}

func (c Checklist) IsCompleted() bool {
	required := []ChecklistItem{
		ChecklistBasic,
		ChecklistContract,
		ChecklistAddress,
		ChecklistDelivery,
		ChecklistPayment,
	}

	for _, item := range required {
		if !c[item] {
			return false
		}
	}

	return true
}

func (i ChecklistItem) isValid() bool {
	switch i {
	case ChecklistBasic,
		ChecklistContract,
		ChecklistAddress,
		ChecklistDelivery,
		ChecklistPayment:
		return true
	}

	return false
}
