package data

type User struct {
	Id          int
	Name        string
	Description string
	Type        string
	Options     []Option
	Quantity    int
	Sold        int
	Cost        int
}

type Option struct {
	Id       int
	Name     string
	Quantity int
	Sold     int
}

func GetUser() {

}

func CreateUser() {

}

func UpdateUser() {

}

func DeleteUser() {

}
