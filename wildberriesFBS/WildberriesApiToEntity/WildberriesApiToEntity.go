package WildberriesApiToEntity

import (
	"encoding/json"
	"firstProject/wildberriesFBS/WildberriesApi"
	"firstProject/wildberriesFBS/entities"
	"log"
)

func GetOrders_FBS(wildberriesKey, supplyId string) entities.Orders {
	jsonString := WildberriesApi.GetOrdersBySupplyId(wildberriesKey, supplyId)
	var orders entities.Orders

	// Декодируем JSON в структуру.
	err := json.Unmarshal([]byte(jsonString), &orders)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	return orders
}
