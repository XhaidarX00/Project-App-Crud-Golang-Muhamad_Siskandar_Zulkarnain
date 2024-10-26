package main

import (
	service "ProjectCRUD/Service"
	utils "ProjectCRUD/Utils"
	view "ProjectCRUD/VIew"
	"context"
	"fmt"
	"time"
)

func main() {
	utils.ClearScreen()
	educationService, err := service.NewEducationService("education_db.json")
	if err != nil {
		fmt.Printf("Error loading database: %v\n", err)
		return
	}

	if educationService.AuthenticateUser() {
		ctx := context.Background()
		deadline := time.Now().Add(3000 * time.Second)
		ctxwithdeadline, cancel := context.WithDeadline(ctx, deadline)
		defer cancel()
		menu(educationService, ctxwithdeadline)
	}
}

func menu(educationService *service.EducationService, ctxwithdeadline context.Context) {
	for {
		select {
		case <-ctxwithdeadline.Done():
			utils.ErrorMessage("Akses Ditolak, Sesi Login Expired!")
			if educationService.AuthenticateUser() {
				ctxWithDeadline, cancel := service.ResetSessionTimeout()
				defer cancel()
				ctxwithdeadline = ctxWithDeadline
				continue
			} else {
				utils.ErrorMessage("Terjadi Kesalahan Saat Login! Hubungi Developer untuk Maintance!")
				service.ExitMainmenu()
			}

		default:
			// time.Sleep(1 * time.Second)
			view.DisplayMenu()

			var choice int
			fmt.Scan(&choice)
			utils.ClearScreen()

			switch choice {
			case 1:
				educationService.CreateCourse()
			case 2:
				educationService.CreateStudent()
			case 3:
				educationService.GetAllStudents()

			case 4:
				educationService.UpdateStudent()

			case 5:
				// Menghapus Siswa
				var id string
				fmt.Print("Masukkan ID Siswa yang akan dihapus: ")
				fmt.Scan(&id)
				educationService.DeleteStudent(id)

			case 6:
				// Membuat Course Baru
				educationService.CreateCourse()

			case 7:
				// Membaca Semua Course
				courses := educationService.GetAllCourses()
				if courses == nil {
					fmt.Println("Tidak ada course yang aktif.")
				} else {
					for _, c := range courses {
						fmt.Printf("ID: %s, Nama: %s, Guru: %s, Kredit: %d\n", c.ID, c.Name, c.Credits)
					}
				}

			case 0:
				utils.SuccesMessage("Program Selesai dan Semoga hari mu menyenangkan!")
				service.ExitMainmenu()

			default:
				fmt.Println("Invalid option!")
			}
		}
	}
}
