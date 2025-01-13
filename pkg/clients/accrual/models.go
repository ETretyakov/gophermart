package accrual

type GoodsCreate struct {
	Match      string `json:"match"`
	Reward     int    `json:"reward"`
	RewardType string `json:"reward_type"`
}

type Goods struct {
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type OrderCreate struct {
	Order string  `json:"order"`
	Goods []Goods `json:"goods"`
}

type OrderRead struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
