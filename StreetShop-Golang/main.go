package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Article struct {
	ID          int
	Name        string
	Price       float64
	Description string
	Detail      string
	Stock       int
	Discount    bool
	Image       string
}

var articles = []Article{
	{ID: 1, Name: "Palace International Hood 'Black'", Price: 29.99, Description: "Palace International Hood 'Black", Detail: "Confectionné dans un tissu doux et confortable pour l'hiver, le pull est un incontournable pour montrer que tu représentes la cité londonienne", Stock: 10, Discount: false, Image: "/static/img/products/noirc.webp"},
	{ID: 2, Name: "Palace Pull a capuche unisexe chasseur", Price: 29.99, Description: "Palace Pull a capuche unisexe chasseur", Detail: "Confectionné dans un tissu doux et confortable pour l'hiver, le pull est un incontournable pour montrer que tu représentes la cité londonienne", Stock: 15, Discount: true, Image: "/static/img/products/vert.webp"},
	{ID: 3, Name: "Palace Pull a capuche marine", Price: 29.99, Description: "Palace Pull a capuche marine", Detail: "Confectionné dans un tissu doux et confortable pour l'hiver, le pull est un incontournable pour montrer que tu représentes la cité londonienne", Stock: 5, Discount: false, Image: "/static/img/products/bleu.webp"},
	{ID: 4, Name: "Palace washed terry 1/4 placket hood mojito", Price: 29.99, Description: "Palace washed terry 1/4 placket hood mojito", Detail: "Confectionné dans un tissu doux et confortable pour l'hiver, le pull est un incontournable pour montrer que tu représentes la cité londonienne", Stock: 20, Discount: true, Image: "/static/img/products/jaune.webp"},
	{ID: 5, Name: "Palace pull crew passepose noir", Price: 19.99, Description: "Palace pull crew passepose noir", Detail: "Confectionné dans un tissu doux et confortable pour l'hiver, le pull est un incontournable pour montrer que tu représentes la cité londonienne", Stock: 0, Discount: false, Image: "/static/img/products/noirs.webp"},
}

func mult(a, b float64) float64 {
	return a * b
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("home.html").Funcs(template.FuncMap{"mult": mult}).ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, articles)
}

func nonPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("non.html").Funcs(template.FuncMap{"mult": mult}).ParseFiles("templates/non.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, articles)
}

func articlePage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	for _, article := range articles {
		if strconv.Itoa(article.ID) == id {
			tmpl, err := template.New("article.html").Funcs(template.FuncMap{"mult": mult}).ParseFiles("templates/article.html")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl.Execute(w, article)
			return
		}
	}
	http.NotFound(w, r)
}

func addProductPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		imageDir := "static/img/products/"
		files, err := ioutil.ReadDir(imageDir)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture du dossier des images", http.StatusInternalServerError)
			return
		}

		var images []string
		for _, file := range files {
			if !file.IsDir() {
				images = append(images, file.Name())
			}
		}

		tmpl, err := template.ParseFiles("templates/add.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, images)
	} else if r.Method == http.MethodPost {
		name := r.FormValue("name")
		price := atof(r.FormValue("price"))
		description := r.FormValue("description")
		detail := r.FormValue("detail")
		stock := atoi(r.FormValue("stock"))
		discount := r.FormValue("discount") == "on"

		selectedImage := r.FormValue("image")

		newArticle := Article{
			ID:          len(articles) + 1,
			Name:        name,
			Price:       price,
			Description: description,
			Detail:      detail,
			Stock:       stock,
			Discount:    discount,
			Image:       "/static/img/products/" + selectedImage,
		}

		articles = append(articles, newArticle)

		http.Redirect(w, r, "/article?id="+strconv.Itoa(newArticle.ID), http.StatusSeeOther)
	}
}

func atof(str string) float64 {
	value, _ := strconv.ParseFloat(str, 64)
	return value
}

func atoi(str string) int {
	value, _ := strconv.Atoi(str)
	return value
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/article", articlePage)
	http.HandleFunc("/add", addProductPage)
	http.HandleFunc("/non", nonPage)

	http.ListenAndServe("localhost:8080", nil)
}
