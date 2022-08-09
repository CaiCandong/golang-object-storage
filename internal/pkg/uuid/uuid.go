package uuid

import (
	"github.com/google/uuid"
)

func GenUUid() string {
	id := uuid.New()
	//2c5fd02e-fe21-4aa9-a1f7-73af1081829e
	return id.String()
}
