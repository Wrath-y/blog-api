package logging

import (
	"blog-api/pkg/def"
	"context"
)

func FromContext(ctx context.Context) *logger {
	l := &logger{}
	if v := ctx.Value(def.RequestID); v != nil {
		l.requestID = v.(string)
	}
	if v := ctx.Value("v1"); v != nil {
		l.v1 = v.(string)
	}
	if v := ctx.Value("v2"); v != nil {
		l.v2 = v.(string)
	}
	if v := ctx.Value("v3"); v != nil {
		l.v3 = v.(string)
	}
	return l
}
