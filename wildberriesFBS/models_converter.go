package wildberriesFBS

import (
	"encoding/json"
	"log"
	"sort"
)

func GetOrders_FBS(wildberriesKey, supplyId string) []OrderWB {
	jsonString := GetOrdersBySupplyId(wildberriesKey, supplyId)
	var orders Orders
	err := json.Unmarshal([]byte(jsonString), &orders)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}
	sortOrdersByArticle(orders.Orders)
	return orders.Orders
}

func GetStickers_FBS(wildberriesKey string, orderId int) StickerWB {
	jsonString := GetCodesByOrderId(wildberriesKey, orderId)
	var stickers StickerWB
	err := json.Unmarshal([]byte(jsonString), &stickers)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}
	return stickers
}

func sortOrdersByArticle(orders []OrderWB) {
	sort.SliceStable(orders, func(i, j int) bool {
		return orders[i].Article < orders[j].Article
	})
}
