package view

import (
	model "ProjectCRUD/Model"
	utils "ProjectCRUD/Utils"
	"encoding/json"
	"fmt"
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

func DisplayMenu() {
	fmt.Println(utils.ColorMessage("yellow", "\n=---------- Menu ----------="))
	fmt.Println("1. " + utils.ColorMessage("green", "Tambah Course"))
	fmt.Println("2. " + utils.ColorMessage("green", "Tambah Siswa"))
	fmt.Println("3. " + utils.ColorMessage("green", "Lihat Semua Siswa"))
	fmt.Println("4. " + utils.ColorMessage("green", "Update Data Siswa"))
	fmt.Println("5. " + utils.ColorMessage("green", "Hapus Siswa"))
	fmt.Println("6. " + utils.ColorMessage("green", "Tambah Course"))
	fmt.Println("7. " + utils.ColorMessage("green", "Lihat Semua Course"))
	fmt.Println("0. " + utils.ColorMessage("red", "Keluar"))
	fmt.Print(utils.ColorMessage("yellow", "Pilih opsi: "))
}
