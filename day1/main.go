package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type CVProject struct {
	Project       string
	Company       string
	Role          string
	ProjectPeriod string
}

type Contact struct {
	Name, Email, Phone string
	WillAttend         bool
}

var responses = make([]*Contact, 0, 10)
var templates = make(map[string]*template.Template, 4)

func LoadTemplates() {
	templateNames := [4]string{"welcome", "form", "list", "thanks"}
	for index, name := range templateNames {
		t, err := template.ParseFiles("layout.html", name+".html")
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		} else {
			panic(err)
		}
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	templates["welcome"].Execute(w, nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	templates["list"].Execute(w, CVProjects)
}

type formData struct {
	*Contact
	Errors []string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates["form"].Execute(w, formData{
			Contact: &Contact{}, Errors: []string{},
		})

	} else if r.Method == http.MethodPost {
		r.ParseForm()
		responseData := Contact{
			Name:  r.FormValue("name"),
			Email: r.FormValue("email"),
			Phone: r.FormValue("phone"),
		}

		errors := []string{}

		if responseData.Name == "" {
			errors = append(errors, "Укажите ваше Имя!")
		}
		if responseData.Email == "" {
			errors = append(errors, "Укажите Вам email!")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Укажите Ваш телефон!")
		}
		if len(errors) > 0 {
			templates["form"].Execute(w, formData{
				Contact: &responseData, Errors: errors,
			})
		} else {
			responses = append(responses, &responseData)
			templates["thanks"].Execute(w, responseData.Name)
		}
	}
}

var CVProjects = []CVProject{
	{
		Project:       "Миграция базы данных на Oracle 19c",
		Company:       "ООО «КиберБабушка»",
		Role:          "Главный внук по базам данных",
		ProjectPeriod: "2021-2023",
	},
	{
		Project:       "Разработка ERP-системы для склада",
		Company:       "ЗАО «ХранимВсе»",
		Role:          "Программист хранимых процедур",
		ProjectPeriod: "2019-2021",
	},
	{
		Project:       "Оптимизация запросов в Oracle",
		Company:       "ПАО «Запрос-миллионник»",
		Role:          "Победитель планировщика",
		ProjectPeriod: "2017-2019",
	},
}

func main() {
	LoadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
