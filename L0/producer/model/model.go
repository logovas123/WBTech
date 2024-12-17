package model

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
	ShardKey          string   `json:"shardkey"`
	SmID              int64    `json:"sm_id" fake:"{number:0,100}"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	ZIP     string `json:"zip"`
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
	Amount       int64  `json:"amount" fake:"{number:0,100000}"`
	PaymentDT    int64  `json:"payment_dt" fake:"{number:0,100000}"`
	Bank         string `json:"bank"`
	DeliveryCost int64  `json:"delivery_cost" fake:"{number:0,1000000}"`
	GoodsTotal   int64  `json:"goods_total" fake:"{number:0,1000000}"`
	CustomFee    int64  `json:"custom_fee" fake:"{number:0,100000}"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id" fake:"{number:0,100000}"`
	TrackNumber string `json:"track_number"`
	Price       int64  `json:"price" fake:"{number:0,1000000}"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int64  `json:"sale" fake:"{number:0,1000000}"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price" fake:"{number:0,1000000}"`
	NmID        int64  `json:"nm_id" fake:"{number:0,1000000}"`
	Brand       string `json:"brand"`
	Status      int64  `json:"status" fake:"{number:0,1000000}"`
}
