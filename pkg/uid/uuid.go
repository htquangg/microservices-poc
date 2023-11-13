package uid

import "github.com/google/uuid"

func UUIDV4() string {
	uid := uuid.New()
	return uid.String()
}
