package middleware

import (
	"InterviewBackendSawitProGolang/generated"
	"context"
	"errors"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/ecdsafile"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"strings"
)

const JWTClaimsContextKey = "jwt_claims"
const KeyID = "backend-sawit-test-id"
const PublicKey = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEEoOvsqqI0Vxww1tybdOvNHWv3NRp
qDWkzCIL+jNh1pfDVAuFi3zE/40DtTldYDTEcBx0AWaEFbwcpYjz7A2HrQ==
-----END PUBLIC KEY-----`

var (
	ErrNoAuthHeader      = errors.New("Authorization header is missing")
	ErrInvalidAuthHeader = errors.New("Authorization header is malformed")
	ErrClaimsInvalid     = errors.New("Provided claims do not match expected scopes")
)

type JWSValidator interface {
	ValidateJws(jwsString string) (jwt.Token, error)
}

type Authenticator struct {
	KeySet jwk.Set
}

func (a *Authenticator) ValidateJws(jwsString string) (jwt.Token, error) {
	return jwt.Parse([]byte(jwsString), jwt.WithKeySet(a.KeySet),
		jwt.WithAudience("backed-sawit-pro-audience"), jwt.WithIssuer("backed-sawit-pro-issuer"))
}

var _ JWSValidator = (*Authenticator)(nil)

func NewMiddleware() (echo.MiddlewareFunc, error) {
	spec, err := generated.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("loading spec: %w", err)
	}
	auth := new(Authenticator)
	if err := auth.Init(); err != nil {
		return nil, err
	}
	validator := middleware.OapiRequestValidatorWithOptions(spec,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(auth),
			},
		})
	return validator, nil
}

func (a *Authenticator) Init() error {
	set := jwk.NewSet()
	pubKey := jwk.NewECDSAPublicKey()
	loadPubKey, err := ecdsafile.LoadEcdsaPublicKey([]byte(PublicKey))
	if err != nil {
		return fmt.Errorf("loading PEM private key: %w", err)
	}
	err = pubKey.FromRaw(loadPubKey)
	if err != nil {
		return fmt.Errorf("parsing jwk key: %w", err)
	}

	err = pubKey.Set(jwk.AlgorithmKey, jwa.ES256)
	if err != nil {
		return fmt.Errorf("setting key algorithm: %w", err)
	}

	err = pubKey.Set(jwk.KeyIDKey, KeyID)
	if err != nil {
		return fmt.Errorf("setting key ID: %w", err)
	}

	set.Add(pubKey)
	a.KeySet = set
	return nil
}

func NewAuthenticator(v JWSValidator) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

func Authenticate(v JWSValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	// Our security scheme is named BearerAuth, ensure this is the case
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	// Now, we need to get the JWS from the request, to match the request expectations
	// against request contents.
	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	// if the JWS is valid, we have a JWT, which will contain a bunch of claims.
	token, err := v.ValidateJws(jws)
	if err != nil {
		return fmt.Errorf("validating JWS: %w", err)
	}

	userID, err := GetClaimsFromToken(token)
	if err != nil {
		return fmt.Errorf("validating JWS: %w", err)
	}

	eCtx := middleware.GetEchoContext(ctx)
	eCtx.Set("user_id", userID)

	return nil
}

func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}

	return strings.TrimPrefix(authHdr, prefix), nil
}

func GetClaimsFromToken(t jwt.Token) (string, error) {
	user, found := t.Get("user")
	if !found {
		return "", nil
	}

	mUser := user.(map[string]interface{})
	return mUser["id"].(string), nil
}
