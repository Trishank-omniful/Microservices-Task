package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderStatus string

const (
	OrderStatusOnHold    OrderStatus = "on_hold"
	OrderStatusNew       OrderStatus = "new_order"
	OrderStatusCanceled  OrderStatus = "canceled"
	OrderStatusCompleted OrderStatus = "completed"
)

type Order struct {
	ID             primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	TenantID       string              `json:"tenant_id" bson:"tenant_id"`
	SellerID       string              `json:"seller_id" bson:"seller_id"`
	ReferenceID    string              `json:"reference_id" bson:"reference_id"`
	Status         OrderStatus         `json:"status" bson:"status"`
	Items          []OrderItem         `json:"items" bson:"items"`
	CustomerInfo   CustomerInfo        `json:"customer_info" bson:"customer_info"`
	ShippingInfo   ShippingInfo        `json:"shipping_info" bson:"shipping_info"`
	PaymentInfo    PaymentInfo         `json:"payment_info" bson:"payment_info"`
	TotalAmount    float64             `json:"total_amount" bson:"total_amount"`
	Currency       string              `json:"currency" bson:"currency"`
	OrderDate      time.Time           `json:"order_date" bson:"order_date"`
	LastUpdated    time.Time           `json:"last_updated" bson:"last_updated"`
	History        []OrderHistoryEvent `json:"history" bson:"history"`
	InvalidRowsCSV string              `json:"invalid_rows_csv,omitempty" bson:"invalid_rows_csv,omitempty"`
}

type OrderItem struct {
	SKUCode   string  `json:"sku_code" bson:"sku_code"`
	HubCode   string  `json:"hub_code" bson:"hub_code"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	UnitPrice float64 `json:"unit_price" bson:"unit_price"`
}

type CustomerInfo struct {
	FirstName string  `json:"first_name" bson:"first_name"`
	LastName  string  `json:"last_name" bson:"last_name"`
	Email     string  `json:"email" bson:"email"`
	Phone     string  `json:"phone" bson:"phone"`
	Address   Address `json:"address" bson:"address"`
}

type ShippingInfo struct {
	Method      string     `json:"method" bson:"method"`
	Cost        float64    `json:"cost" bson:"cost"`
	TrackingID  string     `json:"tracking_id,omitempty" bson:"tracking_id,omitempty"`
	ShippedDate *time.Time `json:"shipped_date,omitempty" bson:"shipped_date,omitempty"`
	Address     Address    `json:"address" bson:"address"`
}

type PaymentInfo struct {
	Method        string     `json:"method" bson:"method"`
	TransactionID string     `json:"transaction_id,omitempty" bson:"transaction_id,omitempty"`
	AmountPaid    float64    `json:"amount_paid" bson:"amount_paid"`
	PaymentStatus string     `json:"payment_status" bson:"payment_status"`
	PaidAt        *time.Time `json:"paid_at,omitempty" bson:"paid_at,omitempty"`
}

type Address struct {
	Street  string `json:"street" bson:"street"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	ZipCode string `json:"zip_code" bson:"zip_code"`
	Country string `json:"country" bson:"country"`
}

type OrderHistoryEvent struct {
	Timestamp   time.Time   `json:"timestamp" bson:"timestamp"`
	NewStatus   OrderStatus `json:"new_status" bson:"new_status"`
	OldStatus   OrderStatus `json:"old_status" bson:"old_status"`
	Description string      `json:"description" bson:"description"`
	Actor       string      `json:"actor" bson:"actor"`
}

type WebhookRegistration struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TenantID    string             `json:"tenant_id" bson:"tenant_id"`
	EventType   string             `json:"event_type" bson:"event_type"`
	CallbackURL string             `json:"callback_url" bson:"callback_url"`
	Secret      string             `json:"secret,omitempty" bson:"secret,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
}

type BulkOrderCSVRequest struct {
	TenantID string `json:"tenant_id" validate:"required"`
	SellerID string `json:"seller_id" validate:"required"`

	S3Path string `json:"s3_path"`
}

type BulkOrderCSVRow struct {
	SKUCode         string  `csv:"sku_code"`
	HubCode         string  `csv:"hub_code"`
	Quantity        int     `csv:"quantity"`
	UnitPrice       float64 `csv:"unit_price"`
	CustomerEmail   string  `csv:"customer_email"`
	CustomerFName   string  `csv:"customer_fname"`
	CustomerLName   string  `csv:"customer_lname"`
	CustomerPhone   string  `csv:"customer_phone"`
	ShippingStreet  string  `csv:"shipping_street"`
	ShippingCity    string  `csv:"shipping_city"`
	ShippingState   string  `csv:"shipping_state"`
	ShippingZip     string  `csv:"shipping_zip"`
	ShippingCountry string  `csv:"shipping_country"`
	ReferenceID     string  `csv:"reference_id"`
}

type InvalidCSVRow struct {
	RowNumber int       `json:"row_number" bson:"row_number"`
	RawData   string    `json:"raw_data" bson:"raw_data"`
	Reason    string    `json:"reason" bson:"reason"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`

	OrderID primitive.ObjectID `json:"order_id,omitempty" bson:"order_id,omitempty"`
}

type OrderListFilter struct {
	TenantID  string      `json:"tenant_id" form:"tenant_id"`
	SellerID  string      `json:"seller_id" form:"seller_id"`
	Status    OrderStatus `json:"status" form:"status"`
	StartDate *time.Time  `json:"start_date" form:"start_date"`
	EndDate   *time.Time  `json:"end_date" form:"end_date"`
	Page      int         `json:"page" form:"page"`
	Limit     int         `json:"limit" form:"limit"`
}

type OrderListResponse struct {
	Orders []Order `json:"orders"`
	Total  int64   `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

type CreateOrderRequest struct {
	TenantID     string       `json:"tenant_id" validate:"required"`
	SellerID     string       `json:"seller_id" validate:"required"`
	ReferenceID  string       `json:"reference_id" validate:"required"`
	Items        []OrderItem  `json:"items" validate:"required,min=1,dive"`
	CustomerInfo CustomerInfo `json:"customer_info" validate:"required"`
	ShippingInfo ShippingInfo `json:"shipping_info" validate:"required"`
	PaymentInfo  PaymentInfo  `json:"payment_info" validate:"required"`
	TotalAmount  float64      `json:"total_amount" validate:"required,min=0"`
	Currency     string       `json:"currency" validate:"required"`
}

type OrderCreatedEvent struct {
	OrderID     string      `json:"order_id"`
	TenantID    string      `json:"tenant_id"`
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"total_amount"`
	Timestamp   time.Time   `json:"timestamp"`
}

type OrderStatusUpdatedEvent struct {
	OrderID   string      `json:"order_id"`
	TenantID  string      `json:"tenant_id"`
	OldStatus OrderStatus `json:"old_status"`
	NewStatus OrderStatus `json:"new_status"`
	Timestamp time.Time   `json:"timestamp"`
}
