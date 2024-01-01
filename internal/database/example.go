package database

func example() {
	//type Product struct {
	//	Name      string
	//	Price     float64
	//	Available bool
	//}
	//
	//func main() {
	//	//url для подключения
	//	connStr := "postgres://postgres:secret@localhost:5432/goDBTest?sslmode=disable"
	//
	//	//открываем подключение
	//	db, err := sql.Open("postgres", connStr)
	//	//закрываем подключение в конце программы
	//	defer db.Close()
	//
	//	//обработка ошибок
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	//проверяем подключение
	//	if err = db.Ping(); err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	//создаём таблицу
	//	createProductTable(db)
	//
	//	data := []Product{}
	//	rows, err := db.Query("SELECT name, avaialble, price FROM product")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	defer rows.Close()
	//	//to scan DB vals
	//	var name string
	//	var available bool
	//	var price float64
	//
	//	for rows.Next() {
	//		err := rows.Scan(&name, &available, &price)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		data = append(data, Product{name, price, available})
	//	}
	//
	//	fmt.Println(data)
	//}
	//
	//// в параметрах указываем с какой базой данных работаем, с каким соединением
	//func createProductTable(db *sql.DB) {
	//	//запрос
	//	query := `CREATE TABLE IF NOT EXISTS product (
	//		id SERIAL PRIMARY KEY,
	//		name VARCHAR(100) NOT NULL,
	//		price NUMERIC(6,2) NOT NULL,
	//		available BOOLEAN,
	//		created timestamp DEFAULT NOW()
	//)`
	//
	//	//выполняем его
	//	_, err := db.Exec(query)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//
	//func insertProduct(db *sql.DB, product Product) int {
	//	//запрос
	//	query := `INSERT INTO product (name,price,available)
	//		VALUES ($1, $2, $3) RETURNING id`
	//	var pk int
	//	//не можем захардкодить как с db.Exec(query) прежде
	//	//вместо знаков доллара, после запроса перечисляем что подставим туда
	//	//считаем возвращаемое значение методом Scan и положим значение в pk
	//
	//	//Scan() помогает вернуть значения из запроса и положить их в другие переменные
	//	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
	//
	//	if err != nil {
	//	log.Fatal(err)
	//}
	//	//возвращаем айди новой записи
	//	return pk
	//}
	//
	////product1 := Product{
	////Name:      "Green Elephant",
	////Price:     228,
	////Available: true,
	////}
	////
	//////возвращаемый айди поместим в переменную pk
	////pk := insertProduct(db, product1)
	////
	////var name string
	////var avaialble bool
	////var price float64
	////
	////query := "SELECT name, avaialble, price FROM product WHERE id = $1"
	//////сканом передаём значения в переменные при помощи указателя
	////err = db.QueryRow(query, pk).Scan(&name, &avaialble, &price)
	////if err != nil {
	////if err == sql.ErrNoRows {
	////log.Fatalf("no rows found with id: ", pk)
	////}
	////log.Fatal(err)
	////}
	////
	////fmt.Println("name: ", name)
	////fmt.Println("available: ", avaialble)
	////fmt.Println("price: ", price)
	////
	////fmt.Println("ID = ", pk)
}
