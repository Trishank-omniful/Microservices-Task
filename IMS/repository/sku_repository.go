package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

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

func getSkuCacheKey(id uint) string {
	return fmt.Sprintf("sku:%d", id)
}

func (r *SkuRepository) GetAllSkus() ([]models.SKU, error) {
	var skus []models.SKU
	result := r.DB.Find(&skus)
	return skus, result.Error
}

func (r *SkuRepository) GetSkuById(id uint) (*models.SKU, error) {
	ctx := context.Background()
	cacheKey := getSkuCacheKey(id)
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
		if success, err := r.Redis.Set(ctx, cacheKey, string(HubJSON), 5*time.Minute); err != nil {
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
		r.Redis.Del(context.Background(), getSkuCacheKey(sku.ID))
		log.Print("SKU Cache Invalidated after Update")
	}
	return result.Error
}

func (r *SkuRepository) DeleteSku(id uint) error {
	var sku models.SKU
	result := r.DB.Delete(&sku, id)
	if result.Error == nil {
		r.Redis.Del(context.Background(), getSkuCacheKey(id))
		log.Print("SKU Cache Invalidated after Update")
	}
	return result.Error
}

func (r *SkuRepository) GetSkuByName(name string) (*models.SKU, error) {
	var sku models.SKU
	result := r.DB.Where("name = ?", name).First(&sku)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("sku not found by name")
	}
	return &sku, result.Error
}
