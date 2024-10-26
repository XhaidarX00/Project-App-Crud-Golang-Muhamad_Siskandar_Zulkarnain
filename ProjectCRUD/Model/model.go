package model

import "time"

// Model Structures
type Student struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Course    []Schedule `json:"course"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Schedule struct {
	Course   string
	Schedule []string `json:"schedule"`
}

type Course struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Teacher    string    `json:"teacher"`
	Credits    int       `json:"credits"`
	Activate   bool      `json:"activation"`
	StudentIDs []string  `json:"student_ids"`
	Schedule   []string  `json:"schedule"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Database struct
type Database struct {
	Students []Student  `json:"students"`
	Courses  []Course   `json:"courses"`
	Login    []Account  `json:"login"`
	Schedule []Schedule `json:"schedule"`
}
