package repository

import (
	"errors"

	"github.com/Trishank-Omniful/Onboarding-Task/models"
	"github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoryRepository struct {
	DB      *gorm.DB
	Redis   *redis.Client
	SKURepo *SkuRepository
	HubRepo *HubRepository
}

func NewInventoryRepository(
	db *gorm.DB,
	redis *redis.Client,
	skuRepo *SkuRepository,
	hubRepo *HubRepository,
) *InventoryRepository {
	return &InventoryRepository{
		DB:      db,
		Redis:   redis,
		SKURepo: skuRepo,
		HubRepo: hubRepo,
	}
}

func (r *InventoryRepository) UpsertInventory(inventory *models.Inventory) error {
	upsert := r.DB.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "hub_id"}, {Name: "sku_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"quantity": inventory.Quantity}),
		},
	).Create(inventory)
	return upsert.Error
}

func (r *InventoryRepository) GetInventoryByHubAndSKU(hubID, skuID uint) (*models.Inventory, error) {
	var inventory models.Inventory
	result := r.DB.Where("hub_id = ? AND sku_id = ?", hubID, skuID).First(&inventory)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &models.Inventory{HubID: hubID, SKUID: skuID, Quantity: 0}, gorm.ErrRecordNotFound
	}
	return &inventory, result.Error
}

func (r *InventoryRepository) GetInventoriesFiltered(skuCode *string, hubID *uint) ([]models.Inventory, error) {
	var inventories []models.Inventory
	query := r.DB.Model(&models.Inventory{})

	if skuCode != nil && *skuCode != "" {
		query = query.Joins("JOIN skus on skus.id = inventories.sku_id").Where("skus.code = ?", *skuCode)
	}

	if hubID != nil && *hubID != 0 {
		query = query.Where("inventories.hub_id = ?", *hubID)
	}

	result := query.Find(&inventories)
	if result.Error != nil {
		return nil, result.Error
	}
	return inventories, nil
}

func (r *InventoryRepository) GetInventory(hubID, skuID string) ([]models.Inventory, error) {
	var inventories []models.Inventory
	query := r.DB.Model(&models.Inventory{}).Preload("Hub").Preload("SKU")

	if hubID != "" {
		query = query.Where("hub_id = ?", hubID)
	}

	if skuID != "" {
		query = query.Where("sku_id = ?", skuID)
	}

	result := query.Find(&inventories)
	return inventories, result.Error
}

func (r *InventoryRepository) ReduceInventory(hubID, skuID uint, quantityToReduce int) error {
	if quantityToReduce <= 0 {
		return errors.New("quantity to reduce must be positive")
	}
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var inventory models.Inventory
		result := tx.Set("gorm:query_option", "FOR UPDATE").Where("hub_id = ? AND sku_id = ?", hubID, skuID).First(&inventory)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("inventory record not found for reduction")
		}
		if result.Error != nil {
			return result.Error
		}
		if inventory.Quantity < quantityToReduce {
			return errors.New("insufficient inventory")
		}
		inventory.Quantity -= quantityToReduce
		if err := tx.Save(&inventory).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *InventoryRepository) UpsertInventoryBatch(inventories []models.Inventory) error {
	if len(inventories) == 0 {
		return nil
	}

	return r.DB.Transaction(func(tx *gorm.DB) error {
		for _, inventory := range inventories {
			upsert := tx.Clauses(
				clause.OnConflict{
					Columns:   []clause.Column{{Name: "hub_id"}, {Name: "sku_id"}},
					DoUpdates: clause.Assignments(map[string]interface{}{"quantity": inventory.Quantity}),
				},
			).Create(&inventory)
			if upsert.Error != nil {
				return upsert.Error
			}
		}
		return nil
	})
}

func (r *InventoryRepository) GetInventoriesByHubAndSKUs(hubID uint, skuIDs []uint) ([]models.Inventory, error) {
	var inventories []models.Inventory
	query := r.DB.Model(&models.Inventory{}).Preload("Hub").Preload("SKU").Where("hub_id = ?", hubID)

	if len(skuIDs) > 0 {
		query = query.Where("sku_id IN (?)", skuIDs)
	}

	result := query.Find(&inventories)
	return inventories, result.Error
}
