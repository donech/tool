package xjwt

import (
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrNoPrivateKeyFile = errors.New("no private key err")
)

type LoginInForm struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginFunc func(LoginInForm) (jwt.MapClaims, error)
type GenerateTokenFunc func(loginFunc LoginFunc) (jwt.Token, error)

type JWTFactory struct {
	singingAlgorithm string
	key              []byte
	publicKeyFile    string
	privateKeyFile   string
	timeout          time.Duration
	loginFunc        LoginFunc

	signKey        interface{}
	decodeKey      interface{}
	privateSignKey *rsa.PrivateKey
	publicSignKey  *rsa.PublicKey
}
type Option func(factory *JWTFactory)

func WithKey(key []byte) Option {
	return func(factory *JWTFactory) {
		factory.key = key
	}
}
func WithPublicKeyFile(key string) Option {
	return func(factory *JWTFactory) {
		factory.publicKeyFile = key
	}
}
func WithPrivateKeyFile(key string) Option {
	return func(factory *JWTFactory) {
		factory.privateKeyFile = key
	}
}
func WithTimeout(duration time.Duration) Option {
	return func(factory *JWTFactory) {
		factory.timeout = duration
	}
}

func WithLoginFunc(f LoginFunc) Option {
	return func(factory *JWTFactory) {
		factory.loginFunc = f
	}
}

func NewJWTFactory(config Config, opts ...Option) (JWTFactory, error) {
	d, err := time.ParseDuration(config.Timeout)
	if err != nil {
		d = time.Minute * 10
	}
	jm := JWTFactory{
		singingAlgorithm: config.SingingAlgorithm,
		key:              []byte(config.Key),
		publicKeyFile:    config.PublicKeyFile,
		privateKeyFile:   config.PrivateKeyFile,
		timeout:          d,
	}
	for _, v := range opts {
		v(&jm)
	}
	err = jm.Init()
	return jm, err
}

func (f *JWTFactory) Init() (err error) {
	f.signKey = f.key
	f.decodeKey = f.key
	if f.useRsaAlgorithm() {
		if f.privateKeyFile != "" {
			err = f.privateKey()
		}
		if f.publicKeyFile != "" {
			err = f.publicKey()
		}
		f.signKey = f.privateSignKey
		f.decodeKey = f.publicSignKey
	}
	return err
}

func (f *JWTFactory) privateKey() error {
	keyData, err := ioutil.ReadFile(f.privateKeyFile)
	if err != nil {
		return err
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		return err
	}
	f.privateSignKey = key
	return nil
}

func (f *JWTFactory) publicKey() error {
	keyData, err := ioutil.ReadFile(f.publicKeyFile)
	if err != nil {
		return err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return err
	}
	f.publicSignKey = key
	return nil
}

func (f JWTFactory) GenerateToken(form LoginInForm) (string, error) {
	claims, err := f.loginFunc(form)
	claims["exp"] = time.Now().Add(f.timeout).Unix()
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod(f.singingAlgorithm), claims)
	s, err := token.SignedString(f.signKey)
	if err != nil {
		return "", err
	}
	return s, err
}

func (f JWTFactory) VerifyToken(token string) bool {
	t, _ := f.parseToken(token)
	return t.Valid
}

func (f JWTFactory) RefreshToken(token string) (string, error) {
	claims, err := f.GetClaims(token)
	if err != nil {
		return "", err
	}
	claims["exp"] = time.Now().Add(f.timeout).Unix()
	t := jwt.NewWithClaims(jwt.GetSigningMethod(f.singingAlgorithm), claims)
	s, err := t.SignedString(f.signKey)
	if err != nil {
		return "", err
	}
	return s, err
}

func (f JWTFactory) GetClaims(token string) (jwt.MapClaims, error) {
	t, err := f.parseToken(token)
	if err != nil {
		return jwt.MapClaims{}, err
	}
	return t.Claims.(jwt.MapClaims), nil
}

func (f JWTFactory) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return f.decodeKey, nil
	})
}

func (f *JWTFactory) useRsaAlgorithm() bool {
	switch f.singingAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}
