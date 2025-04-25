package sqlification

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UUIDArray struct {
	Elements []uuid.UUID
}

func NewUUIDArray(uuids []uuid.UUID) *UUIDArray {
	return &UUIDArray{uuids}
}

func (ua *UUIDArray) Scan(src interface{}) error {
	return pq.GenericArray{A: &ua.Elements}.Scan(src)
}

func (ua UUIDArray) Value() (driver.Value, error) {
	return pq.GenericArray{A: ua.Elements}.Value()
}
