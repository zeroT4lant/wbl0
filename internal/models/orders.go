package models

import (
	"time"
)

type Order struct {
	OrderUid          string    `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Delivery          Delivery  `json:"delivery"`
	Payment           Payment   `json:"payment"`
	Items             []Item    `json:"items"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          string    `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}

func (o Order) CheckOrder() bool {
	if len(o.OrderUid) != 19 {
		return false
	}

	if len(o.TrackNumber) != 14 {
		return false
	}

	if len(o.Entry) != 4 {
		return false
	}

	//_, err := mail.ParseAddress(o.Delivery.Email)
	//if err != nil {
	//	return false
	//}

	//regx := regexp.MustCompile(`\+9\d{9}`)
	//res := regx.MatchString(o.Delivery.Phone)
	//
	//return res

	return true

}
