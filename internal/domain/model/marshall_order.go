package model

import (
	"encoding/json"
	"time"
)

func MarshalOrder(order *Order) ([]byte, error) {
	// Вспомогательная структура для кастомного маршалинга
	aux := struct {
		OrderUID    string `json:"order_uid"`
		TrackNumber string `json:"track_number"`
		Entry       string `json:"entry"`
		Delivery    struct {
			Name    string `json:"name"`
			Phone   string `json:"phone"`
			Zip     string `json:"zip"`
			City    string `json:"city"`
			Address string `json:"address"`
			Region  string `json:"region"`
			Email   string `json:"email"`
		} `json:"delivery"`
		Payment struct {
			Transaction  string `json:"transaction"`
			RequestID    string `json:"request_id"`
			Currency     string `json:"currency"`
			Provider     string `json:"provider"`
			Amount       int    `json:"amount"`
			PaymentDT    int64  `json:"payment_dt"`
			Bank         string `json:"bank"`
			DeliveryCost int    `json:"delivery_cost"`
			GoodsTotal   int    `json:"goods_total"`
			CustomFee    int    `json:"custom_fee"`
		} `json:"payment"`
		Items []struct {
			ChrtID      int    `json:"chrt_id"`
			TrackNumber string `json:"track_number"`
			Price       int    `json:"price"`
			Rid         string `json:"rid"`
			Name        string `json:"name"`
			Sale        int    `json:"sale"`
			Size        string `json:"size"`
			TotalPrice  int    `json:"total_price"`
			NmID        int    `json:"nm_id"`
			Brand       string `json:"brand"`
			Status      int    `json:"status"`
		} `json:"items"`
		Locale            string `json:"locale"`
		InternalSignature string `json:"internal_signature"`
		CustomerID        string `json:"customer_id"`
		DeliveryService   string `json:"delivery_service"`
		ShardKey          string `json:"shardkey"`
		SmID              int    `json:"sm_id"`
		DateCreated       string `json:"date_created"`
		OofShard          string `json:"oof_shard"`
	}{
		OrderUID:          order.OrderUID.String(),
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		ShardKey:          order.ShardKey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated.Format(time.RFC3339),
		OofShard:          order.OofShard,
	}

	// Заполняем Delivery
	aux.Delivery.Name = order.Delivery.Name
	aux.Delivery.Phone = order.Delivery.Phone
	aux.Delivery.Zip = order.Delivery.Zip
	aux.Delivery.City = order.Delivery.City
	aux.Delivery.Address = order.Delivery.Address
	aux.Delivery.Region = order.Delivery.Region
	aux.Delivery.Email = order.Delivery.Email

	// Заполняем Payment
	aux.Payment.Transaction = order.Payment.Transaction.String()
	aux.Payment.RequestID = order.Payment.RequestID
	aux.Payment.Currency = order.Payment.Currency
	aux.Payment.Provider = order.Payment.Provider
	aux.Payment.Amount = order.Payment.Amount
	aux.Payment.PaymentDT = order.Payment.PaymentDT.Unix()
	aux.Payment.Bank = order.Payment.Bank
	aux.Payment.DeliveryCost = order.Payment.DeliveryCost
	aux.Payment.GoodsTotal = order.Payment.GoodsTotal
	aux.Payment.CustomFee = order.Payment.CustomFee

	// Заполняем Items
	for _, item := range order.Items {
		aux.Items = append(aux.Items, struct {
			ChrtID      int    `json:"chrt_id"`
			TrackNumber string `json:"track_number"`
			Price       int    `json:"price"`
			Rid         string `json:"rid"`
			Name        string `json:"name"`
			Sale        int    `json:"sale"`
			Size        string `json:"size"`
			TotalPrice  int    `json:"total_price"`
			NmID        int    `json:"nm_id"`
			Brand       string `json:"brand"`
			Status      int    `json:"status"`
		}{
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			Rid:         item.Rid,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.NmID,
			Brand:       item.Brand,
			Status:      item.Status,
		})
	}

	return json.Marshal(aux)
}
