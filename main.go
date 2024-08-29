package main

import (
	"firstProject/wildberriesFBS/WildberriesApiToEntity"
	"fmt"
)

func main() {
	wildberriesKey := ""

	supplyId := "WB-GI-104193936"

	json := WildberriesApiToEntity.GetOrders_FBS(wildberriesKey, supplyId)

	// Выводим результат.
	for _, order := range json.Orders {
		fmt.Printf("Order UID: %s, Price: %d\n", order.ID, order.Article)
	}
}
