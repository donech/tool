package middleware

import (
	"errors"
	"strings"

	"github.com/donech/tool/xjwt"
	"github.com/gin-gonic/gin"
)

var TokenNotFoundErr = errors.New("token not found err")

type JWTMiddleware struct {
	factory     xjwt.JWTFactory
	tokenLookup string
}

type Option func(middleware *JWTMiddleware)

func WithFactory(factory xjwt.JWTFactory) Option {
	return func(middleware *JWTMiddleware) {
		middleware.factory = factory
	}
}

func WithTokenLookup(tokenLookup string) Option {
	return func(middleware *JWTMiddleware) {
		middleware.tokenLookup = tokenLookup
	}
}

func NewJWTMiddleware(opts ...Option) JWTMiddleware {
	m := JWTMiddleware{}
	for _, o := range opts {
		o(&m)
	}
	return m
}

func (j JWTMiddleware) GenerateTokenHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		f := xjwt.LoginInForm{}
		err := ctx.ShouldBind(&f)
		if err != nil {
			ctx.JSON(200, gin.H{"code": 1, "msg": err.Error()})
			return
		}
		token, err := j.factory.GenerateToken(f)
		if err != nil {
			ctx.JSON(200, gin.H{"code": 1, "msg": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"code": 0, "msg": "success", "token": token})
	}
}

func (j JWTMiddleware) MiddleWareImpl() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := j.GetToken(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"msg": "get jwt-token failed: " + err.Error()})
			return
		}
		claims, err := j.factory.GetClaims(token)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"msg": "get jwt-claims failed: " + err.Error()})
			return
		}
		ctx.Set("jwt", claims)
		ctx.Next()
	}
}

func (j JWTMiddleware) GetToken(ctx *gin.Context) (token string, err error) {
	scopes := strings.Split(j.tokenLookup, ",")
	for _, scope := range scopes {
		if len(token) > 0 {
			break
		}
		parts := strings.Split(scope, ":")
		if len(parts) != 2 {
			continue
		}
		method, key := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch method {
		case "header":
			token, err = getTokenFromHeader(ctx, key)
		}
	}
	if token == "" {
		err = TokenNotFoundErr
	}
	return token, err
}

func getTokenFromHeader(ctx *gin.Context, key string) (token string, err error) {
	token = ctx.Request.Header.Get(key)
	return
}
