package validator

import (
	"testing"
)

// Table-driven tests with subtests for User validation
func TestValidateUser(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		expectValid bool
		expectError string // field name that should have error
	}{
		{
			name: "Valid user",
			user: &User{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "Pass123!@#",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				Nickname:    "johnny",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: true,
		},
		{
			name: "Invalid username - too short",
			user: &User{
				Username:    "jo",
				Email:       "john@example.com",
				Password:    "Pass123!@#",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Username",
		},
		{
			name: "Invalid username - special chars",
			user: &User{
				Username:    "john@doe",
				Email:       "john@example.com",
				Password:    "Pass123!@#",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Username",
		},
		{
			name: "Invalid email - missing @",
			user: &User{
				Username:    "john_doe",
				Email:       "johnexample.com",
				Password:    "Pass123!@#",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Email",
		},
		{
			name: "Invalid email - double dots",
			user: &User{
				Username:    "john_doe",
				Email:       "john..doe@example.com",
				Password:    "Pass123!@#",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Email",
		},
		{
			name: "Invalid password - too short",
			user: &User{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "Pass1!",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Password",
		},
		{
			name: "Invalid password - no uppercase",
			user: &User{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "pass123!@#",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Password",
		},
		{
			name: "Invalid password - no digit",
			user: &User{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "Password!!!",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Password",
		},
		{
			name: "Invalid age - too young",
			user: &User{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "Pass123!@#",
				Phone:       "+12345678901",
				Age:         15,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Age",
		},
		{
			name: "Invalid age - too old",
			user: &User{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "Pass123!@#",
				Phone:       "+12345678901",
				Age:         150,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Age",
		},
		{
			name: "Invalid birth date format",
			user: &User{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "Pass123!@#",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "15-01-1999",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "BirthDate",
		},
		{
			name: "Invalid phone - no country code",
			user: &User{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "Pass123!@#",
				Phone:       "12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "US",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "Phone",
		},
		{
			name: "Invalid country code - too long",
			user: &User{
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "Pass123!@#",
				Phone:       "+12345678901",
				Age:         25,
				BirthDate:   "1999-01-15",
				FullName:    "John Doe",
				CountryCode: "USA",
				PostalCode:  "12345",
			},
			expectValid: false,
			expectError: "CountryCode",
		},
		{
			name: "Missing required fields",
			user: &User{
				Username:    "john_doe",
				Email:       "",
				Password:    "",
				Phone:       "",
				Age:         0,
				BirthDate:   "",
				FullName:    "",
				CountryCode: "",
				PostalCode:  "",
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateUser(tt.user)
			
			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got %v", tt.expectValid, result.Valid)
			}
			
			if !tt.expectValid && tt.expectError != "" {
				if _, exists := result.Errors[tt.expectError]; !exists {
					t.Errorf("Expected error for field %s, but got errors: %v", tt.expectError, result.Errors)
				}
			}
			
			// Print errors for debugging in case of unexpected validation
			if !tt.expectValid && result.Errors != nil {
				t.Logf("Validation errors: %v", result.Errors)
			}
		})
	}
}

// Table-driven tests for standalone validators
func TestStandaloneValidators(t *testing.T) {
	tests := []struct {
		name      string
		validator func() error
		expectErr bool
	}{
		{
			name: "Valid username",
			validator: func() error {
				return ValidateUsername("john_doe_123")
			},
			expectErr: false,
		},
		{
			name: "Invalid username - too short",
			validator: func() error {
				return ValidateUsername("jo")
			},
			expectErr: true,
		},
		{
			name: "Invalid username - special chars",
			validator: func() error {
				return ValidateUsername("john@doe")
			},
			expectErr: true,
		},
		{
			name: "Valid email",
			validator: func() error {
				return ValidateEmail("test@example.com")
			},
			expectErr: false,
		},
		{
			name: "Invalid email - no @",
			validator: func() error {
				return ValidateEmail("testexample.com")
			},
			expectErr: true,
		},
		{
			name: "Invalid email - no domain",
			validator: func() error {
				return ValidateEmail("test@")
			},
			expectErr: true,
		},
		{
			name: "Valid password",
			validator: func() error {
				return ValidatePassword("StrongP@ss123")
			},
			expectErr: false,
		},
		{
			name: "Invalid password - too short",
			validator: func() error {
				return ValidatePassword("Weak1!")
			},
			expectErr: true,
		},
		{
			name: "Invalid password - no uppercase",
			validator: func() error {
				return ValidatePassword("weakpass123!")
			},
			expectErr: true,
		},
		{
			name: "Valid phone",
			validator: func() error {
				return ValidatePhone("+12345678901")
			},
			expectErr: false,
		},
		{
			name: "Valid phone with spaces",
			validator: func() error {
				return ValidatePhone("+1 234 567 8901")
			},
			expectErr: false,
		},
		{
			name: "Invalid phone - no plus",
			validator: func() error {
				return ValidatePhone("12345678901")
			},
			expectErr: true,
		},
		{
			name: "Valid age",
			validator: func() error {
				return ValidateAge(25, 18, 120)
			},
			expectErr: false,
		},
		{
			name: "Invalid age - too young",
			validator: func() error {
				return ValidateAge(16, 18, 120)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.validator()
			if (err != nil) != tt.expectErr {
				t.Errorf("Expected error=%v, got %v (error: %v)", tt.expectErr, err != nil, err)
			}
		})
	}
}

// Benchmark tests
func BenchmarkValidateUser(b *testing.B) {
	user := &User{
		Username:    "john_doe",
		Email:       "john@example.com",
		Password:    "Pass123!@#",
		Phone:       "+12345678901",
		Age:         25,
		BirthDate:   "1999-01-15",
		FullName:    "John Doe",
		Nickname:    "johnny",
		CountryCode: "US",
		PostalCode:  "12345",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidateUser(user)
	}
}

func BenchmarkValidateStandalone(b *testing.B) {
	b.Run("Username", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ValidateUsername("john_doe_123")
		}
	})
	
	b.Run("Email", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ValidateEmail("test@example.com")
		}
	})
	
	b.Run("Password", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ValidatePassword("StrongP@ss123")
		}
	})
	
	b.Run("Phone", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ValidatePhone("+12345678901")
		}
	})
}