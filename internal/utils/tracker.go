package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"runtime"

	"github.com/krobus00/storage-service/internal/constant"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	reFn = regexp.MustCompile(`([^/]+)/?$`)
)

func Trace() (string, int, string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "?", 0, "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return file, line, "?"
	}

	return file, line, fn.Name()
}

func getUserIDFromCtx(ctx context.Context) string {
	ctxUserID := ctx.Value(constant.KeyUserIDCtx)

	userID := fmt.Sprintf("%v", ctxUserID)
	if userID == "" {
		return constant.GuestID
	}
	return userID
}

func SanitizePayload(structa interface{}, censoredField map[string]bool) (res []byte, err error) {
	value := reflect.ValueOf(structa).Elem()
	typeSt := reflect.TypeOf(structa)
	size := value.NumField()

	res = append(res, '{')
	for i := 0; i < size; i++ {
		structValue := value.Field(i)
		fieldName := typeSt.Elem().Field(i).Tag.Get("json")
		// skip if field doesn't have json tag
		if fieldName == "" {
			continue
		}

		fieldValue := "REDACTED"
		if !censoredField[fieldName] {
			value, err := json.Marshal(structValue.Interface())
			if err != nil {
				return []byte{}, err
			}
			fieldValue = string(value)
		}

		res = append(res, '"')
		res = append(res, []byte(fieldName)...)
		res = append(res, '"')
		res = append(res, ':')
		res = append(res, (fieldValue)...)
		if i+1 != size {
			res = append(res, ',')
		}
	}
	res = append(res, '}')
	return res, nil
}

func NewSpan(ctx context.Context, fn string) (context.Context, trace.Span) {
	tr := otel.Tracer("")
	fn = reFn.FindString(fn)
	ctx, span := tr.Start(ctx, fn)
	span.SetAttributes(
		attribute.String("user_id", getUserIDFromCtx(ctx)),
	)
	return ctx, span
}

func SetSpanBody(span trace.Span, raw any) {
	body, err := SanitizePayload(raw, constant.RedactedField)
	if err != nil {
		return
	}
	span.SetAttributes(
		attribute.String("body", string(body)),
	)
}
