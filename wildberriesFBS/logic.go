package wildberriesFBS

import (
	"encoding/base64"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func GetReadyFile(wildberriesKey, supplyId string) error {
	orders := GetOrders_FBS(wildberriesKey, supplyId)
	var ordersSlice []string
	for _, order := range orders {
		stickers := GetStickers_FBS(wildberriesKey, order.ID)
		decodeToPDF(stickers.Stickers[0].File, stickers.Stickers[0].OrderId, order.SKUs[0])
		ordersSlice = append(ordersSlice, "wildberriesFBS/ready/"+strconv.Itoa(order.ID)+".pdf")
	}
	err := mergePDFsInDirectory(ordersSlice, "wildberriesFBS/"+supplyId+".pdf")
	if err != nil {
		return err
	}
	if !fileExists("wildberriesFBS/" + supplyId + ".pdf") {
		err = fmt.Errorf("Такого файла не существует")
	}
	return err
}

func decodeToPNG(base64String string, orderId int) string {
	// Ваш base64 закодированный контент
	base64Data := base64String

	// Декодирование base64 в байты
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		fmt.Println("Ошибка при декодировании base64:", err)
	}

	// Определите путь и имя файла для сохранения
	filePath := "wildberriesFBS/codes/" + strconv.Itoa(orderId) + ".png" // Замените на желаемое имя файла и расширение

	// Открытие файла для записи
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
	}
	defer file.Close()

	// Запись данных в файл
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
	}

	return filePath
}

func decodeToPDF(base64String string, orderId int, sku string) {
	pageWidthMM := 75.0
	pageHeightMM := 120.0
	// Создание нового PDF-документа
	pdf := gofpdf.New("P", "mm", "", "")
	// Добавление страницы с заданными размерами
	pdf.AddPageFormat("P", gofpdf.SizeType{pageWidthMM, pageHeightMM})
	// Путь к первому PNG-файлу
	imgPath1 := decodeToPNG(base64String, orderId)
	// Добавление первого изображения в PDF (без изменения размера изображения)
	pdf.ImageOptions(imgPath1, (75-58)/2, 13, 58, 40, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	// Путь ко второму PNG-файлу
	imgPath2 := "wildberriesFBS/barcodes/" + sku + ".jpg" // Пример использования другого orderId для второго PNG
	// Добавление второго изображения в PDF на новую позицию
	if !fileExists(imgPath2) {
		imgPath2 = "wildberriesFBS/barcodes/0.jpg"
	}
	pdf.ImageOptions(imgPath2, (75-58)/2, 67, 58, 40, false, gofpdf.ImageOptions{ImageType: "JPG"}, 0, "")
	// Сохранение PDF-документа
	err := pdf.OutputFileAndClose("wildberriesFBS/ready/" + strconv.Itoa(orderId) + ".pdf")
	if err != nil {
		log.Fatalf("Error saving PDF: %s", err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir() // Проверяем, что это файл, а не директория
}

func mergePDFsInDirectory(orderSlice []string, outputFile string) error {

	// Проверяем, есть ли PDF-файлы для объединения
	if len(orderSlice) == 0 {
		return fmt.Errorf("нет PDF-файлов в директории")
	}

	// Формируем команду для выполнения merge PDF через pdfcpu
	args := append([]string{"merge", outputFile}, orderSlice...)
	cmd := exec.Command("pdfcpu", args...)

	// Запуск команды
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка выполнения pdfcpu: %v, %s", err, string(output))
	}
	return nil
}

func Clean_files(supplyId string) {
	err := os.RemoveAll("/Users/a123/GolandProjects/WildberriesGO/wildberriesFBS/codes")
	if err != nil {
		fmt.Println(err)
	}
	err = os.RemoveAll("/Users/a123/GolandProjects/WildberriesGO/wildberriesFBS/ready")
	if err != nil {
		fmt.Println(err)
	}
	err = os.Mkdir("/Users/a123/GolandProjects/WildberriesGO/wildberriesFBS/codes", 0755) // 0755 - это права доступа к директории (чтение, запись, выполнение)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Mkdir("/Users/a123/GolandProjects/WildberriesGO/wildberriesFBS/ready", 0755) // 0755 - это права доступа к директории (чтение, запись, выполнение)
	if err != nil {
		fmt.Println(err)
	}
	err = os.Remove("/Users/a123/GolandProjects/WildberriesGO/wildberriesFBS/" + supplyId + ".pdf")
	if err != nil {
		fmt.Println(err)
	}

}
