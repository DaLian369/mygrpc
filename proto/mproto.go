package proto

const (
	ExchangeTopic = "exchange"
)

type Order struct {
	Id        int64  `json:"order"`
	From      int64  `json:"from"`
	FromMoney int64  `json:"from_money"`
	To        int64  `json:"to"`
	ToMoney   int64  `json:"to_money"`
	Ext       string `json:"ext"`
	Ct        int64  `json:"ct"`
}
