package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GenerateUUID(pre string) string {
	return strings.ToLower(fmt.Sprintf("%s_%s", pre, uuid.New().String()))
}
