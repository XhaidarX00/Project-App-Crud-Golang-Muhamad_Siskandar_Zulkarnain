package view

import (
	model "ProjectCRUD/Model"
	utils "ProjectCRUD/Utils"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func PrintStudentToJson(data []model.Student) {
	utils.ClearScreen()
	fmt.Println(utils.ColorMessage("yellow", "\n============= Data Json Students ============="))
	time.Sleep(1 * time.Second)
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func PrintCourseToJson(data []model.Course) {
	utils.ClearScreen()
	fmt.Println(utils.ColorMessage("yellow", "\n============= Data Json Course ============="))
	time.Sleep(1 * time.Second)
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(jsonData))
}

func DisplayDataChosenIndexOne(title string, data []map[string]interface{}, keys []string) string {
	var id string

	for {
		msg := fmt.Sprintf("\n============== ◉  %s ◉  ==============", title)
		utils.PrintColorMsg("yellow", msg)
		fmt.Println(strings.Repeat("-", 50))

		for i, result := range data {
			fmt.Printf("%d. ", i+1)
			for _, key := range keys {
				value, ok := result[key]
				if !ok {
					value = "N/A"
				}
				fmt.Printf("%v | ", value)
			}
			fmt.Println()
		}

		fmt.Println(strings.Repeat("-", 50))
		fmt.Print(utils.ColorMessage("yellow", "Masukkan nomor siswa yang akan diupdate: "))
		fmt.Scan(&id)
		utils.ClearScreen()

		intInput, err := strconv.Atoi(id)
		if err != nil || intInput < 1 || intInput > len(data) {
			msg := "Input harus berupa angka yang valid dan harus sesuai dengan nomor daftar siswa."
			utils.ErrorMessage(msg)
			continue
		}

		if selectedData, ok := data[intInput-1]["ID"].(string); ok {
			return selectedData
		} else {
			utils.ErrorMessage("ID siswa tidak valid.")
			continue
		}
	}
}

// func DisplayData(title string, data []map[string]interface{}, keys []string) {
// 	msg := fmt.Sprintf("\n============== ◉  %s ◉  ==============", title)
// 	utils.PrintColorMsg("yellow", msg)
// 	fmt.Println(strings.Repeat("-", 50))
// 	for i, result := range data {
// 		fmt.Printf("%d. ", i+1)
// 		for _, key := range keys {
// 			value, ok := result[key]
// 			if !ok {
// 				value = "N/A"
// 			}
// 			fmt.Printf("%v | ", value)
// 		}
// 		fmt.Println()
// 	}
// 	fmt.Println(strings.Repeat("-", 50))
// }

func DisplayData(title string, data []map[string]interface{}, keys []string) {
	msg := fmt.Sprintf("\n============== ◉  %s ◉  ==============", title)
	utils.PrintColorMsg("yellow", msg)
	fmt.Println(strings.Repeat("-", 50))
	for i, result := range data {
		fmt.Printf("%d. ", i+1)
		for _, key := range keys {
			value, ok := result[key]
			if !ok {
				value = "N/A"
			}
			// Using reflect to get the type and value of the data
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.String:
				fmt.Printf("%s | ", val.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fmt.Printf("%d | ", val.Int())
			case reflect.Float32, reflect.Float64:
				fmt.Printf("%f | ", val.Float())
			case reflect.Bool:
				fmt.Printf("%t | ", val.Bool())
			default:
				fmt.Printf("%v | ", val.Interface())
			}
		}
		fmt.Println()
	}
	fmt.Println(strings.Repeat("-", 50))
}

func DisplayTempChoice(tempChoice []string) {
	for _, msg := range tempChoice {
		fmt.Println(msg)
	}
}

func DisplayChoise(title string, data []map[string]interface{}, keys []string) []map[string]interface{} {
	var input string
	var tempChoice []string
	var selectedData []map[string]interface{}
	defer func() {
		tempChoice = []string{}
	}()

	for {
		DisplayData(title, data, keys)
		fmt.Print(utils.ColorMessage("yellow", "Masukkan Nomor atau ketik 'done' untuk selesai: "))
		fmt.Scan(&input)
		utils.ClearScreen()

		if strings.ToLower(input) == "done" {
			return selectedData
		}

		intInput, err := strconv.Atoi(input)
		if err != nil || intInput < 1 || intInput > len(data) {
			msg := fmt.Sprintf("Input harus berupa angka yang valid dan tidak boleh lebih dari %d\n", len(data))
			utils.ErrorMessage(msg)
			continue
		}

		result := data[intInput-1]
		selectedData = append(selectedData, result)

		if value, ok := result[keys[0]].(string); ok {
			msg := fmt.Sprintf("Anda memilih menu: %s", utils.ColorMessage("green", value))
			tempChoice = append(tempChoice, msg)
		} else {
			msg := fmt.Sprintf("Anda memilih menu: %s", utils.ColorMessage("green", fmt.Sprintf("%v", result[keys[0]])))
			tempChoice = append(tempChoice, msg)
		}

		DisplayTempChoice(tempChoice)
	}
}

func DisplayMenu() {
	fmt.Println(utils.ColorMessage("yellow", "\n=---------- Menu ----------="))
	fmt.Println("1. " + utils.ColorMessage("green", "Tambah Course"))
	fmt.Println("2. " + utils.ColorMessage("green", "Lihat Semua Course"))
	fmt.Println("3. " + utils.ColorMessage("green", "Tambah Siswa"))
	fmt.Println("4. " + utils.ColorMessage("green", "Lihat Semua Siswa"))
	fmt.Println("5. " + utils.ColorMessage("green", "Update Data Siswa"))
	fmt.Println("6. " + utils.ColorMessage("green", "Hapus Siswa"))
	fmt.Println("0. " + utils.ColorMessage("red", "Keluar"))
	fmt.Print(utils.ColorMessage("yellow", "Pilih opsi: "))
}
