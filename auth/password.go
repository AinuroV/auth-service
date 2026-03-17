package auth

import "golang.org/x/crypto/bcrypt"

// Функция для хэширования пароля
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// Функция проверки пароля пользователя
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Функция проверка сложности пароля
func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasNumber := false
	hasUpper := false
	hasSpecial := false

	for _, ch := range password {
		switch {
		case ch >= '0' && ch <= '9':
			hasNumber = true
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= '!' && ch <= '/', ch >= ':' && ch <= '@':
			hasSpecial = true
		}
	}

	return hasNumber && hasUpper && hasSpecial
}
