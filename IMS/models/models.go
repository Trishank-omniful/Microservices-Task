package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Hub struct {
	gorm.Model
	Name         string `gorm:"type:varchar(255);not null;uniqueIndex" json:"name"`
	Address      string `gorm:"type:varchar(512);not null" json:"address"`
	City         string `gorm:"type:varchar(100)" json:"city"`
	State        string `gorm:"type:varchar(100)" json:"state"`
	Country      string `gorm:"type:varchar(100)" json:"country"`
	PostalCode   string `gorm:"type:varchar(10)" json:"postal_code"`
	ContactName  string `gorm:"type:varchar(255)" json:"contact_name"`
	ContactEmail string `gorm:"type:varchar(255)" json:"contact_email"`
}

type SKU struct {
	gorm.Model
	Code        string          `gorm:"type:varchar(255);not null;uniqueIndex" json:"code"`
	Name        string          `gorm:"type:varchar(255);not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	TenantId    string          `gorm:"type:varchar(255);not null;index" json:"tenant_id"`
	SellerId    string          `gorm:"type:varchar(255);not null;index" json:"seller_id"`
	Category    string          `gorm:"type:varchar(100)" json:"category"`
	Price       sql.NullFloat64 `json:"price"`
}

type Inventory struct {
	gorm.Model
	HubID    uint `gorm:"not null;uniqueIndex:idx_sku_hub" json:"hub_id"`
	Hub      Hub  `gorm:"foreignKey:HubID;constraint:OnDelete:CASCADE" json:"hub"`
	SKUID    uint `gorm:"column:sku_id;not null;uniqueIndex:idx_sku_hub" json:"sku_id"`
	SKU      SKU  `gorm:"foreignKey:SKUID;constraint:OnDelete:CASCADE" json:"sku"`
	Quantity int  `gorm:"not null;default:0" json:"quantity"`
}
