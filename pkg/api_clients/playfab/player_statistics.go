package playfab

type PlayerStatisticValueResponse struct {
	StatisticName string `json:"StatisticName"`
	Value         int    `json:"Value"`
	Version       int    `json:"Version"`
}
