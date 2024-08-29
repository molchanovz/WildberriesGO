package main

import (
	"fmt"
)

/*ну вот логика такая:
берем все заказы из поставки, закидывает в массив
по этому массиву проходим, для каждого заказа получаем куаркод в svg
перекидываем все svg в png или в pdf
находим баркод у каждого заказа, соединяем его с куаркодом, чтобы они были все на одном стикере*/

func main() {
	wildberriesKey := ""

	supplyId := "WB-GI-104193936"

	json := GetOrders_FBS(wildberriesKey, supplyId)

	// Выводим результат.
	for _, order := range json.Orders {
		fmt.Printf("Order UID: %s, Price: %d\n", order.ID, order.Article)
	}
}
