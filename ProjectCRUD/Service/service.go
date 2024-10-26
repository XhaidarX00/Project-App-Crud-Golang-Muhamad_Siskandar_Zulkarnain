package service

import (
	model "ProjectCRUD/Model"
	utils "ProjectCRUD/Utils"
	view "ProjectCRUD/VIew"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type EducationService struct {
	db     *model.Database
	dbFile string
}

// NewEducationService creates a new instance of EducationService
func NewEducationService(dbFile string) (*EducationService, error) {
	if dbFile == "" {
		dbFile = "education_db.json"
	}

	service := &EducationService{
		db: &model.Database{
			Students: []model.Student{},
			Courses:  []model.Course{},
			Login:    []model.Account{},
			Schedule: []model.Schedule{},
		},
		dbFile: dbFile,
	}

	if err := service.loadDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize service: %w", err)
	}

	return service, nil
}

// Database operations
func (s *EducationService) loadDatabase() error {
	data, err := ioutil.ReadFile(s.dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return s.saveDatabase()
		}
		return err
	}
	return json.Unmarshal(data, s.db)
}

func (s *EducationService) saveDatabase() error {

	data, err := json.MarshalIndent(s.db, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}
	return ioutil.WriteFile(s.dbFile, data, 0644)
}

// Account operations
func (s *EducationService) InitializeAccounts() {
	if len(s.db.Login) == 0 {
		s.db.Login = append(s.db.Login, model.Account{
			Username: "admin",
			Password: "admin123",
		})
		s.saveDatabase()
	}
}

// login handler
func (s *EducationService) AuthenticateUser() bool {
	s.InitializeAccounts()
	var username string
	var password string

	for {
		fmt.Println("\n=---------------= 🔑 Login 🔑 =---------------=")
		fmt.Print(utils.ColorMessage("yellow", "Masukkan username : "))
		fmt.Scan(&username)
		fmt.Print(utils.ColorMessage("yellow", "Masukkan password : "))
		fmt.Scan(&password)
		utils.ClearScreen()

		for _, account := range s.db.Login {
			if account.Username == username {
				if account.Password == password {
					utils.SuccesMessage("Selamat Datang Di Program Management Course Golang...")
					return true
				}
				utils.ErrorMessage("Password Salah")
				continue
			} else {
				utils.ErrorMessage("Pengguna Tidak Ditemukan")
				continue
			}
		}
	}
}

// Student operations

// Create Student
func (s *EducationService) CreateStudent() {
	var student model.Student

	if len(s.db.Courses) == 0 {
		utils.ErrorMessage("Course Tidak Tersedia")
		return
	}

	dataC := utils.ConvertSliceToMap(s.db.Courses)
	dkC := []string{"CreatedAt", "UpdatedAt"}
	keysC := utils.GetStructKeys(model.Student{}, dkC)

	fmt.Print(utils.ColorMessage("yellow", "Masukkan Nama Siswa : "))
	fmt.Scan(&student.Name)
	s.DisplayCountIndex("DATA COURSE", dataC, keysC, &student)

	student.ID = fmt.Sprintf("STD%d", len(s.db.Students)+1)
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()
	s.db.Students = append(s.db.Students, student)

	s.saveDatabase()
	view.PrintStudentToJson(s.db.Students)
}

// Display Students
func (s *EducationService) GetAllStudents() {
	if len(s.db.Students) == 0 {
		utils.ErrorMessage("Database Siswa Kosong!")
		return
	}
	view.PrintStudentToJson(s.db.Students)
}

// Update Student
func (s *EducationService) UpdateStudent() {
	var id string
	var skipCourse string
	var updates model.Student

	if len(s.db.Students) == 0 {
		if err := s.loadDatabase(); err != nil {
			return
		}
	}

	data := utils.ConvertSliceToMap(s.db.Students)
	dk := []string{"CreatedAt", "UpdatedAt"}
	keys := utils.GetStructKeys(model.Student{}, dk)

	dataC := utils.ConvertSliceToMap(s.db.Courses)
	dkC := []string{"CreatedAt", "UpdatedAt"}
	keysC := utils.GetStructKeys(model.Student{}, dkC)

	for {
		id = DisplayDataChosenIndexOne("DATA SISWA", data, keys)
		fmt.Print(utils.ColorMessage("yellow", "Update Nama Baru (skip jika tidak diubah) : "))
		fmt.Scan(&updates.Name)
		fmt.Print(utils.ColorMessage("yellow", "Update Course (y / t) : "))
		fmt.Scan(&skipCourse)
		if strings.ToLower(skipCourse) != "t" {
			s.DisplayCountIndex("DATA COURSE", dataC, keysC, &updates)
		}
		utils.ClearScreen()

		err := s.validateStudent(updates)
		if !err {
			continue
		} else {
			break
		}
	}

	for i, student := range s.db.Students {
		if student.ID == id {
			if strings.ToLower(updates.Name) != "skip" {
				s.db.Students[i].Name = updates.Name
			}
			s.db.Students[i].UpdatedAt = time.Now()
			s.saveDatabase()
			return
		}
	}
	utils.ErrorMessage("Siswa Tidak Ditemukan Gagal Update!")
}

func (s *EducationService) DeleteStudent(id string) {
	for i, student := range s.db.Students {
		if student.ID == id {
			s.db.Students = append(s.db.Students[:i], s.db.Students[i+1:]...)
			s.saveDatabase()
			utils.SuccesMessage("Berhasil Update Data Siswa")
			return
		}
	}
	utils.ErrorMessage("Siswa Tidak Ditemukan!")
}

// // Course operations
func (s *EducationService) CreateCourse() {
	var course model.Course
	var hari string

	for {
		fmt.Print("Masukkan Nama Course : ")
		fmt.Scan(&course.Name)
		for _, cr := range s.db.Courses {
			if cr.Name == course.Name {
				utils.ErrorMessage("Nama Course Sudah Terdaftar")
				return
			}
		}

		fmt.Print("Masukkan Nama Guru : ")
		fmt.Scan(&course.Teacher)
		fmt.Print("Masukkan Harga : ")
		fmt.Scan(&course.Credits)
		course.Activate = false
		course.StudentIDs = []string{}

		// Input jadwal course
		course.Schedule = []string{}
		for {
			fmt.Print("Masukkan Jadwal Course (atau ketik 'selesai' untuk berhenti): ")
			fmt.Scan(&hari)
			if hari == "selesai" {
				break
			}
			course.Schedule = append(course.Schedule, hari)
		}

		err := s.validateCourse(course)
		if !err {
			continue
		} else {
			break
		}
	}

	course.ID = fmt.Sprintf("CRS%d", len(s.db.Courses)+1)
	course.Activate = true
	course.CreatedAt = time.Now()
	course.UpdatedAt = time.Now()
	course.StudentIDs = []string{}

	s.db.Courses = append(s.db.Courses, course)
	s.saveDatabase()
	view.PrintCourseToJson(s.db.Courses)
}

func (s *EducationService) AddStudentToCourse(courseID string, studentID string) {
	for i, course := range s.db.Courses {
		if course.ID == courseID {
			s.db.Courses[i].StudentIDs = append(s.db.Courses[i].StudentIDs, studentID)
			s.db.Courses[i].UpdatedAt = time.Now()
			s.saveDatabase()
			utils.SuccesMessage("Student berhasil ditambahkan ke course")
			return
		}
	}

	msg := fmt.Sprintf("Course dengan ID '%s' tidak ditemukan", courseID)
	utils.ErrorMessage(msg)
}

func (s *EducationService) GetAllCourses() []model.Course {
	if len(s.db.Courses) == 0 {
		utils.ErrorMessage("Tidak Ada Course Yang Aktif")
		return nil
	}

	courses := make([]model.Course, len(s.db.Courses))
	copy(courses, s.db.Courses)
	return courses
}

// Validation helpers
func (s *EducationService) validateStudent(student model.Student) bool {
	if student.Name == "" {
		utils.ErrorMessage("Nama Siswa Tidak Boleh Kosong!")
		return false
	}

	return true
}

func (s *EducationService) validateCourse(course model.Course) bool {
	if course.Name == "" {
		utils.ErrorMessage("Nama Kursus Tidak Boleh Kosong!")
		return false
	}
	if course.Teacher == "" {
		utils.ErrorMessage("Nama Guru Tidak Boleh Kosong!")
		return false
	}
	if course.Credits < 1 {
		utils.ErrorMessage("Credits Tidak Boleh Kosong atau kurang dari 0!")
		return false
	}

	return true
}

func ResetSessionTimeout() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	deadline := time.Now().Add(10 * time.Second)
	return context.WithDeadline(ctx, deadline)
}

func ExitMainmenu() {
	defer os.Exit(0)
	utils.ClearScreen()
	utils.SuccesMessage("Keluar dari Program\n")
}

func (s *EducationService) DisplayCountIndex(title string, data []map[string]interface{}, keys []string, student *model.Student) {
	var input string
	var tempChoice []string
	defer func() {
		tempChoice = []string{}
	}()

	for {
		DisplayData(title, data, keys)
		fmt.Print(utils.ColorMessage("yellow", "Masukkan Kelas Course atau ketik `done` untuk berhenti : "))
		fmt.Scan(&input)
		utils.ClearScreen()

		if strings.ToLower(input) == "done" {
			return
		}

		intInput, err := strconv.Atoi(input)
		if err != nil || intInput < 1 || intInput > len(data) {
			msg := fmt.Sprintf("Input harus berupa angka yang valid dan tidak boleh lebih dari %d\n", len(data))
			utils.ErrorMessage(msg)
			continue
		}

		selectedCourse := s.db.Courses[intInput-1]

		newSchedule := model.Schedule{
			Course:   selectedCourse.Name,
			Schedule: selectedCourse.Schedule,
		}

		if isScheduleConflict(student.Course, newSchedule) {
			conflictMsg := fmt.Sprintf("Course '%s' bentrok dengan course yang sudah dipilih sebelumnya", selectedCourse.Name)
			utils.ErrorMessage(conflictMsg)
			continue
		}

		student.Course = append(student.Course, newSchedule)

		msg := fmt.Sprintf("Anda memilih course: %s", utils.ColorMessage("green", selectedCourse.Name))
		tempChoice = append(tempChoice, msg)
		displayTempChoice(tempChoice)

		s.AddStudentToCourse(selectedCourse.ID, student.ID)
	}

}

func isScheduleConflict(existingSchedules []model.Schedule, newSchedule model.Schedule) bool {
	for _, existing := range existingSchedules {
		for _, newTime := range newSchedule.Schedule {
			for _, existingTime := range existing.Schedule {
				if newTime == existingTime && newTime != "" {
					return true
				}
			}
		}
	}
	return false
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

		displayTempChoice(tempChoice)
	}
}

func displayTempChoice(tempChoice []string) {
	for _, msg := range tempChoice {
		fmt.Println(msg)
	}
}

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
			fmt.Printf("%v | ", value)
		}
		fmt.Println()
	}
	fmt.Println(strings.Repeat("-", 50))
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
