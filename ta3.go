package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"net/http"

)

type House struct {
	Size  float64
	Rooms float64
	Price float64
}

func readCSV(filename string) ([]House, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	houses := make([]House, len(records)-1)
	for i, record := range records[1:] {
		size, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, err
		}
		rooms, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		houses[i] = House{Size: size, Rooms: rooms, Price: price}
	}

	return houses, nil
}

func generateTestData(num int) []House {
	rand.Seed(time.Now().UnixNano())
	testData := make([]House, num)
	for i := 0; i < num; i++ {
		size := 120 + rand.Float64()*(500-120)
		rooms := 4 + rand.Float64()*(13-4)
		price := size*1200 + rooms*500
		testData[i] = House{Size: size, Rooms: rooms, Price: price}
	}
	return testData
}

func mean(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func calculateCoefficients(trainData []House) (float64, float64, float64) {
	size := len(trainData)
	covarianceSizePrice := make(chan float64)
	covarianceRoomsPrice := make(chan float64)
	varianceSize := make(chan float64)
	varianceRooms := make(chan float64)

	var wg sync.WaitGroup

	sizes := make([]float64, size)
	rooms := make([]float64, size)
	prices := make([]float64, size)

	for i, house := range trainData {
		sizes[i] = house.Size
		rooms[i] = house.Rooms
		prices[i] = house.Price
	}

	meanSize := mean(sizes)
	meanRooms := mean(rooms)
	meanPrice := mean(prices)

	wg.Add(4)
	go func() {
		defer wg.Done()
		var sum float64
		for i := 0; i < size; i++ {
			sum += (sizes[i] - meanSize) * (prices[i] - meanPrice)
		}
		covarianceSizePrice <- sum / float64(size)
	}()

	go func() {
		defer wg.Done()
		var sum float64
		for i := 0; i < size; i++ {
			sum += (rooms[i] - meanRooms) * (prices[i] - meanPrice)
		}
		covarianceRoomsPrice <- sum / float64(size)
	}()

	go func() {
		defer wg.Done()
		var sum float64
		for i := 0; i < size; i++ {
			sum += math.Pow(sizes[i]-meanSize, 2)
		}
		varianceSize <- sum / float64(size)
	}()

	go func() {
		defer wg.Done()
		var sum float64
		for i := 0; i < size; i++ {
			sum += math.Pow(rooms[i]-meanRooms, 2)
		}
		varianceRooms <- sum / float64(size)
	}()

	go func() {
		wg.Wait()
		close(covarianceSizePrice)
		close(covarianceRoomsPrice)
		close(varianceSize)
		close(varianceRooms)
	}()

	covSP := <-covarianceSizePrice
	covRP := <-covarianceRoomsPrice
	varS := <-varianceSize
	varR := <-varianceRooms

	beta1 := covSP / varS
	beta2 := covRP / varR
	beta0 := meanPrice - beta1*meanSize - beta2*meanRooms

	return beta0, beta1, beta2
}

func predict(size, rooms, beta0, beta1, beta2 float64) float64 {
	return beta0 + beta1*size + beta2*rooms
}

func main() {
	//houses, err := readCSV("house_prices.csv")
	//https://github.com/parzzd/ta3/blob/main/house_prices.csv
	//https://raw.githubusercontent.com/parzzd/ta3/main/house_prices.csv



	

	/*
	// Realizar la solicitud HTTP para obtener el contenido del archivo CSV
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error al obtener el archivo CSV:", err)
		return
	}
	defer resp.Body.Close()


	// Verificar el código de estado de la respuesta HTTP
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error al obtener el archivo CSV: código de estado", resp.StatusCode)
		return
	}

	//------------------------------------------------------------------------------

	houses, err := readCSV("house_prices.csv")

	if err != nil {
		log.Fatal(err)
	}

*/


	//-----------------------------------
	//-----------------------------------
	//-----------------------------------

	//LINEAS AGREGADAS


	url := "https://raw.githubusercontent.com/parzzd/ta3/main/house_prices.csv"

	// Realizar la solicitud HTTP para obtener el contenido del archivo CSV
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error al obtener el archivo CSV:", err)
	}
	defer resp.Body.Close()


	// Lee CSV
	reader := csv.NewReader(resp.Body)

	// Lee todas las filas del CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error al leer el archivo CSV:", err)
	}

	var houses []House
    for _, record := range records {
        size,_ := strconv.ParseFloat(record[0], 64)
        
        rooms,_ := strconv.ParseFloat(record[1], 64)
        
        price,_ := strconv.ParseFloat(record[2], 64)
    

        houses = append(houses, House{
            Size:  size,
            Rooms: rooms,
            Price: price,
        })
    }





//-----------------------------------
//-----------------------------------
//-----------------------------------

	trainSize := len(houses)
	trainData := houses[:trainSize ]

	testData := generateTestData(1000)

	beta0, beta1, beta2 := calculateCoefficients(trainData)

	fmt.Printf("Beta0: %.2f, Beta1: %.2f, Beta2: %.2f\n", beta0, beta1, beta2)

	var totalError float64
	for _, house := range testData {
		predictedPrice := predict(house.Size, house.Rooms, beta0, beta1, beta2)
		totalError += math.Pow(predictedPrice-house.Price, 2)
	}

	mse := totalError / float64(len(testData))
	fmt.Printf("MSE: %.2f\n", mse)
}
