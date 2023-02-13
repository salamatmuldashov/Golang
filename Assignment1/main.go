package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func menu() int {
	var x int
	fmt.Println("Enter 0 to exit")
	fmt.Println("Enter 1 to registration")
	fmt.Println("Enter 2 to authorization")
	fmt.Scanln(&x)
	return x
}

func userMenu() int {
	var x int
	fmt.Println("Enter 0 to log out")
	fmt.Println("Enter 1 to search")
	fmt.Println("Enter 2 to rate item")
	fmt.Println("Enter 3 filter")
	fmt.Println("Enter 4 to sort")
	fmt.Println("Enter 5 to see list of items")
	fmt.Scanln(&x)
	return x
}

type user struct {
	name     string
	surname  string
	password string
	mail     string
}

var users []user

func registration(newName, newSurname, newPassword string) {
	newUser := user{
		name:     newName,
		surname:  newSurname,
		password: newPassword,
	}
	newUser.constructMail(newName, newSurname)

	users = append(users, newUser)
}

func authorization(mail, password string) user {
	var user user
	for i := 0; i < len(users); i++ {
		if users[i].getMail() == mail {
			if users[i].getPassword() == password {
				user = users[i]
				fmt.Println("Welcome " + user.getName() + " " + user.getSurname())
				break
			} else {
				fmt.Println("Password is not correct! Try again!")
			}
		}
	}
	return user
}

func findUser(mail string) bool {
	for i := 0; i < len(users); i++ {
		if users[i].getMail() == mail {
			return true
		}
	}
	return false
}

// Setters, Getters
func (u *user) getName() string {
	return u.name
}
func (u *user) getMail() string {
	return u.mail
}
func (u *user) setMail(newMail string) {
	u.mail = newMail
}
func (u *user) getPassword() string {
	return u.password
}

func (u *user) setName(newName string) {
	u.name = newName
}

func (u *user) setSurname(newSurname string) {
	u.surname = newSurname
}

func (u *user) setPassword(newPassword string) {
	u.password = newPassword
}
func (u *user) getSurname() string {
	return u.surname
}

// Setters, Getters

func (u *user) constructMail(name, surname string) {
	var mail string
	firstChar := name[0:1]
	firstChar = strings.ToLower(firstChar)
	curSurname := strings.ToLower(surname)
	mail = firstChar + "_" + curSurname
	u.setMail(mail)
}

func (u *user) search(name string) {
	fmt.Println()
	for i := 0; i < len(items); i++ {
		if items[i].getName() == name {
			fmt.Println(items[i].itemsInfo() + "\n")
			return
		}
	}
	fmt.Println("Item does not exist!")
}

func (u *user) rate(name string, rating float64) {
	for i := 0; i < len(items); i++ {
		if items[i].getName() == name {
			items[i].setRating(rating)
			return
		}
	}
	fmt.Println("Item does not exist!")
}

type item struct {
	name   string
	rating float64
	price  float64
}

var items = []item{
	item{
		name:   "Phone",
		rating: 0.0,
		price:  1000.0,
	},
	item{
		name:   "Watch",
		rating: 0.0,
		price:  500.0,
	},
	item{
		name:   "Computer",
		rating: 0.0,
		price:  1500.0,
	},
	item{
		name:   "TV",
		rating: 0.0,
		price:  2000.0,
	},
}

func (i *item) getName() string {
	return i.name
}
func (i *item) getRating() float64 {
	return i.rating
}
func (i *item) getPrice() float64 {
	return i.price
}
func (i *item) setName(newName string) {
	i.name = newName
}
func (i *item) setRating(newRating float64) {
	i.rating = newRating
}
func (i *item) setPrice(newPrice float64) {
	i.price = newPrice
}

func (i *item) itemsInfo() string {
	rating := fmt.Sprint(i.rating)
	price := fmt.Sprint(i.price)
	return "Name: " + i.getName() + "\nRating: " + rating + "\nPrice: " + price
}

func list() {
	fmt.Println()
	for i := 0; i < len(items); i++ {
		fmt.Println(items[i].itemsInfo())
		fmt.Println()
	}

}

func (u *user) filterByRating(rating float64, str string) {
	var res []item
	if str == "1" {
		for i := 0; i < len(items); i++ {
			if items[i].getRating() <= rating {
				res = append(res, items[i])
			}
		}

	} else if str == "2" {
		for i := 0; i < len(items); i++ {
			if items[i].getRating() >= rating {
				res = append(res, items[i])
			}
		}

	} else {
		fmt.Println("Error!")
		return
	}
	fmt.Println()
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i].itemsInfo())
		fmt.Println()
	}

}

func (u *user) filterByPrice(price float64, str string) {
	var res []item
	if str == "1" {
		for i := 0; i < len(items); i++ {
			if items[i].getPrice() <= price {
				res = append(res, items[i])
			}
		}

	} else if str == "2" {
		for i := 0; i < len(items); i++ {
			if items[i].getPrice() >= price {
				res = append(res, items[i])
			}
		}

	} else {
		fmt.Println("Error!")
		return
	}
	fmt.Println()
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i].itemsInfo())
		fmt.Println()
	}

}

func (u *user) sortByRating(str int) {
	res := items
	if str == 1 {
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].rating < res[j].rating
		})
	} else if str == 2 {
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].rating > res[j].rating
		})
	} else {
		fmt.Println("Error!")
		return
	}
	fmt.Println()
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i].itemsInfo())
		fmt.Println()
	}

}

func (u *user) sortByPrice(str int) {
	res := items
	if str == 1 {
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].price < res[j].price
		})
	} else if str == 2 {
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].price > res[j].price
		})
	} else {
		fmt.Println("Error!")
		return
	}
	fmt.Println()
	for i := 0; i < len(res); i++ {
		fmt.Println(res[i].itemsInfo())
		fmt.Println()
	}

}

func main() {
	command := menu()
	for true {
		if command == 0 {
			os.Exit(0)
		} else if command == 1 {
			var name, surname, password string
			fmt.Println("Enter name")
			fmt.Scanln(&name)
			fmt.Println("Enter surname")
			fmt.Scanln(&surname)
			fmt.Println("Enter password")
			fmt.Scanln(&password)
			registration(name, surname, password)
		} else if command == 2 {
			var mail, password string
			fmt.Println("Enter mail")
			fmt.Scanln(&mail)
			fmt.Println("Enter password")
			fmt.Scanln(&password)
			if findUser(mail) == true {
				user := authorization(mail, password)
				if user.getPassword() == password {
					command := userMenu()
					for true {
						if command == 0 {
							break
						} else if command == 1 {
							var name string
							fmt.Println("Enter name of item")
							fmt.Scanln(&name)
							user.search(name)
						} else if command == 2 {
							var name string
							var rating float64
							fmt.Println("Enter name of item")
							fmt.Scanln(&name)
							fmt.Println("Rate the item")
							fmt.Scanln(&rating)
							user.rate(name, rating)
						} else if command == 3 {
							var command int
							fmt.Println("Enter 1 to filter by rating")
							fmt.Println("Enter 2 to filter by price")
							fmt.Scanln(&command)
							if command == 1 {
								var rating float64
								fmt.Println("Rating:")
								fmt.Scanln(&rating)
								var str string
								fmt.Println("1:Under\n2:Over")
								fmt.Scanln(&str)
								user.filterByRating(rating, str)
							} else if command == 2 {
								var price float64
								fmt.Println("Price:")
								fmt.Scanln(&price)
								var str string
								fmt.Println("1:Under\n2:Over")
								fmt.Scanln(&str)
								user.filterByPrice(price, str)
							} else {
								fmt.Println("Error!")
							}
						} else if command == 4 {
							var command int
							fmt.Println("Enter 1 to sort by rating")
							fmt.Println("Enter 2 to sort by price")
							fmt.Scanln(&command)
							if command == 1 {
								var command int
								fmt.Println("1:Ascending\n2:Descending")
								fmt.Scanln(&command)
								user.sortByRating(command)
							} else if command == 2 {
								var command int
								fmt.Println("1:Ascending\n2:Descending")
								fmt.Scanln(&command)
								user.sortByPrice(command)
							} else {
								fmt.Println("Error!")
							}
						} else if command == 5 {
							list()
						} else {
							fmt.Println("Error command!")
						}
						command = userMenu()
					}
				}

			} else {
				fmt.Println("Such a user does not exist! Please register!")
			}

		} else {
			fmt.Println("Error! Unknown command!")
		}
		command = menu()
	}

}
