package marketdata

type GetMarketDataPayload struct {
	Symbol string `form:"symbol" valid:"optional"`
}

type AggTrade struct {
	EventType        string `json:"e"` // Event type
	EventTime        int64  `json:"E"` // Event time
	AggregateTradeID int64  `json:"a"` // Aggregate trade ID
	Symbol           string `json:"s"` // Symbol
	Price            string `json:"p"` // Price
	Quantity         string `json:"q"` // Quantity
	FirstTradeID     int64  `json:"f"` // First trade ID
	LastTradeID      int64  `json:"l"` // Last trade ID
	TradeTime        int64  `json:"T"` // Trade time
	IsMarketMaker    bool   `json:"m"` // Is the buyer the market maker
}

type MarkPriceStream struct {
	Stream string        `json:"stream"` // Stream name
	Data   MarkPriceData `json:"data"`   // Data details
}

type MarkPriceData struct {
	EventType            string `json:"e"` // Event type
	EventTime            int64  `json:"E"` // Event time
	Symbol               string `json:"s"` // Symbol
	MarkPrice            string `json:"p"` // Current mark price
	IndexPrice           string `json:"P"` // Current index price
	EstimatedSettlePrice string `json:"i"` // Estimated settle price
	FundingRate          string `json:"r"` // Funding rate
	NextFundingTime      int64  `json:"T"` // Next funding time
}

type BookTickerMessage struct {
	Stream string         `json:"stream"`
	Data   BookTickerData `json:"data"`
}

type BookTickerData struct {
	EventType       string `json:"e"` // Event type, e.g., "bookTicker"
	UpdateID        int64  `json:"u"` // Update ID
	Symbol          string `json:"s"` // Symbol, e.g., "DOTUSDT"
	BidPrice        string `json:"b"` // Best bid price
	BidQuantity     string `json:"B"` // Best bid quantity
	AskPrice        string `json:"a"` // Best ask price
	AskQuantity     string `json:"A"` // Best ask quantity
	TransactionTime int64  `json:"T"` // Transaction time
	EventTime       int64  `json:"E"` // Event time
}
