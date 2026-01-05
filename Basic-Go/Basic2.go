package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

/*
-----------------------------------
CONSTANTS & VARIABLES
-----------------------------------
*/

const AppName = "Go Backend Fundamentals Practice"

var appVersion = "1.0.0"

/*
-----------------------------------
STRUCTS
-----------------------------------
*/

// User represents a basic entity (like DB model)
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Custom error
var ErrUserNotFound = errors.New("user not found")

/*
-----------------------------------
INTERFACE (VERY IMPORTANT)
-----------------------------------
*/

type UserRepository interface {
	Create(user User) (User, error)
	GetByID(id int) (User, error)
	List() []User
}

/*
-----------------------------------
IN-MEMORY REPOSITORY
-----------------------------------
*/

type InMemoryUserRepo struct {
	mu     sync.Mutex
	users  map[int]User
	nextID int
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users:  make(map[int]User),
		nextID: 1,
	}
}

func (r *InMemoryUserRepo) Create(user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	user.CreatedAt = time.Now()

	r.users[user.ID] = user
	r.nextID++

	return user, nil
}

func (r *InMemoryUserRepo) GetByID(id int) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepo) List() []User {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := make([]User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result
}

/*
-----------------------------------
SERVICE LAYER
-----------------------------------
*/

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(name, email string) (User, error) {
	if name == "" || email == "" {
		return User{}, errors.New("name or email cannot be empty")
	}

	user := User{
		Name:  name,
		Email: email,
	}

	return s.repo.Create(user)
}

func (s *UserService) GetUser(ctx context.Context, id int) (User, error) {
	select {
	case <-ctx.Done():
		return User{}, ctx.Err()
	default:
		return s.repo.GetByID(id)
	}
}

/*
-----------------------------------
UTILITY FUNCTIONS
-----------------------------------
*/

// Variadic function
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

/*
-----------------------------------
GOROUTINES & CHANNELS
-----------------------------------
*/

func asyncLogger(ch <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range ch {
		log.Println("ASYNC LOG:", msg)
	}
}

/*
-----------------------------------
MAIN FUNCTION
-----------------------------------
*/

func main() {
	fmt.Println(AppName, "v"+appVersion)

	// Context with timeout (very common in backend)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	repo := NewInMemoryUserRepo()
	service := NewUserService(repo)

	// Channel & goroutine
	logChannel := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go asyncLogger(logChannel, &wg)

	// Create users
	users := []struct {
		name  string
		email string
	}{
		{"Gaurav", "gaurav@example.com"},
		{"Amit", "amit@example.com"},
	}

	for _, u := range users {
		user, err := service.RegisterUser(u.name, u.email)
		if err != nil {
			log.Println("Error:", err)
			continue
		}

		logChannel <- fmt.Sprintf("User created: %+v", user)
	}

	// Get user
	user, err := service.GetUser(ctx, 1)
	if err != nil {
		log.Println("Get user error:", err)
	} else {
		fmt.Println("Fetched User:", user)
	}

	// JSON marshal
	jsonData, err := json.MarshalIndent(repo.List(), "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All Users JSON:")
	fmt.Println(string(jsonData))

	// Use utility function
	fmt.Println("Sum result:", Sum(1, 2, 3, 4, 5))

	// Close channel & wait
	close(logChannel)
	wg.Wait()

	fmt.Println("Program finished cleanly")
}
