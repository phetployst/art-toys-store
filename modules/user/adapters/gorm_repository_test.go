package adapters

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	createUserQuery         = `INSERT INTO "users" ("created_at","updated_at","deleted_at","username","email","password_hash","role") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
	isUniqueUserQuery       = `SELECT * FROM "users" WHERE (email = $1 OR username = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`
	getUserAccountByIdQuery = `SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`
)

func TestCreateUser_gormRepo(t *testing.T) {

	t.Run("create user successfully", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		user := entities.User{Username: "phetploy", Email: "phetploy@example.com", PasswordHash: "password1234hash", Role: "user"}

		mock.ExpectBegin()
		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(createUserQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Username, user.Email, user.PasswordHash, user.Role).
			WillReturnRows(row)
		mock.ExpectCommit()

		want := uint(1)
		got, err := repo.CreateUser(&user)

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("create user given error during query", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		user := entities.User{Username: "phetploy", Email: "phetploy@example.com", PasswordHash: "password1234hash", Role: "user"}

		mock.ExpectBegin()
		mock.ExpectQuery(createUserQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), user.Username, user.Email, user.PasswordHash, user.Role).
			WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		want := uint(0)
		got, err := repo.CreateUser(&user)

		assert.Error(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})
}

func TestIsUniqueUser_gormRepo(t *testing.T) {
	t.Run("returns true when username and email is unique", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		user := entities.User{Email: "chopper@example.com", Username: "tonytony"}

		mock.ExpectQuery(isUniqueUserQuery).
			WithArgs(user.Email, user.Username, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result := repo.IsUniqueUser(user.Email, user.Username)

		assert.True(t, result)
	})

	t.Run("returns false when username and email already exists", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		user := entities.User{Email: "chopper@example.com", Username: "tonytony"}

		rows := sqlmock.NewRows([]string{"id", "email", "username"})
		rows.AddRow(1, user.Email, user.Username)

		mock.ExpectQuery(isUniqueUserQuery).
			WithArgs(user.Email, user.Username, 1).
			WillReturnRows(rows)

		result := repo.IsUniqueUser(user.Email, user.Username)

		assert.False(t, result)
	})
}

func TestGetUserAccountById_gormRepo(t *testing.T) {

	t.Run("get user account by id given users exist in the database", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(1)

		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "role"})
		rows.AddRow(userID, "phetploy", "phetploy@example.com", "password1234", "user")

		mock.ExpectQuery(getUserAccountByIdQuery).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		got, err := repo.GetUserAccountById(userID)

		want := &entities.User{
			Model:    gorm.Model{ID: 1, CreatedAt: got.CreatedAt, UpdatedAt: got.UpdatedAt, DeletedAt: got.DeletedAt},
			Username: "phetploy",
			Email:    "phetploy@example.com",
			Role:     "user",
		}

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("get user by id given user does not exist", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(2)

		mock.ExpectQuery(getUserAccountByIdQuery).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetUserAccountById(userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	})
}
