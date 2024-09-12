//nolint:gochecknoglobals,forcetypeassert
package log

import (
	"context"
	"fmt"
	"strings"
)

type ContextKey struct{}

type Field struct {
	Name  string
	Value any
}

var (
	ContextKeyVal = ContextKey{}
)

type Fields []Field

func (f Fields) String() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString(" [ ")
	for i, field := range f {
		strBuilder.WriteString(field.Name)
		strBuilder.WriteString(": ")
		strBuilder.WriteString(fmt.Sprintf("%v", field.Value))
		if i != len(f)-1 {
			strBuilder.WriteString(", ")
		}
	}

	strBuilder.WriteString(" ]\n")
	return strBuilder.String()
}

func GetFields(ctx context.Context) Fields {
	if ctx.Value(ContextKeyVal) == nil {
		return []Field{}
	}
	return ctx.Value(ContextKeyVal).(Fields)
}

func AddField(ctx context.Context, value Field) context.Context {
	ctx = context.WithValue(ctx, ContextKeyVal, append(GetFields(ctx), value))
	return ctx
}

func AddKeyVal(ctx context.Context, key string, val any) context.Context {
	return AddField(ctx, Field{Name: key, Value: val})
}
