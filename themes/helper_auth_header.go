package themes

import (
	"fmt"
	"github.com/aubm/postmanerator/postman"
)

func helperAuthHeader(auth postman.Auth)  (string, error)  {
	switch  auth.Type {
	case "bearer":
		return fmt.Sprintf("Authorization: Bearer %v", auth.Params[0].Value), nil
	default:
		return "", fmt.Errorf("unsupported auth type %v. Please implement it", auth.Type)
	}
}