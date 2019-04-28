package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

var (
	errTokenHelperSetup      = errors.New("token helper: impossible to setup the token helper twice")
	errTokenHelperTypeAssert = errors.New("token helper: impossible to type assert token string")
)

var Token TokenHelperInterface

type TokenHelper struct {
	secret                 string
	issuer                 string
	storeTokenDuration     int
	anonymousTokenDuration int
	customerTokenDuration  int
}

type TokenHelperInterface interface {
	CreateToken(audience, storeName string, storeID int) *JWTToken
	EncodeToken(jwtToken *JWTToken) (string, error)
	DecodeToken(encodedToken string) (*JWTToken, error)
	AppendEncodedTokenToHeader(c echo.Context, encodedJWTToken string)
}

// JWTToken is a custom claim used to generate a token for the session, customer and cart
type JWTToken struct {
	jwt.StandardClaims
	CustomerID    uint   `json:"customerId,omitempty"`
	CustomerName  string `json:"customerName,omitempty"`
	CustomerEmail string `json:"customerEmail,omitempty"`
}

// TODO: change JWTTOken uint's to int

func SetupTokenHelper(secret, issuer string, storeTokenDuration, anonymousTokenDuration, customerTokenDuration int) error {
	if Token != nil {
		return errTokenHelperSetup
	}
	Token = &TokenHelper{secret, issuer, storeTokenDuration, anonymousTokenDuration, customerTokenDuration}
	return nil
}

func (s *TokenHelper) CreateToken(audience, storeName string, storeID int) *JWTToken {
	now := time.Now()
	jwtToken := &JWTToken{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewV4().String(),
			IssuedAt:  now.Unix(),
			Issuer:    s.issuer,
			Audience:  audience,
			ExpiresAt: now.Add(time.Duration(s.storeTokenDuration) * time.Hour * 24).Unix(),
		},
	}

	return jwtToken
}

// refreshToken is used to refresh the JWT Token.
// Since all JWT tokens are immutable, a new ID is created and also the expiration date is updated
func (s *TokenHelper) refreshToken(jwtToken *JWTToken, duration int) {
	now := time.Now()
	jwtToken.Id = uuid.NewV4().String()
	jwtToken.IssuedAt = now.Unix()
	jwtToken.ExpiresAt = now.Add(time.Duration(duration) * time.Hour * 24).Unix()
}

// AppendCustomerToToken creates a new token appending the customer info
func (s *TokenHelper) AppendCustomerToToken(jwtToken JWTToken, id uint, name, email string) *JWTToken {
	s.refreshToken(&jwtToken, s.customerTokenDuration)
	jwtToken.Subject = strconv.FormatUint(uint64(id), 10)
	jwtToken.CustomerID = id
	jwtToken.CustomerName = name
	jwtToken.CustomerEmail = email
	return &jwtToken
}

func (s *TokenHelper) EncodeToken(jwtToken *JWTToken) (string, error) {
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtToken)
	encodedToken, err := tokenClaim.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("token helper: error when encoding token. %s", err)
	}
	return encodedToken, nil
}

func (s *TokenHelper) DecodeToken(encodedToken string) (*JWTToken, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &JWTToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token helper: token decode error. %s", err)
	}
	// Type assert token to internal token
	jwtToken, ok := token.Claims.(*JWTToken)
	if ok == false {
		return nil, errTokenHelperTypeAssert
	}
	return jwtToken, nil
}

func (s *TokenHelper) AppendEncodedTokenToHeader(c echo.Context, encodedJWTToken string) {
	c.Response().Header().Add("Api-Token", encodedJWTToken)
}
