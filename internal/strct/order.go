package order

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.2 --name=Value
type Data struct {
	OrderUID    string `json:"order_uid" validate:"required"`    // "b563feb7b2b84b6test" // [a-z0-9]{18,20}
	TrackNumber string `json:"track_number" validate:"required"` // "WBILMTESTTRACK" // [A-Z]{12-14}
	Entry       string `json:"entry,omitempty"`                  // "WBIL"
	Delivery    struct {
		Name    string `json:"name,omitempty" validate:"alpha"`  // "Test Testov"
		Phone   string `json:"phone,omitempty" validate:"e164"`  // "+9720000000"
		Zip     string `json:"zip,omitempty"`                    // "2639809"
		City    string `json:"city,omitempty"`                   // "Kiryat Mozkin"
		Address string `json:"address,omitempty"`                // Ploshad Mira 15"
		Region  string `json:"region,omitempty"`                 // "Kraiot"
		Email   string `json:"email,omitempty" validate:"email"` // "test@gmail.com"
	} `json:"delivery,omitempty"`
	Payment struct {
		Transaction  string `json:"transaction" validate:"required"`       // "b563feb7b2b84b6test"
		RequestID    string `json:"request_id,omitempty"`                  // ""
		Currency     string `json:"currency"  validate:"required,iso4217"` // "USD"
		Provider     string `json:"provider,omitempty"`                    // "wbpay"
		Amount       int    `json:"amount" validate:"required,gte=0"`      // 1817
		PaymentDt    int    `json:"payment_dt,omitempty"`                  // 1637907727
		Bank         string `json:"bank,omitempty" validate:"alpha"`       // "alpha"
		DeliveryCost int    `json:"delivery_cost,omitempty"`               // 1500
		GoodsTotal   int    `json:"goods_total,omitempty"`                 // 317
		CustomFee    int    `json:"custom_fee,omitempty"`                  // 0
	} `json:"payment,omitempty"`
	Items []struct {
		ChrtID      int    `json:"chrt_id,omitempty"`                             // 9934930
		TrackNumber string `json:"track_number,omitempty"`                        // "WBILMTESTTRACK"
		Price       int    `json:"price,omitempty" validate:"numeric,gte=0"`      // 453
		Rid         string `json:"rid,omitempty"`                                 // "ab4219087a764ae0btest"
		Name        string `json:"name,omitempty"`                                // "Mascaras"
		Sale        int    `json:"sale,omitempty"`                                // 30
		Size        string `json:"size,omitempty"`                                // "0"
		TotalPrice  int    `json:"total_price,omitempty" validate:"numeric,gt=0"` //  317
		NmID        int    `json:"nm_id,omitempty"`                               // 2389212
		Brand       string `json:"brand,omitempty"`                               // "Vivienne Sabo"
		Status      int    `json:"status,omitempty"`                              // 202
	} `json:"items,omitempty"`
	Locale            string    `json:"locale,omitempty"`             // "en"
	InternalSignature string    `json:"internal_signature,omitempty"` // ""
	CustomerID        string    `json:"customer_id,omitempty"`        // "test"
	DeliveryService   string    `json:"delivery_service,omitempty"`   // "meest"
	Shardkey          string    `json:"shardkey,omitempty"`           // "9"
	SmID              int       `json:"sm_id,omitempty"`              // 99
	DateCreated       time.Time `json:"date_created,omitempty"`       // "2021-11-26T06:22:19Z"
	OofShard          string    `json:"oof_shard,omitempty"`          // "1"
}

func (d Data) Value() (driver.Value, error) {
	v, err := json.Marshal(d)
	if err != nil {
		err = errors.Join(errors.New("failed to get data Value"), err)
	}
	return v, err
}

func (d *Data) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	if err := json.Unmarshal(b, &d); err != nil {
		err = errors.Join(errors.New("failed to Scan data to struct"), err)
		return err
	}
	if err := d.Validate(); err != nil {
		err = errors.Join(errors.New("failed to Validate data to struct order.Data"), err)
		return err
	}
	return nil
}

func (d *Data) Validate() error {
	validate := validator.New()
	return validate.Struct(d)
}

// For now it validates not all fields, just:
// OrderUID, TrackNumber, Delivery.Name, Delivery.Phone, Delivery.Email, Payment.Transaction,
// Payment.Currency, Payment.Amount, Payment.Bank, Items.Price, Items.TotalPrice
// If needed, we can expand it in future
