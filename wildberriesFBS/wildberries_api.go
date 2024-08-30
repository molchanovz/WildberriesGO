package wildberriesFBS

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

func GetCodesByOrderId(wildberriesKey string, orderId int) string {

	params := url.Values{}
	params.Add("type", "png")
	params.Add("width", "58")
	params.Add("height", "40")

	// Основной URL
	baseURL := "https://marketplace-api.wildberries.ru/api/v3/orders/stickers"

	// Добавление параметров к URL
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	type RequestBody struct {
		Orders []int `json:"orders"`
	}

	data := RequestBody{
		Orders: []int{orderId},
	}

	// Преобразование данных в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка при преобразовании данных в JSON:", err)
	}

	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonData))

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
