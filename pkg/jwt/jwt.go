package jwt

import (
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/ecdsafile"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

// PrivateKey is an ECDSA private key which was generated with the following
// command:
//
//	openssl ecparam -name prime256v1 -genkey -noout -out ecprivatekey.pem
//
// We are using a hard coded key here in this example, but in real applications,
// you would never do this. Your JWT signing key must never be in your application,
// only the public key.
const PrivateKey = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEINqRJn6HsK+Ad1Bggsdkfvpr4rg02qN0zH8x3fWtn8bXoAoGCCqGSM49
AwEHoUQDQgAEEoOvsqqI0Vxww1tybdOvNHWv3NRpqDWkzCIL+jNh1pfDVAuFi3zE
/40DtTldYDTEcBx0AWaEFbwcpYjz7A2HrQ==
-----END EC PRIVATE KEY-----`
const KeyID = "backend-sawit-test-id"

// SignToken takes a JWT and signs it with our private key, returning a JWS.
func SignToken(t jwt.Token) ([]byte, error) {
	privKey, err := ecdsafile.LoadEcdsaPrivateKey([]byte(PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("loading PEM private key: %w", err)
	}

	hdr := jws.NewHeaders()
	if err := hdr.Set(jws.AlgorithmKey, jwa.ES256); err != nil {
		return nil, fmt.Errorf("setting algorithm: %w", err)
	}
	if err := hdr.Set(jws.TypeKey, "JWT"); err != nil {
		return nil, fmt.Errorf("setting type: %w", err)
	}
	if err := hdr.Set(jws.KeyIDKey, KeyID); err != nil {
		return nil, fmt.Errorf("setting Key ID: %w", err)
	}
	return jwt.Sign(t, jwa.ES256, privKey, jwt.WithHeaders(hdr))
}

// CreateJWSWithClaims is a helper function to create JWT's with the specified
// claims.
func CreateJWSWithClaims(user map[string]interface{}) ([]byte, error) {
	t := jwt.New()
	err := t.Set(jwt.IssuerKey, "backed-sawit-pro-issuer")
	if err != nil {
		return nil, fmt.Errorf("setting issuer: %w", err)
	}
	err = t.Set(jwt.AudienceKey, "backed-sawit-pro-audience")
	if err != nil {
		return nil, fmt.Errorf("setting audience: %w", err)
	}
	err = t.Set("user", user)
	if err != nil {
		return nil, fmt.Errorf("setting permissions: %w", err)
	}
	return SignToken(t)
}
