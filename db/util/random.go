package util

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generated a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomNDigitNumber generates a random n digit number as a string.
func RandomNDigitNumber(n int) string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Ensure the first digit is not zero
	firstDigit := rand.Intn(9) + 1
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(firstDigit))

	// Generate the remaining digits, which can include zero
	for i := 1; i < n; i++ {
		sb.WriteString(strconv.Itoa(rand.Intn(10)))
	}

	return sb.String()
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOtp(num int, str int) string {
	intpart := RandomNDigitNumber(num)
	strpart := RandomString(str)

	return fmt.Sprintf(string(intpart) + strpart)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{USD, EUR, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// Function to get Random Existing Account from the Database
func GetRandomAccountId(db *sql.DB) (int64, error) {
	var accountID int64
	// Here you type in the query
	query := ``
	err := db.QueryRow(query).Scan(&accountID)
	if err != nil {
		return 0, err
	}
	return accountID, nil
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
