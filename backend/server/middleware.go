package server

import (
	"avito/auth"
	"avito/auth/jwt"
	"avito/domain"
	"avito/log"
	"avito/utils"
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
)

type Middleware struct {
	logger log.Logger
}

type BodyKey struct{}
type AuthKey struct{}

func NewMiddleware(logger log.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (md *Middleware) Wrap(endpointFunc any) httprouter.Handle {
	return md.HandleResponse(md.HandleLogging(md.HandleCalling(endpointFunc)))
}

func (md *Middleware) WrapAuthed(endpointFunc any) httprouter.Handle {
	return md.HandleResponse(md.HandleLogging(md.HandleAuth(md.HandleCalling(endpointFunc))))
}

func (md *Middleware) HandleResponse(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		handler(w, r, params)
		err := domain.GetError(r)
		if err != nil {
			w.WriteHeader(int(err.Status))
			w.Write([]byte(err.String()))
			return
		}
		w.WriteHeader(int(domain.SuccessCode))
		body := r.Context().Value(BodyKey{})
		if body == nil {
			return
		}

		if err := utils.EncodeJson(w, body); err != nil {
			md.logger.Error(r.Context(), "could not encode body")
		}
	}
}

func (md *Middleware) HandleLogging(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		handler(w, r, params)
		err := domain.GetError(r)
		if err == nil {
			return
		}
		warnMsg := fmt.Sprintf("\n%s \nMethod: %s URL: %s ", err.String(), r.Method, r.URL)
		md.logger.Warn(r.Context(), warnMsg)
	}
}

func (md *Middleware) HandleAuth(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		claims, err := jwt.ParseToken(r.Header.Get("Authorization"))
		if err != nil {
			domErr := domain.NewHTTPError(err, "could not authenticate: "+err.Error(), domain.UnauthorizedCode)
			domain.SetError(r, domErr)
			return
		}

		claimsCtx := context.WithValue(r.Context(), AuthKey{}, *claims)
		*r = *r.WithContext(claimsCtx)

		handler(w, r, params)
	}
}

func (md *Middleware) HandleCalling(endpointFunc any) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var (
			argsIn  []reflect.Value
			argsOut []reflect.Value

			funcType  = reflect.TypeOf(endpointFunc)
			funcValue = reflect.ValueOf(endpointFunc)
		)

		if funcType.Kind() != reflect.Func {
			panic("endpoint must be a function")
		}
		if funcType.NumIn() > 3 {
			panic("endpoint can not have more than two parameters")
		}

		argsIn = append(argsIn, reflect.ValueOf(r.Context()))

		if funcType.NumIn() > 2 {
			// has second param
			bodyType := funcType.In(1)
			secondParam, ctx, err := utils.DecodeJson(r.Context(), r.Body, bodyType)
			*r = *r.WithContext(ctx)
			if err != nil {
				domErr := domain.NewHTTPError(err, "could not parse: "+err.Error(), domain.BadRequestCode)
				domain.SetError(r, domErr)
				return
			}

			ctx, err = utils.Validate(r.Context(), r.Body)
			*r = *r.WithContext(ctx)
			if err != nil {
				domErr := domain.NewHTTPError(err, "body validation: "+err.Error(), domain.BadRequestCode)
				domain.SetError(r, domErr)
				return
			}
			argsIn = append(argsIn, reflect.ValueOf(secondParam).Elem())
		}

		requestData := domain.RequestData{Params: params}
		if claims := r.Context().Value(AuthKey{}); claims != nil {
			requestData.Claims = claims.(auth.Claims)
		}
		argsIn = append(argsIn, reflect.ValueOf(requestData))

		argsOut = funcValue.Call(argsIn)

		outErr := argsOut[len(argsOut)-1]
		outVal := argsOut[0]

		if !outErr.IsNil() && outErr.Kind() == reflect.Pointer {
			outErr = outErr.Elem()
			httpError := outErr.Interface().(domain.HTTPError)
			domain.SetError(r, httpError)
		}

		if len(argsOut) <= 1 {
			return
		}

		bodyCtx := context.WithValue(r.Context(), BodyKey{}, outVal.Interface())
		*r = *r.WithContext(bodyCtx)
	}
}
