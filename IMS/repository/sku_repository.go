package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Trishank-Omniful/Onboarding-Task/constants"
	"github.com/Trishank-Omniful/Onboarding-Task/models"
	"github.com/omniful/go_commons/redis"
	"gorm.io/gorm"
)

type SkuRepository struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewSkuRepository(db *gorm.DB, redis *redis.Client) *SkuRepository {
	return &SkuRepository{DB: db, Redis: redis}
}

func getSKUIDCacheKey(id uint) string {
	return fmt.Sprintf("%s%d", constants.CacheKeySKUID, id)
}

func (r *SkuRepository) GetAllSkus() ([]models.SKU, error) {
	var skus []models.SKU
	result := r.DB.Find(&skus)
	return skus, result.Error
}

func (r *SkuRepository) GetSkuById(id uint) (*models.SKU, error) {
	ctx := context.Background()
	cacheKey := getSKUIDCacheKey(id)
	var sku models.SKU

	val, err := r.Redis.Get(ctx, cacheKey)
	if err == nil {
		if jsonErr := json.Unmarshal([]byte(val), &sku); jsonErr == nil {
			log.Printf("SKU retrieved from cache: %d", id)
			return &sku, nil
		}
		log.Println("Failed to unmarshal SKU from Redis: ", err)
	} else if err != r.Redis.Nil {
		log.Println("Redis error while fetching SKU: ", err)
	}
	sku = models.SKU{}
	result := r.DB.First(&sku, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	HubJSON, jsonErr := json.Marshal(sku)
	if jsonErr != nil {
		log.Println("Failed to marshal SKU for Redis", jsonErr)
	} else {
		if success, err := r.Redis.Set(ctx, cacheKey, string(HubJSON), time.Duration(constants.CacheTTLSKUs)*time.Minute); err != nil {
			log.Print("Failed to Set SKU in Redis", err)
		} else if success {
			log.Print("SKU set in Redis Cache: ", cacheKey)
		}
	}
	return &sku, result.Error
}

func (r *SkuRepository) CreateSku(sku *models.SKU) error {
	result := r.DB.Create(sku)
	return result.Error
}

func (r *SkuRepository) UpdateSku(sku *models.SKU) error {
	result := r.DB.Model(sku).Updates(sku)
	if result.Error == nil {
		r.Redis.Del(context.Background(), getSKUIDCacheKey(sku.ID))
		log.Print("SKU Cache Invalidated after Update")
	}
	return result.Error
}

func (r *SkuRepository) DeleteSku(id uint) error {
	var sku models.SKU
	result := r.DB.Delete(&sku, id)
	if result.Error == nil {
		r.Redis.Del(context.Background(), getSKUIDCacheKey(id))
		log.Print("SKU Cache Invalidated after Update")
		return nil
	}
	return result.Error
}

func (r *SkuRepository) GetSkusByTenantAndSeller(tenantID string, sellerID string, skuCodes []string) ([]models.SKU, error) {
	query := r.DB.Model(&models.SKU{})

	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if sellerID != "" {
		query = query.Where("seller_id = ?", sellerID)
	}

	if len(skuCodes) > 0 {
		query = query.Where("code IN (?)", skuCodes)
	}

	var skus []models.SKU
	err := query.Find(&skus).Error
	return skus, err
}

func (r *SkuRepository) CreateSKUsBatch(skus []models.SKU) error {
	if len(skus) == 0 {
		return nil
	}
	result := r.DB.CreateInBatches(skus, 100)
	return result.Error
}

func (r *SkuRepository) GetSKUsByIDs(ids []uint) ([]models.SKU, error) {
	var skus []models.SKU
	result := r.DB.Where("id IN (?)", ids).Find(&skus)
	return skus, result.Error
}

func (r *SkuRepository) GetSKUsByCodes(codes []string) ([]models.SKU, error) {
	var skus []models.SKU
	result := r.DB.Where("code IN (?)", codes).Find(&skus)
	return skus, result.Error
}
