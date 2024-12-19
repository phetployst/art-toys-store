package adapters

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/phetployst/art-toys-store/modules/product/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	insertProductQuery  = `INSERT INTO "products" ("created_at","updated_at","deleted_at","name","description","price","stock","image_url","active") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`
	getAllProductQuery  = `SELECT * FROM "products" WHERE "products"."deleted_at" IS NULL`
	getProductByIdQuery = `SELECT * FROM "products" WHERE "products"."id" = $1 AND "products"."deleted_at" IS NULL ORDER BY "products"."id" LIMIT $2`
)

func TestInsertProduct_gormRepo(t *testing.T) {
	t.Run("insert new product successfully", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewProductRepository(gormDB)

		newProduct := &entities.Product{
			Name: "Molly Classic", Description: "The iconic Molly figure, loved by art toy collectors worldwide.",
			Price: 340.99, Stock: 30, ImageURL: "https://example.com/images/molly-classic.jpg", Active: true,
		}

		mock.ExpectBegin()
		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(insertProductQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), newProduct.Name, newProduct.Description, newProduct.Price, newProduct.Stock, newProduct.ImageURL, newProduct.Active).
			WillReturnRows(row)
		mock.ExpectCommit()

		got, err := repo.InsertProduct(newProduct)

		want := &entities.Product{
			Model:       gorm.Model{ID: 1, CreatedAt: got.CreatedAt, UpdatedAt: got.UpdatedAt, DeletedAt: got.DeletedAt},
			Name:        "Molly Classic",
			Description: "The iconic Molly figure, loved by art toy collectors worldwide.",
			Price:       340.99,
			Stock:       30,
			ImageURL:    "https://example.com/images/molly-classic.jpg",
			Active:      true,
		}

		assert.NoError(t, err)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("insert new product given user does not exist", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewProductRepository(gormDB)

		newProduct := &entities.Product{
			Name: "Molly Classic", Description: "The iconic Molly figure, loved by art toy collectors worldwide.",
			Price: 340.99, Stock: 30, ImageURL: "https://example.com/images/molly-classic.jpg", Active: true,
		}

		mock.ExpectQuery(insertProductQuery).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), newProduct.Name, newProduct.Description, newProduct.Price, newProduct.Stock, newProduct.ImageURL, newProduct.Active).
			WillReturnError(gorm.ErrRecordNotFound)

		got, err := repo.InsertProduct(newProduct)

		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestGetAllProducts_gormRepo(t *testing.T) {
	t.Run("get all product successfully", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewProductRepository(gormDB)

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "price", "stock", "image_url", "active",
		}).AddRow(1, "Dimoo Starry Night", "Dimoo inspired by Van Gogh's 'Starry Night,' featuring a dreamy and artistic design.", 49.99, 25, "https://example.com/images/dimoo-starry-night.jpg", true).
			AddRow(2, "Pucky Forest Fairy", "A magical art toy figure from Pucky, with a whimsical forest fairy design.", 44.99, 40, "https://example.com/images/pucky-forest-fairy.jpg", true)

		mock.ExpectQuery(getAllProductQuery).WillReturnRows(rows)

		got, err := repo.GetAllProduct()

		assert.NoError(t, err)

		want := []entities.Product{
			{Model: gorm.Model{ID: 1}, Name: "Dimoo Starry Night", Description: "Dimoo inspired by Van Gogh's 'Starry Night,' featuring a dreamy and artistic design.", Price: 49.99, Stock: 25, ImageURL: "https://example.com/images/dimoo-starry-night.jpg", Active: true},
			{Model: gorm.Model{ID: 2}, Name: "Pucky Forest Fairy", Description: "A magical art toy figure from Pucky, with a whimsical forest fairy design.", Price: 44.99, Stock: 40, ImageURL: "https://example.com/images/pucky-forest-fairy.jpg", Active: true},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("get all product given database error", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewProductRepository(gormDB)

		mock.ExpectQuery(getAllProductQuery).
			WillReturnError(errors.New("database error"))

		_, err := repo.GetAllProduct()

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
	})
}

func TestGetProductById_gormRepo(t *testing.T) {
	t.Run("get product by id successfully", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewProductRepository(gormDB)

		rows := sqlmock.NewRows([]string{
			"id", "name", "description", "price", "stock", "image_url", "active",
		}).AddRow(1, "Dimoo Starry Night", "Dimoo inspired by Van Gogh's 'Starry Night,' featuring a dreamy and artistic design.", 49.99, 25, "https://example.com/images/dimoo-starry-night.jpg", true)

		mock.ExpectQuery(getProductByIdQuery).
			WithArgs("1", 1).
			WillReturnRows(rows)

		got, err := repo.GetProductById("1")

		assert.NoError(t, err)

		want := &entities.Product{
			Model:       gorm.Model{ID: 1, CreatedAt: time.Time{}, UpdatedAt: time.Time{}},
			Name:        "Dimoo Starry Night",
			Description: "Dimoo inspired by Van Gogh's 'Starry Night,' featuring a dreamy and artistic design.",
			Price:       49.99,
			Stock:       25,
			ImageURL:    "https://example.com/images/dimoo-starry-night.jpg",
			Active:      true,
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v but want %v", got, want)
		}
	})

	t.Run("get product by id given database error", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
		repo := NewProductRepository(gormDB)

		mock.ExpectQuery(getProductByIdQuery).
			WithArgs("1", 1).
			WillReturnError(errors.New("database error"))

		got, err := repo.GetProductById("1")

		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
