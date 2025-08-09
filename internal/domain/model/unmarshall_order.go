package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func UnmarshalOrder(data []byte) (*Order, error) {
	// Вспомогательная структура для кастомного парсинга
	type Alias struct {
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
	}

	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	// Парсинг UUID
	orderUID, err := uuid.Parse(aux.OrderUID)
	if err != nil {
		return nil, fmt.Errorf("invalid order_uid: %w", err)
	}

	paymentTransaction, err := uuid.Parse(aux.Payment.Transaction)
	if err != nil {
		return nil, fmt.Errorf("invalid payment transaction: %w", err)
	}

	// Парсинг дат
	dateCreated, err := time.Parse(time.RFC3339, aux.DateCreated)
	if err != nil {
		return nil, fmt.Errorf("invalid date_created format: %w", err)
	}

	paymentDT := time.Unix(aux.Payment.PaymentDT, 0)

	// Формируем итоговую структуру
	order := &Order{
		OrderUID:          orderUID,
		TrackNumber:       aux.TrackNumber,
		Entry:             aux.Entry,
		Locale:            aux.Locale,
		InternalSignature: aux.InternalSignature,
		CustomerID:        aux.CustomerID,
		DeliveryService:   aux.DeliveryService,
		ShardKey:          aux.ShardKey,
		SmID:              aux.SmID,
		DateCreated:       dateCreated,
		OofShard:          aux.OofShard,
		Delivery: Delivery{
			OrderUID: orderUID,
			Name:     aux.Delivery.Name,
			Phone:    aux.Delivery.Phone,
			Zip:      aux.Delivery.Zip,
			City:     aux.Delivery.City,
			Address:  aux.Delivery.Address,
			Region:   aux.Delivery.Region,
			Email:    aux.Delivery.Email,
		},
		Payment: Payment{
			Transaction:  paymentTransaction,
			RequestID:    aux.Payment.RequestID,
			Currency:     aux.Payment.Currency,
			Provider:     aux.Payment.Provider,
			Amount:       aux.Payment.Amount,
			PaymentDT:    paymentDT,
			Bank:         aux.Payment.Bank,
			DeliveryCost: aux.Payment.DeliveryCost,
			GoodsTotal:   aux.Payment.GoodsTotal,
			CustomFee:    aux.Payment.CustomFee,
		},
	}

	// Добавляем items
	for _, item := range aux.Items {
		order.Items = append(order.Items, Item{
			OrderUID:    orderUID,
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

	return order, nil
}
