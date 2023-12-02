package uid

import (
	"github.com/google/uuid"
)

type UIDManager struct {
	sf *Sonyflake
}

var uidManager = NewManager()

func GetManager() *UIDManager {
	return uidManager
}

func NewManager() *UIDManager {
	return &UIDManager{
		sf: NewSonyflake(),
	}
}

func (m *UIDManager) ID() string {
	return m.sf.ID()
}

func (m *UIDManager) UUID() string {
	return uuid.New().String()
}
