package entities

type OrderWB struct {
	OrderUID              string   `json:"orderUid"`
	Article               string   `json:"article"`
	ColorCode             string   `json:"colorCode"`
	RID                   string   `json:"rid"`
	CreatedAt             string   `json:"createdAt"` // Можно использовать time.Time, если хотите парсить даты
	Offices               []string `json:"offices"`
	SKUs                  []string `json:"skus"`
	ID                    int      `json:"id"`
	WarehouseID           int      `json:"warehouseId"`
	NmID                  int      `json:"nmId"`
	ChrtID                int      `json:"chrtId"`
	Price                 int      `json:"price"`
	ConvertedPrice        int      `json:"convertedPrice"`
	CurrencyCode          int      `json:"currencyCode"`
	ConvertedCurrencyCode int      `json:"convertedCurrencyCode"`
	CargoType             int      `json:"cargoType"`
	IsZeroOrder           bool     `json:"isZeroOrder"`
}
type Orders struct {
	Orders []OrderWB `json:"orders"`
}
