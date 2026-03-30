package validator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps go-playground/validator with custom validation rules
type CustomValidator struct {
	validate *validator.Validate
}

// NewCustomValidator creates a new validator instance with custom rules
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Register custom validations
	v.RegisterValidation("username", validateUsername)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("email_strict", validateEmailStrict)
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("age_range", validateAgeRange)
	v.RegisterValidation("date_format", validateDateFormat)

	return &CustomValidator{validate: v}
}

// Validate performs validation on a struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

// ValidateVar validates a single variable
func (cv *CustomValidator) ValidateVar(i interface{}, tag string) error {
	return cv.validate.Var(i, tag)
}

// Custom validation functions

// validateUsername checks username: 3-20 chars, alphanumeric and underscores only
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}

// validatePassword checks password strength:
// - at least 8 characters
// - at least one uppercase letter
// - at least one lowercase letter
// - at least one digit
// - at least one special character
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}

// validateEmailStrict checks email format with strict rules
func validateEmailStrict(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	
	// Basic email validation with regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return false
	}
	
	// Check for consecutive dots
	if strings.Contains(email, "..") {
		return false
	}
	
	// Check local part length (max 64 chars)
	parts := strings.Split(email, "@")
	if len(parts[0]) > 64 {
		return false
	}
	
	return true
}

// validatePhone validates phone numbers (international format)
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	
	// Remove spaces and hyphens
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	
	// Check international format: +XX... or 00XX...
	matched, _ := regexp.MatchString(`^(\+[0-9]{1,3}|00[0-9]{1,3})[0-9]{6,14}$`, phone)
	return matched
}

// validateAgeRange checks if age is between min and max (inclusive)
func validateAgeRange(fl validator.FieldLevel) bool {
	age := fl.Field().Int()
	
	// Get the parameter string
	param := fl.Param()
	
	// Parse min and max from param (format: "18 120" with space)
	// We'll use space as separator to avoid conflict with comma
	parts := strings.Fields(param)
	if len(parts) != 2 {
		return false
	}
	
	min, err1 := strconv.ParseInt(parts[0], 10, 64)
	max, err2 := strconv.ParseInt(parts[1], 10, 64)
	
	if err1 != nil || err2 != nil {
		return false
	}
	
	return age >= min && age <= max
}

// validateDateFormat checks date in YYYY-MM-DD format
func validateDateFormat(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	
	// Check format with regex
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateStr)
	if !matched {
		return false
	}
	
	// Parse and validate actual date
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// User struct with validation tags
// For age_range, use space as separator to avoid issues with comma
type User struct {
	Username    string `json:"username" validate:"required,username"`
	Email       string `json:"email" validate:"required,email_strict"`
	Password    string `json:"password" validate:"required,password"`
	Phone       string `json:"phone" validate:"required,phone"`
	Age         int    `json:"age" validate:"required,age_range=18 120"`
	BirthDate   string `json:"birth_date" validate:"required,date_format"`
	FullName    string `json:"full_name" validate:"required,min=2,max=100"`
	Nickname    string `json:"nickname" validate:"omitempty,min=2,max=30"`
	CountryCode string `json:"country_code" validate:"required,len=2,alpha"`
	PostalCode  string `json:"postal_code" validate:"required,min=3,max=10"`
}

// ValidationResult represents validation result with errors
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors map[string]string `json:"errors,omitempty"`
}

// ValidateUser validates a user and returns structured result
func ValidateUser(user *User) *ValidationResult {
	cv := NewCustomValidator()
	err := cv.Validate(user)
	
	if err == nil {
		return &ValidationResult{Valid: true}
	}
	
	// Extract validation errors
	errors := make(map[string]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()
			
			// Format param for display (replace space with -)
			displayParam := strings.ReplaceAll(param, " ", "-")
			
			var message string
			switch tag {
			case "required":
				message = "This field is required"
			case "username":
				message = "Username must be 3-20 characters and contain only letters, numbers, and underscores"
			case "password":
				message = "Password must be at least 8 characters and contain uppercase, lowercase, digit, and special character"
			case "email_strict":
				message = "Invalid email format"
			case "phone":
				message = "Invalid phone number format (use international format: +1234567890)"
			case "age_range":
				message = "Age must be between " + displayParam + " years"
			case "date_format":
				message = "Date must be in YYYY-MM-DD format"
			case "min":
				message = "Minimum length is " + param + " characters"
			case "max":
				message = "Maximum length is " + param + " characters"
			case "len":
				message = "Must be exactly " + param + " characters"
			case "alpha":
				message = "Must contain only letters"
			default:
				message = "Validation failed for " + tag
			}
			
			errors[field] = message
		}
	}
	
	return &ValidationResult{
		Valid:  false,
		Errors: errors,
	}
}

// ValidateEmail standalone email validation
func ValidateEmail(email string) error {
	cv := NewCustomValidator()
	return cv.ValidateVar(email, "email_strict")
}

// ValidatePassword standalone password validation
func ValidatePassword(password string) error {
	cv := NewCustomValidator()
	return cv.ValidateVar(password, "password")
}

// ValidateUsername standalone username validation
func ValidateUsername(username string) error {
	cv := NewCustomValidator()
	return cv.ValidateVar(username, "username")
}

// ValidatePhone standalone phone validation
func ValidatePhone(phone string) error {
	cv := NewCustomValidator()
	return cv.ValidateVar(phone, "phone")
}

// ValidateAge validates age with custom range
func ValidateAge(age int, min, max int) error {
	cv := NewCustomValidator()
	tag := fmt.Sprintf("age_range=%d %d", min, max)
	return cv.ValidateVar(age, tag)
}