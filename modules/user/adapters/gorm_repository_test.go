package adapters

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/phetployst/art-toys-store/modules/user/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	createUserQuery                = `INSERT INTO "users" ("created_at","updated_at","deleted_at","username","email","password_hash","role") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`
	isUniqueUserQuery              = `SELECT * FROM "users" WHERE (email = $1 OR username = $2) AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $3`
	getUserAccountByIdQuery        = `SELECT * FROM "users" WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`
	getUserAccountByUsernameQuery  = `SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`
	insertUserCredentialQuery      = `INSERT INTO "credentials" ("created_at","updated_at","deleted_at","user_id","refresh_token","expires_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`
	getUserCredentialByUserIdQuery = `SELECT * FROM "credentials" WHERE (user_id = $1 AND deleted_at IS NULL) AND "credentials"."deleted_at" IS NULL ORDER BY "credentials"."id" LIMIT $2`
	deleteUserCredentialQuery      = `DELETE FROM "credentials" WHERE user_id = $1`
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

func TestGetUserAccountByUsername_gormRepo(t *testing.T) {
	t.Run("get user account by username given users exist in the database", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		username := "phetploy"

		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "role"})
		rows.AddRow(uint(13), "phetploy", "phetploy@example.com", "password1234", "user")

		mock.ExpectQuery(getUserAccountByUsernameQuery).
			WithArgs(username, 1).
			WillReturnRows(rows)

		got, err := repo.GetUserByUsername(username)

		want := &entities.User{
			Model:    gorm.Model{ID: 13, CreatedAt: got.CreatedAt, UpdatedAt: got.UpdatedAt, DeletedAt: got.DeletedAt},
			Username: "phetploy",
			Email:    "phetploy@example.com",
			Role:     "user",
		}

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("get user by username given user does not exist", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		username := "phetploy"

		mock.ExpectQuery(getUserAccountByUsernameQuery).
			WithArgs(username, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := repo.GetUserByUsername(username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	})

	t.Run("get user by username given database error", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		username := "phetploy"

		mock.ExpectQuery(getUserAccountByUsernameQuery).
			WithArgs(username, 1).
			WillReturnError(errors.New("database error"))

		user, err := repo.GetUserByUsername(username)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.EqualError(t, err, "database error")
	})
}

func TestInsertUserCredential_gormRepo(t *testing.T) {
	t.Run("insert user credential successfully", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userCredential := entities.Credential{UserID: uint(14), RefreshToken: "refreshToken1234", ExpiresAt: time.Time{}}

		mock.ExpectBegin()
		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(insertUserCredentialQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userCredential.UserID, userCredential.RefreshToken, userCredential.ExpiresAt).
			WillReturnRows(row)
		mock.ExpectCommit()

		err := repo.InsertUserCredential(&userCredential)

		assert.NoError(t, err)

	})

	t.Run("insert user credential given error during query", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userCredential := entities.Credential{UserID: uint(14), RefreshToken: "refreshToken1234", ExpiresAt: time.Time{}}

		mock.ExpectBegin()
		mock.ExpectQuery(insertUserCredentialQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userCredential.UserID, userCredential.RefreshToken, userCredential.ExpiresAt).
			WillReturnError(errors.New("database error"))
		mock.ExpectCommit()

		err := repo.InsertUserCredential(&userCredential)

		assert.Error(t, err)
	})
}

func TestGetUserCredentialByUserId_gormRepo(t *testing.T) {
	t.Run("get credential by user ID given user login exis in database", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(14)

		rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "role"})
		rows.AddRow(uint(14), "phetploy", "phetploy@example.com", "password1234", "user")

		mock.ExpectQuery(getUserCredentialByUserIdQuery).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		err := repo.GetUserCredentialByUserId(userID)

		assert.NoError(t, err)

	})

	t.Run("get credential by user ID given user login not exis in database", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(14)

		mock.ExpectQuery(getUserCredentialByUserIdQuery).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		err := repo.GetUserCredentialByUserId(userID)

		assert.Error(t, err)

	})

	t.Run("get credential by user ID given error during query", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(14)

		mock.ExpectQuery(getUserCredentialByUserIdQuery).
			WithArgs(userID, 1).
			WillReturnError(errors.New("database error"))

		err := repo.GetUserCredentialByUserId(userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")

	})
}

func TestDeleteUserCredential_gormRepo(t *testing.T) {
	t.Run("delete user credential by user ID given success", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(19)

		mock.ExpectBegin()
		mock.ExpectExec(deleteUserCredentialQuery).
			WithArgs(userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteUserCredential(userID)

		assert.NoError(t, err)

	})

	t.Run("delete user credential by user ID given error during guery", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(19)

		mock.ExpectBegin()
		mock.ExpectExec(deleteUserCredentialQuery).
			WithArgs(userID).
			WillReturnError(errors.New("database error"))
		mock.ExpectCommit()

		err := repo.DeleteUserCredential(userID)

		assert.Error(t, err)
	})
}
