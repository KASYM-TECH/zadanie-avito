package utils

import (
	"avito/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"reflect"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func DecodeJson(ctx context.Context, reader io.Reader, objType reflect.Type) (any, context.Context, error) {
	v := reflect.New(objType).Interface()
	if err := json.NewDecoder(reader).Decode(v); err != nil {
		return nil, nil, fmt.Errorf("decode body: %w", err)
	}

	return v, ctx, nil
}

func EncodeJson(w io.Writer, val any) error {
	if err := json.NewEncoder(w).Encode(val); err != nil {
		return fmt.Errorf("encode body: %w", err)
	}
	return nil
}

func Validate(ctx context.Context, obj any) (context.Context, error) {
	if err := validate.Struct(obj); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			ctx = log.AddField(ctx, log.Field{Name: "Field", Value: err.StructField()})
			ctx = log.AddField(ctx, log.Field{Name: "ActualTag", Value: err.ActualTag()})
			ctx = log.AddField(ctx, log.Field{Name: "Tag", Value: err.Tag()})
			ctx = log.AddField(ctx, log.Field{Name: "Param", Value: err.Param()})
			ctx = log.AddField(ctx, log.Field{Name: "Value", Value: err.Value()})
		}
		return ctx, fmt.Errorf("invalid of type: %T", obj)
	}
	return ctx, nil
}
