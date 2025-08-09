package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
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
}

type Item struct {
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
}

type Order struct {
	OrderUID          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Item   `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	Shardkey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

func main() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:29092"
	}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "order_created"
	}
	intervalSec := 60
	if s := os.Getenv("GENERATOR_INTERVAL_SEC"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			intervalSec = v
		}
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	rand.Seed(time.Now().UnixNano())
	log.Printf("order-generator started: broker=%s topic=%s interval=%ds\n", broker, topic, intervalSec)

	// handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("shutdown signal received, exiting...")
		cancel()
	}()

	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer ticker.Stop()

	// send one immediately, then every tick
	sendOrder(ctx, writer)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sendOrder(ctx, writer)
		}
	}
}

func sendOrder(ctx context.Context, writer *kafka.Writer) {
	order := generateOrder()

	data, err := json.Marshal(order)
	if err != nil {
		log.Printf("marshal error: %v\n", err)
		return
	}

	msg := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: data,
	}

	// try to write with a timeout context to avoid hanging forever
	writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := writer.WriteMessages(writeCtx, msg); err != nil {
		log.Printf("failed to write message: %v\n", err)
		return
	}

	log.Printf("Sent order order_uid=%s\n", order.OrderUID)
}

func generateOrder() Order {
	orderUID := uuid.New().String()
	track := "WBILM" + randomDigits(6) + "TRACK"
	entry := "WBIL"

	items := make([]Item, 1)
	price := rand.Intn(900) + 100
	amount := 1
	items[0] = Item{
		ChrtID:      rand.Intn(9_999_999),
		TrackNumber: track,
		Price:       price,
		Rid:         uuid.New().String(),
		Name:        randomProductName(),
		Sale:        rand.Intn(50),
		Size:        "0",
		TotalPrice:  price * amount,
		NmID:        rand.Intn(9_999_999),
		Brand:       randomBrand(),
		Status:      200 + rand.Intn(10),
	}

	payment := Payment{
		Transaction:  orderUID,
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       items[0].TotalPrice,
		PaymentDT:    time.Now().Unix(),
		Bank:         randomBank(),
		DeliveryCost: 1500,
		GoodsTotal:   items[0].TotalPrice,
		CustomFee:    0,
	}

	del := Delivery{
		Name:    "Test Testov",
		Phone:   "+972" + randomDigits(7),
		Zip:     randomDigits(7),
		City:    randomCity(),
		Address: "Ploshad Mira " + strconv.Itoa(rand.Intn(100)),
		Region:  "Kraiot",
		Email:   "test+" + randomDigits(3) + "@gmail.com",
	}

	return Order{
		OrderUID:          orderUID,
		TrackNumber:       track,
		Entry:             entry,
		Delivery:          del,
		Payment:           payment,
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test_customer_" + randomDigits(3),
		DeliveryService:   "meest",
		Shardkey:          strconv.Itoa(rand.Intn(10)),
		SmID:              rand.Intn(200),
		DateCreated:       time.Now().UTC().Format(time.RFC3339),
		OofShard:          "1",
	}
}

func randomDigits(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('0' + rand.Intn(10))
	}
	return string(b)
}

func randomProductName() string {
	names := []string{"Mascaras", "Laptop", "Phone", "Headphones", "Keyboard", "Monitor", "Watch"}
	return names[rand.Intn(len(names))]
}

func randomBrand() string {
	brands := []string{"Vivienne Sabo", "BrandX", "BrandY", "Acme"}
	return brands[rand.Intn(len(brands))]
}

func randomBank() string {
	banks := []string{"alpha", "bankA", "bankB"}
	return banks[rand.Intn(len(banks))]
}

func randomCity() string {
	cities := []string{"Kiryat Mozkin", "Tel Aviv", "Haifa", "Jerusalem"}
	return cities[rand.Intn(len(cities))]
}
