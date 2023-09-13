package user

import (
	"context"
	"errors"
	"net/http"

	connector "github.com/SamstreamJet/dialogo/server/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id         int
	Username   string
	Email      string
	Password   string
	Created_at string
}

func Login(ctx context.Context, email string, password string) error {
	user, err := GetUserByEmail(ctx, email)
	if err != nil {
		ctx = context.WithValue(ctx, "httpStatus", http.StatusInternalServerError)
		return err
	}

	if user.Email == "" {
		ctx = context.WithValue(ctx, "httpStatus", http.StatusUnauthorized)
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx = context.WithValue(ctx, "httpStatus", http.StatusUnauthorized)
		return err
	}
	return nil
}

func Register(ctx context.Context, email string, username string, password string) error {
	user := new(User)

	// TODO: ADD VALIDATION FOR EMAIL HERE
	user.Email = email
	user.Username = username

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		ctx = context.WithValue(ctx, "httpStatus", http.StatusInternalServerError)
		return err
	}
	user.Password = string(hashedPasswordBytes)

	err = create(user)
	if err != nil {
		ctx = context.WithValue(ctx, "httpStatus", http.StatusInternalServerError)
		return err
	}

	return nil
}

func create(user *User) error {
	_, err := connector.Query(
		`insert into "user" ("email", "username", "password")
		values ($1, $2, $3)`,
		user.Email, user.Username, user.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(ctx context.Context, email string) (User, error) {
	user := []User{}
	query := `select * from "user" where email = $1`

	err := connector.Select(&user, query, email)
	if err != nil {
		return *new(User), err
	}
	if len(user) == 0 {
		return *new(User), nil
	}
	return user[0], nil
}

func GetAllUsers(ctx context.Context) ([]User, error) {
	users := []User{}
	conn, err := connector.GetConnx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = conn.SelectContext(ctx, &users, `select * from "user"`)
	if err != nil {
		return nil, err
	}

	// err := connector.Select(&users, `select * from "user"`)
	// if err != nil {
	// 	return nil, err
	// }

	// for result.Next() {
	// 	var user User
	// 	if err := result.Scan(&user.id, &user.password, &user.email, &user.email, &user.created_at); err != nil {
	// 		return nil, err
	// 	}
	// 	users = append(users, user)
	// }
	// if err := result.Err(); err != nil {
	// 	return nil, err
	// }

	return users, nil
}
