package main

import (
	"encoding/json"
	"log"
)

func GetOrders_FBS(wildberriesKey, supplyId string) Orders {
	jsonString := GetOrdersBySupplyId(wildberriesKey, supplyId)
	var orders Orders

	// Декодируем JSON в структуру.
	err := json.Unmarshal([]byte(jsonString), &orders)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	return orders
}
