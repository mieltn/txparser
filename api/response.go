package api

type CurrentBlockResponse struct {
	BlockNumber uint64 `json:"block_number"`
}

type SubscribeResponse struct {
	Ok bool `json:"ok"`
}

type Transaction struct {
	Hash   string `json:"hash"`
	Amount uint64 `json:"amount"`
	From   string `json:"from"`
	To     string `json:"to"`
}

type GetTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

type Error struct {
	Message string `json:"message"`
}
