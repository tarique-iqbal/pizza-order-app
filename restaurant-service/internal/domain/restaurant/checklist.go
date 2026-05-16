package restaurant

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

func (c Checklist) Complete(item ChecklistItem) {
	c[item] = true
}

func (c Checklist) Reopen(item ChecklistItem) {
	c[item] = false
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
