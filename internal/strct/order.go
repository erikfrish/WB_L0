package order

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Data struct {
	OrderUID    string `json:"order_uid,omitempty"`
	TrackNumber string `json:"track_number,omitempty"`
	Entry       string `json:"entry,omitempty"`
	Delivery    struct {
		Name    string `json:"name,omitempty"`
		Phone   string `json:"phone,omitempty"`
		Zip     string `json:"zip,omitempty"`
		City    string `json:"city,omitempty"`
		Address string `json:"address,omitempty"`
		Region  string `json:"region,omitempty"`
		Email   string `json:"email,omitempty"`
	} `json:"delivery,omitempty"`
	Payment struct {
		Transaction  string `json:"transaction,omitempty"`
		RequestID    string `json:"request_id,omitempty"`
		Currency     string `json:"currency,omitempty"`
		Provider     string `json:"provider,omitempty"`
		Amount       int    `json:"amount,omitempty"`
		PaymentDt    int    `json:"payment_dt,omitempty"`
		Bank         string `json:"bank,omitempty"`
		DeliveryCost int    `json:"delivery_cost,omitempty"`
		GoodsTotal   int    `json:"goods_total,omitempty"`
		CustomFee    int    `json:"custom_fee,omitempty"`
	} `json:"payment,omitempty"`
	Items []struct {
		ChrtID      int    `json:"chrt_id,omitempty"`
		TrackNumber string `json:"track_number,omitempty"`
		Price       int    `json:"price,omitempty"`
		Rid         string `json:"rid,omitempty"`
		Name        string `json:"name,omitempty"`
		Sale        int    `json:"sale,omitempty"`
		Size        string `json:"size,omitempty"`
		TotalPrice  int    `json:"total_price,omitempty"`
		NmID        int    `json:"nm_id,omitempty"`
		Brand       string `json:"brand,omitempty"`
		Status      int    `json:"status,omitempty"`
	} `json:"items,omitempty"`
	Locale            string    `json:"locale,omitempty"`
	InternalSignature string    `json:"internal_signature,omitempty"`
	CustomerID        string    `json:"customer_id,omitempty"`
	DeliveryService   string    `json:"delivery_service,omitempty"`
	Shardkey          string    `json:"shardkey,omitempty"`
	SmID              int       `json:"sm_id,omitempty"`
	DateCreated       time.Time `json:"date_created,omitempty"`
	OofShard          string    `json:"oof_shard,omitempty"`
}

func (d Data) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Data) Scan(value any) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &d)
}

// type MultipleData []Data

// func (m *MultipleData) Scan(values ...any) error {
// 	var err error
// 	for i, value := range values {
// 		b, ok := value.([]byte)
// 		if !ok {
// 			return errors.New("type assertion to []byte failed")
// 		}
// 		m1 := make(MultipleData, len(values))
// 		err = json.Unmarshal(b, &m1[i])
// 	}
// 	return err
// }
