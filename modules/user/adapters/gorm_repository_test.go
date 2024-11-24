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
	getRefreshTokenByUserIDQuery   = `SELECT * FROM "credentials" WHERE user_id = $1 AND "credentials"."deleted_at" IS NULL ORDER BY created_at DESC,"credentials"."id" LIMIT $2`
	getUserProfileByIDQuery        = `SELECT * FROM "user_profiles" WHERE (user_id = $1 AND deleted_at IS NULL) AND "user_profiles"."deleted_at" IS NULL ORDER BY "user_profiles"."id" LIMIT $2`
	updateUserProfileQuery         = `UPDATE "user_profiles" SET "updated_at"=$1,"user_id"=$2,"username"=$3,"first_name"=$4,"last_name"=$5,"email"=$6,"street"=$7,"city"=$8,"state"=$9,"postal_code"=$10,"country"=$11,"profile_picture_url"=$12 WHERE user_id = $13 AND "user_profiles"."deleted_at" IS NULL`
	getAllUserProfileQuery         = `SELECT * FROM "user_profiles" WHERE "user_profiles"."deleted_at" IS NULL`
	insertUserProfileQuery         = `INSERT INTO "user_profiles" ("created_at","updated_at","deleted_at","user_id","username","first_name","last_name","email","street","city","state","postal_code","country","profile_picture_url") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`
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

func TestGetRefreshTokenByUserID_gormRepo(t *testing.T) {
	t.Run("get refresh token by user ID when record exists", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(14)
		expectedToken := "sample_refresh_token"

		rows := sqlmock.NewRows([]string{"user_id", "refresh_token", "created_at"}).
			AddRow(userID, expectedToken, time.Now())

		mock.ExpectQuery(getRefreshTokenByUserIDQuery).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		token, err := repo.GetRefreshTokenByUserID(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("get refresh token by user ID when record does not exist", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(14)

		mock.ExpectQuery(getRefreshTokenByUserIDQuery).
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		token, err := repo.GetRefreshTokenByUserID(userID)

		assert.Error(t, err)
		assert.Equal(t, "", token)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("get refresh token by user ID when query fails", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userID := uint(14)

		mock.ExpectQuery(getRefreshTokenByUserIDQuery).
			WithArgs(userID, 1).
			WillReturnError(errors.New("database error"))

		token, err := repo.GetRefreshTokenByUserID(userID)

		assert.Error(t, err)
		assert.Equal(t, "", token)
		assert.EqualError(t, err, "database error")
	})
}

func TestGetUserProfileByID_gormRepo(t *testing.T) {
	t.Run("successfully retrieves user profile by ID", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		rows := sqlmock.NewRows([]string{
			"user_id", "username", "first_name", "last_name", "email",
			"profile_picture_url", "street", "city", "state", "postal_code", "country",
		})
		rows.AddRow(
			14, "phetploy", "Phet", "Ploy", "phetploy@example.com",
			"https://example.com/profiles/14.jpg", "123 Green Lane", "Bangkok", "Central", "10110", "Thailand",
		)

		mock.ExpectQuery(getUserProfileByIDQuery).
			WithArgs("14", 1).
			WillReturnRows(rows)

		got, err := repo.GetUserProfileByID("14")

		want := &entities.UserProfile{
			UserID:            14,
			Username:          "phetploy",
			FirstName:         "Phet",
			LastName:          "Ploy",
			Email:             "phetploy@example.com",
			ProfilePictureURL: "https://example.com/profiles/14.jpg",
			Address: entities.Address{
				Street:     "123 Green Lane",
				City:       "Bangkok",
				State:      "Central",
				PostalCode: "10110",
				Country:    "Thailand",
			},
		}

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("user profile not found", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		mock.ExpectQuery(getUserProfileByIDQuery).
			WithArgs("18", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		profile, err := repo.GetUserProfileByID("18")

		assert.Error(t, err)
		assert.Nil(t, profile)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("database error during query", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		mock.ExpectQuery(`SELECT * FROM "user_profiles" WHERE user_id = \$1 AND deleted_at IS NULL ORDER BY "user_profiles"."user_id" LIMIT 1`).
			WithArgs("21", 1).
			WillReturnError(errors.New("database error"))

		profile, err := repo.GetUserProfileByID("21")

		assert.Error(t, err)
		assert.Nil(t, profile)
	})
}

func TestUpdateUserProfile_gormRepo(t *testing.T) {
	t.Run("successfully updates user profile", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		updateInput := &entities.UserProfile{
			UserID: 14, Username: "phetploy", FirstName: "Phet", LastName: "Ploy", Email: "phetploy@example.com",
			Address:           entities.Address{Street: "123 Green Lane", City: "Bangkok", State: "Central", PostalCode: "10110", Country: "Thailand"},
			ProfilePictureURL: "https://example.com/profiles/14.jpg",
		}

		mock.ExpectBegin()

		mock.ExpectExec(updateUserProfileQuery).
			WithArgs(
				sqlmock.AnyArg(), 14, "phetploy", "Phet", "Ploy", "phetploy@example.com",
				"123 Green Lane", "Bangkok", "Central", "10110", "Thailand",
				"https://example.com/profiles/14.jpg", 14,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		updatedProfile, err := repo.UpdateUserProfile(updateInput)

		assert.NoError(t, err)

		if !reflect.DeepEqual(updatedProfile, updateInput) {
			t.Errorf("got %v but want %v", updatedProfile, updateInput)
		}

	})

	t.Run("user profile not found during update", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userProfile := &entities.UserProfile{
			UserID:   18,
			Username: "chopper",
		}

		mock.ExpectExec(updateUserProfileQuery).
			WithArgs(18, "chopper").
			WillReturnResult(sqlmock.NewResult(0, 0))

		updatedProfile, err := repo.UpdateUserProfile(userProfile)

		assert.Error(t, err)
		assert.Nil(t, updatedProfile)
	})

	t.Run("database error during update", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userProfile := &entities.UserProfile{
			UserID:   21,
			Username: "tonytony",
		}

		mock.ExpectExec(updateUserProfileQuery).
			WithArgs(21, "tonytony").
			WillReturnError(errors.New("database error"))

		updatedProfile, err := repo.UpdateUserProfile(userProfile)

		assert.Error(t, err)
		assert.Nil(t, updatedProfile)
	})
}

func TestGetAllUserProfile_gormRepo(t *testing.T) {
	t.Run("get all user profiles successfully", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := &gormUserRepository{db: gormDB}

		rows := sqlmock.NewRows([]string{
			"user_id", "username", "first_name", "last_name", "email",
			"street", "city", "state", "postal_code", "country",
			"profile_picture_url",
		}).AddRow(
			31, "phetploy", "Phet", "Ploy", "phetploy@example.com",
			"123 Green Lane", "Bangkok", "Central", "10110", "Thailand",
			"https://example.com/profiles/31.jpg",
		).AddRow(
			32, "tonytonychopper", "Tony", "Chopper", "tonychopper@example.com",
			"456 Blue Street", "Chiang Mai", "North", "50200", "Thailand",
			"https://example.com/profiles/32.jpg",
		)

		mock.ExpectQuery(getAllUserProfileQuery).
			WillReturnRows(rows)

		count, profiles, err := repo.GetAllUserProfile()

		expectedProfiles := []entities.UserProfile{
			{UserID: 31, Username: "phetploy", FirstName: "Phet", LastName: "Ploy", Email: "phetploy@example.com",
				Address:           entities.Address{Street: "123 Green Lane", City: "Bangkok", State: "Central", PostalCode: "10110", Country: "Thailand"},
				ProfilePictureURL: "https://example.com/profiles/31.jpg"},
			{UserID: 32, Username: "tonytonychopper", FirstName: "Tony", LastName: "Chopper", Email: "tonychopper@example.com",
				Address:           entities.Address{Street: "456 Blue Street", City: "Chiang Mai", State: "North", PostalCode: "50200", Country: "Thailand"},
				ProfilePictureURL: "https://example.com/profiles/32.jpg"},
		}

		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
		assert.Equal(t, expectedProfiles, profiles)
	})

	t.Run("database error while fetching profiles", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := &gormUserRepository{db: gormDB}

		mock.ExpectQuery(getAllUserProfileQuery).
			WillReturnError(errors.New("database error"))

		count, profiles, err := repo.GetAllUserProfile()

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		assert.Equal(t, int64(0), count)
		assert.Nil(t, profiles)
	})
}

func TestInsertUserProfile_gormRepo(t *testing.T) {
	t.Run("insert user profile successfully", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userProfile := &entities.UserProfile{
			UserID: 31, Username: "phetploy", FirstName: "Phet", LastName: "Ploy", Email: "phetploy@example.com",
			Address:           entities.Address{Street: "123 Green Lane", City: "Bangkok", State: "Central", PostalCode: "10110", Country: "Thailand"},
			ProfilePictureURL: "https://example.com/profiles/31.jpg",
		}

		mock.ExpectBegin()
		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(insertUserProfileQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 31, "phetploy", "Phet", "Ploy", "phetploy@example.com",
				"123 Green Lane", "Bangkok", "Central", "10110", "Thailand", "https://example.com/profiles/31.jpg").
			WillReturnRows(row)
		mock.ExpectCommit()

		err := repo.InsertUserProfile(userProfile)

		assert.NoError(t, err)
	})

	t.Run("database error during query", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewUserRepository(gormDB)

		userProfile := &entities.UserProfile{
			UserID: 31, Username: "phetploy", FirstName: "Phet", LastName: "Ploy", Email: "phetploy@example.com",
			Address:           entities.Address{Street: "123 Green Lane", City: "Bangkok", State: "Central", PostalCode: "10110", Country: "Thailand"},
			ProfilePictureURL: "https://example.com/profiles/31.jpg",
		}

		mock.ExpectQuery(insertUserProfileQuery).
			WillReturnError(errors.New("database error"))

		err := repo.InsertUserProfile(userProfile)

		assert.Error(t, err)
	})

}
