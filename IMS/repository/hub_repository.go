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

type HubRepository struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewHubRepository(db *gorm.DB, redis *redis.Client) *HubRepository {
	return &HubRepository{DB: db, Redis: redis}
}

func getHubCacheKey(id uint) string {
	return fmt.Sprintf("hub:%d", id)
}

func (r *HubRepository) GetAllHubs() ([]models.Hub, error) {
	var hubs []models.Hub
	result := r.DB.Find(&hubs)
	return hubs, result.Error
}

func (r *HubRepository) GetHubById(id uint) (*models.Hub, error) {
	ctx := context.Background()
	cacheKey := getHubCacheKey(id)
	var hub models.Hub

	val, err := r.Redis.Get(ctx, cacheKey)
	if err == nil {
		if jsonErr := json.Unmarshal([]byte(val), &hub); jsonErr == nil {
			log.Printf("HUB retrieved from cache: %d", id)
			return &hub, nil
		}
		log.Println("Failed to unmarshal HUB from Redis: ", err)
	} else if err != r.Redis.Nil {
		log.Println("Redis error while fetching HUB: ", err)
	}
	hub = models.Hub{}
	result := r.DB.First(&hub, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	HubJSON, jsonErr := json.Marshal(hub)
	if jsonErr != nil {
		log.Println("Failed to marshal HUB for Redis", jsonErr)
	} else {
		if success, err := r.Redis.Set(ctx, cacheKey, string(HubJSON), 5*time.Minute); err != nil {
			log.Print("Failed to Set HUB in Redis", err)
		} else if success {
			log.Print("HUB set in Redis Cache: ", cacheKey)
		}
	}
	return &hub, result.Error
}

func (r *HubRepository) CreateHub(hub *models.Hub) error {
	result := r.DB.Create(hub)
	return result.Error
}

func (r *HubRepository) UpdateHub(hub *models.Hub) error {
	result := r.DB.Model(hub).Updates(hub)
	if result.Error == nil {
		r.Redis.Del(context.Background(), getHubCacheKey(hub.ID))
		log.Print("SKU Cache Invalidated after Update")
	}
	return result.Error
}

func (r *HubRepository) DeleteHub(id uint) error {
	var hub models.Hub
	result := r.DB.Delete(&hub, id)
	if result.Error == nil {
		r.Redis.Del(context.Background(), getHubCacheKey(id))
		log.Print("SKU Cache Invalidated after Update")
	}
	return result.Error
}

func (r *HubRepository) GetHubByName(name string) (*models.Hub, error) {
	var hub models.Hub
	result := r.DB.Where("name = ?", name).First(&hub)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("hub not found by name")
	}
	return &hub, result.Error
}
