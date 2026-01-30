package handlers

import (
	"crud/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.CreateUser

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, "Invalid json body", http.StatusBadRequest)
			return
		}

		// Input validation

		if len(user.Name) < 1 {
			err := models.ReturnApiUser{
				Success: false,
				Message: "name field shouldn't be empty",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		if len(user.Email) < 6 {
			err := models.ReturnApiUser{
				Success: false,
				Message: "email field shouldn't be empty",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return

		}
		var UserExists bool
		query := `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);
	   `

		err = db.QueryRow(query, user.Email).Scan(&UserExists)

		if err != nil {
			http.Error(w, "Internal server error, try later", http.StatusInternalServerError)
			return
		}

		if UserExists {
			err := models.ReturnApiUser{
				Success: false,
				Message: "user with this email already exists, choose different email",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}

		if len(user.Password) < 6 {
			err := models.ReturnApiUser{
				Success: false,
				Message: "password must be atleast 6charcter long",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return

		}

		//Hash password(Good Security)

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		//Insert data into database
		query = `
		INSERT INTO users (name,email,password) VALUES ($1,$2,$3);
	   `

		_, err = db.Exec(query, user.Name, user.Email, string(hashed))

		if err != nil {
			http.Error(w, "Internal server error, try later", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(models.ReturnApiUser{
			Success: true,
			Message: "Account created successfully",
			Data: []models.User{
				{
					Name:     user.Name,
					Email:    user.Email,
					Password: string(hashed),
				},
			},
		})

	}
}

func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil {
			err := models.ReturnApiUser{
				Success: false,
				Message: "Id must be a number",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err)
			return
		}
		query := `
	    SELECT id,name,created_at FROM users WHERE id = $1 LIMIT 1;
	  `
		err = db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.CreatedAt)

		if err != nil {
			err := models.ReturnApiUser{
				Success: false,
				Message: "user with this id doesn't exists",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(err)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.ReturnApiUser{
			Success: true,
			Message: "user fetched successfully",
			Data:    []models.User{user},
		})

	}
}

func GetAllUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data []models.User
		var users models.User
		query := `SELECT id,name,created_at FROM users`

		row, err := db.Query(query)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error, try later", http.StatusInternalServerError)
			return
		}
		defer row.Close()

		for row.Next() {
			err = row.Scan(&users.Id, &users.Name, &users.CreatedAt)

			if err != nil {
				log.Println(err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(models.ReturnApiUser{
					Success: false,
					Message: "Failed to fetch users",
				})
				return
			}

			data = append(data, users)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.ReturnApiUser{
			Success: true,
			Message: "Users fetched successfully",
			Data:    data,
		})

	}
}

func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			w.Header().Set("content-tyoe", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "Invalid incoming json",
			})
			return
		}

		if user.Email == "" {
			w.Header().Set("content-tyoe", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "Email field shouldn't be empty",
			})
			return
		}

		if user.Password == "" {
			w.Header().Set("content-tyoe", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "Password field shouldn't be empty",
			})
			return
		}

		query := `SELECT name,password FROM users WHERE email = $1`

		var hashed string

		err = db.QueryRow(query, user.Email).Scan(&user.Name, &hashed)

		if err != nil {
			w.Header().Set("content-tyoe", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "Invalid email entered",
			})
			log.Println(err)
			return

		}

		//Check password validty

		err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(user.Password))

		if err != nil {
			w.Header().Set("content-tyoe", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "Invalid password entered",
			})
			return
		}

		w.Header().Set("content-tyoe", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.ReturnApiUser{
			Success: true,
			Message: "Welcome Home," + user.Name,
		})

	}
}

func UpdateEmail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var email models.UpdateEmail

		err := json.NewDecoder(r.Body).Decode(&email)

		if err != nil {
			http.Error(w, "Invalid incoming json", http.StatusBadRequest)
			return
		}

		if email.OldEmail == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "old_email field shouldn't be empty",
			})
			return
		}

		var exists bool

		query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

		err = db.QueryRow(query, email.OldEmail).Scan(&exists)

		if !exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "Given old_email address doesn't exists",
			})
			return
		}

		if email.NewEmail == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "new_email field shouldn't be empty",
			})
			return
		}

		if exists {
			query = `UPDATE users SET email = $1 WHERE email= $2`
			_, err = db.Exec(query, email.NewEmail, email.OldEmail)

			if err != nil {
				http.Error(w, "Internal server error, try later", http.StatusInternalServerError)
				log.Println(err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: true,
				Message: "Email has been changed successfully",
				Data: []models.User{
					{
						OldEmail: email.OldEmail,
						NewEmail: email.NewEmail,
					},
				},
			})
		}
	}
}

func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, "Invalid incoming json", http.StatusBadRequest)
			return
		}

		if user.Email == "" {
			w.Header().Set("content-tyoe", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "Email field shouldn't be empty",
			})
			return
		}

		var exists bool

		query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

		err = db.QueryRow(query, user.Email).Scan(&exists)

		if !exists {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ReturnApiUser{
				Success: false,
				Message: "Given email address doesn't exists",
			})
			return
		}

		query = `DELETE FROM users WHERE email = $1`

		_, err = db.Exec(query, user.Email)

		if err != nil {
			http.Error(w, "Internal server error, try later", http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("content-tyoe", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.ReturnApiUser{
			Success: true,
			Message: "User with this email has been deleted successfully",
		})

	}
}
