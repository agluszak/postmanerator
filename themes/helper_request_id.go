package themes

import (
	"github.com/aubm/postmanerator/postman"
)

func helperRequestId(request postman.Request) string {
	return request.Method + helperSlugify(request.Name)
}
