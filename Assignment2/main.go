package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Item struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Price   float64   `json:"price"`
	Rating  float64   `json:"rating"`
	Ratings []float64 `json:"ratings"`
}

type Database struct {
	Users []*User `json:"users"`
	Items []*Item `json:"items"`
}

type System struct {
	Database *Database
}

func (s *System) Register(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error parsing request body: %v", err)
		return
	}
	for _, existingUser := range s.Database.Users {
		if existingUser.Username == user.Username {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "user %q already exists", user.Username)
			return
		}
	}
	user.ID = len(s.Database.Users) + 1
	s.Database.Users = append(s.Database.Users, &user)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *System) Authorize(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error parsing request body: %v", err)
		return
	}
	for _, existingUser := range s.Database.Users {
		if existingUser.Username == user.Username && existingUser.Password == user.Password {
			fmt.Fprint(w, "Welcome "+existingUser.Username+"\n")
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprint(w, "invalid username or password\n")
}

func (s *System) SearchItemsByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	var items []*Item
	for _, item := range s.Database.Items {
		if name == "" || item.Name == name {
			items = append(items, item)
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func (s *System) FilterItems(w http.ResponseWriter, r *http.Request) {
	minPrice := r.URL.Query().Get("min_price")
	maxPrice := r.URL.Query().Get("max_price")
	minRating := r.URL.Query().Get("min_rating")
	maxRating := r.URL.Query().Get("max_rating")
	var minPriceFloat, maxPriceFloat, minRatingFloat, maxRatingFloat float64
	fmt.Sscan(minPrice, &minPriceFloat)
	fmt.Sscan(maxPrice, &maxPriceFloat)
	fmt.Sscan(minRating, &minRatingFloat)
	fmt.Sscan(maxRating, &maxRatingFloat)
	var filteredItems []*Item
	for _, item := range s.Database.Items {
		if (minPrice == "" || item.Price >= minPriceFloat) &&
			(maxPrice == "" || item.Price <= maxPriceFloat) &&
			(minRating == "" || item.Rating >= minRatingFloat) &&
			(maxRating == "" || item.Rating <= maxRatingFloat) {
			filteredItems = append(filteredItems, item)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(filteredItems)
}

func (s *System) RateItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("item_id")
	itemRating := r.URL.Query().Get("rating")
	curItemID, err1 := strconv.Atoi(itemID)
	curItemRating, err2 := strconv.ParseFloat(itemRating, 64)

	if err1 != nil {
		fmt.Println(err1)
	}
	if err2 != nil {
		fmt.Println(err2)
	}

	var item *Item
	for _, existingItem := range s.Database.Items {
		if existingItem.ID == curItemID {
			item = existingItem
			break
		}
	}

	// fmt.Println(item)
	// var rating float64
	// err := json.NewDecoder(r.Body).Decode(&rating)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	fmt.Fprintf(w, "error parsing request body: %v", err)
	// 	return
	// }

	item.Ratings = append(item.Ratings, curItemRating)
	var totalRating float64
	for _, r := range item.Ratings {
		totalRating += r
	}
	item.Rating = totalRating / float64(len(item.Ratings))

	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(item)
}

func main() {
	database := &Database{
		Users: []*User{
			{ID: 101, Username: "john", Password: "password123"},
		},
		Items: []*Item{
			{ID: 1, Name: "Item 1", Price: 10.99, Rating: 4.5, Ratings: []float64{4.0, 5.0}},
			{ID: 2, Name: "Item 2", Price: 25.50, Rating: 3.2, Ratings: []float64{2.5, 3.5, 4.0}},
			{ID: 3, Name: "Item 3", Price: 5.99, Rating: 2.0, Ratings: []float64{1.5, 2.5}},
		},
	}

	system := &System{Database: database}
	// fs := http.FileServer(http.Dir("static"))
	// http.Handle("/", fs)

	// http.HandleFunc("/after", func(w http.ResponseWriter, r *http.Request) {
	// 	name := r.FormValue("username")
	// 	age := r.FormValue("userage")
	// 	fmt.Fprintf(w, "Name: %s; Age: %s", name, age)
	// })
	http.HandleFunc("/register", system.Register)
	http.HandleFunc("/authorize", system.Authorize)
	http.HandleFunc("/items/search", system.SearchItemsByName)
	http.HandleFunc("/items/filter", system.FilterItems)
	http.HandleFunc("/items/rate", system.RateItem)
	fmt.Println("Listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

// curl post request for authorize
// curl -X POST -H "Authorization: Bearer {TOKEN}" -H "Content-Type: application/json" -d '{"Username": "Sala", "password": "123"}' http://localhost:8080/authorize

// curl post request for register
// curl -X POST -H "Content-Type: application/json" -d '{ "Username": "Sala", "Password": "123"}' http://localhost:8080/register

// others by url
