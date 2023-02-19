package models

import (
  "database/sql"
  "errors"
  "github.com/go-sql-driver/mysql"
  "golang.org/x/crypto/bcrypt"
  "strings"
  "time"
)

// User type. Notice how the field names and types align
// with the columns in the database "users" table?
type User struct {
  ID             int
  Name           string
  Email          string
  HashedPassword []byte
  Created        time.Time
}

// UserModel type which wraps a database connection pool.
type UserModel struct {
  DB *sql.DB
}

// Insert method to add a new record to the "users" table.
func (m *UserModel) Insert(name, email, password string) error {
  // Create a bcrypt hash of the plain-text password.
  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
  if err != nil {
	 return err
  }

  stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

  // Use the Exec() method to insert the user details and hashed password
  // into the users table.
  _, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
  if err != nil {
	 // If this returns an error, we use the errors.As() function to check
	 // whether the error has the type *mysql.MySQLError. If it does, the
	 // error will be assigned to the mySQLError variable. We can then check
	 // whether or not the error relates to our users_uc_email key by
	 // checking if the error code equals 1062 and the contents of the error
	 // message string. If it does, we return an ErrDuplicateEmail error.
	 var mySQLError *mysql.MySQLError

	 if errors.As(err, &mySQLError) {
		if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
		  return ErrDuplicateEmail
		}
	 }

	 return err
  }

  return nil
}

// Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
  return 0, nil
}

// Exists method to check if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
  return false, nil
}