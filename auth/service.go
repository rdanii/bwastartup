package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	// GenerateToken digunakan untuk membuat token JWT dengan user ID yang diberikan.
	GenerateToken(userID int) (string, error)

	// ValidateToken digunakan untuk memvalidasi token JWT yang diberikan.
	ValidateToken(encodeToken string) (*jwt.Token, error)
}

type jwtService struct{}

var SECRET_KEY = []byte("shbvsefekkasakbcdbee")

// NewService membuat instance baru dari jwtService.
func NewService() *jwtService {
	return &jwtService{}
}

// GenerateToken membuat token JWT dengan menggunakan user ID yang diberikan.
func (s *jwtService) GenerateToken(userID int) (string, error) {
	// Membuat claim untuk token JWT dengan user ID.
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	// Membuat token JWT dengan claim yang sudah disiapkan.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	// Menandatangani token dengan menggunakan SECRET_KEY.
	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil
}

// ValidateToken memvalidasi token JWT yang diberikan menggunakan SECRET_KEY.
func (s *jwtService) ValidateToken(encodeToken string) (*jwt.Token, error) {
	// Parse token JWT yang diberikan.
	token, err := jwt.Parse(encodeToken, func(token *jwt.Token) (interface{}, error) {
		// Memeriksa apakah metode penandatanganan token sesuai dengan metode HMAC.
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		// Jika metode penandatanganan tidak sesuai, kembalikan error.
		if !ok {
			return nil, errors.New("invalid token")
		}

		// Mengembalikan SECRET_KEY sebagai kunci untuk memverifikasi token.
		return SECRET_KEY, nil
	})

	// Mengembalikan token dan error, jika ada.
	if err != nil {
		return token, err
	}

	return token, nil
}
