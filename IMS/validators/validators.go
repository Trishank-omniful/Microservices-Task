package validators

import (
	"errors"
	"strings"

	"github.com/Trishank-Omniful/Onboarding-Task/models"
)

func ValidateHub(hub *models.Hub) error {
	if strings.TrimSpace(hub.Name) == "" {
		return errors.New("hub name is required")
	}

	if strings.TrimSpace(hub.Address) == "" {
		return errors.New("hub address is required")
	}

	if len(hub.Name) > 255 {
		return errors.New("hub name too long (max 255 characters)")
	}

	if len(hub.Address) > 512 {
		return errors.New("hub address too long (max 512 characters)")
	}

	return nil
}

func ValidateSKU(sku *models.SKU) error {
	if strings.TrimSpace(sku.Code) == "" {
		return errors.New("SKU code is required")
	}

	if strings.TrimSpace(sku.Name) == "" {
		return errors.New("SKU name is required")
	}

	if strings.TrimSpace(sku.TenantId) == "" {
		return errors.New("tenant ID is required")
	}

	if strings.TrimSpace(sku.SellerId) == "" {
		return errors.New("seller ID is required")
	}

	if len(sku.Code) > 255 {
		return errors.New("SKU code too long (max 255 characters)")
	}

	if len(sku.Name) > 255 {
		return errors.New("SKU name too long (max 255 characters)")
	}

	return nil
}

func ValidateInventory(inventory *models.Inventory) error {
	if inventory.HubID == 0 {
		return errors.New("hub ID is required")
	}

	if inventory.SKUID == 0 {
		return errors.New("SKU ID is required")
	}

	if inventory.Quantity < 0 {
		return errors.New("inventory quantity cannot be negative")
	}

	return nil
}

func ValidateBatchSize(size int) error {
	if size == 0 {
		return errors.New("batch cannot be empty")
	}

	if size > 1000 {
		return errors.New("batch size too large (max 1000 items)")
	}

	return nil
}
