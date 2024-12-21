package adapters

import (
	"errors"

	"github.com/phetployst/art-toys-store/modules/product/entities"
	"github.com/phetployst/art-toys-store/modules/product/usecase"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type gormProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) usecase.ProductRepository {
	return &gormProductRepository{db}
}

func (r *gormProductRepository) InsertProduct(product *entities.Product) (*entities.Product, error) {
	if result := r.db.Create(&product); result.Error != nil {
		return nil, result.Error
	}

	return product, nil
}

func (r *gormProductRepository) GetAllProduct() ([]entities.Product, error) {
	var products []entities.Product

	result := r.db.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *gormProductRepository) GetProductById(id string) (*entities.Product, error) {
	product := new(entities.Product)

	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *gormProductRepository) UpdateProduct(product *entities.Product, id string) (*entities.Product, error) {
	if result := r.db.Model(&entities.Product{}).
		Where("id = ?", id).
		Updates(product); result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}

func (r *gormProductRepository) UpdateStock(id string, count int) (int, error) {
	var newStock int

	err := r.db.Transaction(func(tx *gorm.DB) error {
		product := &entities.Product{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(product, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("product not found")
			}
			return errors.New("failed to retrieve product")
		}

		if product.Stock < count {
			return errors.New("insufficient stock")
		}

		newStock = product.Stock - count

		updates := map[string]interface{}{
			"stock":  newStock,
			"active": newStock > 0,
		}
		if err := tx.Model(product).Updates(updates).Error; err != nil {
			return errors.New("failed to update product stock")
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return newStock, nil
}

func (r *gormProductRepository) SearchProducts(keyword string) ([]entities.Product, error) {
	var products []entities.Product

	if err := r.db.Where("(name ILIKE ? OR description ILIKE ?) AND active = ?",
		"%"+keyword+"%", "%"+keyword+"%", true).
		Find(&products).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}

		return nil, err
	}

	return products, nil
}
