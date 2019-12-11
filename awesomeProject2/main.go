package main

import (
	memorycache "awesomeProject2/cache"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Product struct {
	Id      int
	Model   string
	Company string
	Price   int
}
type Producer struct {
	Id      int
	Company string
	Numberofmodels   int
}

var database *sql.DB

func ShowAllProducts(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from productdb.Products")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	products := []Product{}

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, products)
}
func ShowAllProducers(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from productdb.Producers")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	producers := []Producer{}

	for rows.Next() {
		p := Producer{}
		err := rows.Scan(&p.Id, &p.Company, &p.Numberofmodels)
		if err != nil {
			fmt.Println(err)
			continue
		}
		producers = append(producers, p)
	}

	tmpl, _ := template.ParseFiles("templates/indexProducers.html")
	tmpl.Execute(w, producers)
}
func ShowReqProduct(w http.ResponseWriter, r *http.Request, model string) {
	str := "SELECT * FROM productdb.Products WHERE model = " + "'" + model + "'" + ";"
	rows, err := database.Query(str)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	products := []Product{}
	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, products)
}
func GetInfo(db *sql.DB, company string) (string, error) {
	str := "SELECT * FROM productdb.Products WHERE company = " + "'" + company + "'" + ";"
	res := db.QueryRow(str)
	prod := Product{}
	if res != nil {
		err := res.Scan(&prod.Id, &prod.Model, &prod.Company, &prod.Price)
		fmt.Print(err)
		s := prod.Model + ", Price==" + strconv.Itoa(prod.Price)
		if prod.Model == "" {
			return "No information", nil
		} else {
			return s, err
		}
	}
	return "No information", nil
}
func main() {
	altdatabase, err := sql.Open("mysql", "root:211user504UI@/productdb")
	fmt.Println("Работать ли с cache?\n да - введите 1\n нет - введите 0")
	var check int
	fmt.Fscan(os.Stdin, &check)
	if check == 1 {
		fmt.Println("Сохраним в кеше информацию о компании с бд. Время жизни контейнера одна минута, как и нашей записи.")
		cache := memorycache.New(time.Minute, 10*time.Minute)
		str, _ := GetInfo(altdatabase, "Lenovo")
		cache.Set("myKey", str, time.Minute)
		fmt.Println("Хотим вывести информацию о компании?\n да - введите 1\n нет - введите 0")
		fmt.Fscan(os.Stdin, &check)
		if check == 1 {
			i, b := cache.Get("myKey")
			fmt.Printf("%s %t", "Информация существует?", b)
			fmt.Printf("%s %s", "\nИнформация: ", i)
		} else {
			fmt.Println("Хотим удалить из кеша или подождать пока закончиться время жизни?\n Удалить - 1\n Ждать - 0")
			fmt.Fscan(os.Stdin, &check)
			if check == 0 {
				fmt.Println("Подождем минуту, чтобы отчистился кэш")
				time.Sleep(time.Minute)
				i, b := cache.Get("myKey")
				fmt.Printf("%s %t", "Информация существует?", b)
				fmt.Printf("%s %s", "\nИнформация: ", i)
				fmt.Println("\nПроверим действительно ли у объекта закончилось время жизни\nвыведем массив ключей у которых закончилось время")
				fmt.Println(cache.ExpiredKeys())
			} else if check == 1 {
				cache.Delete("myKey")
				i, b := cache.Get("myKey")
				fmt.Printf("%s %t", "Информация существует?", b)
				fmt.Printf("%s %s", "\nИнформация: ", i)
				fmt.Println("\nПроверим действительно ли удалили то что хранилось в кеше под нашим ключем или у объекта закончилось время жизни\nвыведем массив ключей у которых закончилось время")
				fmt.Println(cache.ExpiredKeys())
			}
		}
	}

	if err != nil {
		log.Println(err)
	}
	database = altdatabase
	defer database.Close()
	http.HandleFunc("/", ShowAllProducers)

	fmt.Println("\nServer is listening...")
	http.ListenAndServe(":8181", nil)
}
