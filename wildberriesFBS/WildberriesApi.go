package main

import (
	"io"
	"log"
	"net/http"
)

func GetOrdersBySupplyId(wildberriesKey, supplyId string) string {

	url := "https://marketplace-api.wildberries.ru/api/v3/supplies/" + supplyId + "/orders"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalf("Ошибка создания запроса: %v", err)
	}

	// Устанавливаем необходимые заголовки (если нужны)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", wildberriesKey) // Замените на ваш API токен, если нужен

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Ошибка: получен статус %d", resp.StatusCode)
	}

	// Читаем тело ответа
	jsonString, _ := io.ReadAll(resp.Body)

	// Выводим ответ
	return string(jsonString)
}
