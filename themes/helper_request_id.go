package themes

import (
	"github.com/aubm/postmanerator/postman"
	"strings"
)

func helperRequestId(request postman.Request) string {
	return strings.ToLower(request.Method) + "-" + helperSlugify(request.Name)
}
