package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"user_service/pkg/storage"

	"github.com/go-chi/chi"
)

// Get - Заглушка
func Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("root")); err != nil {
			log.Fatalln(err)
		}
	}
}

// GetAllUsers  - http хендлер выводит всех пользователей
func GetAllUsers(s *storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		users := s.GetAllUsers()

		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Fatalln(err)
		}
	}
}

// CreateUser - http хендлер создания пользователя
func CreateUser(s *storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(r.Body)

		request := map[string]interface{}{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Fatalln(err)
		}

		name := fmt.Sprintf("%v", request["name"])
		age, err := strconv.Atoi(fmt.Sprintf("%v", request["age"]))
		if err != nil {
			log.Fatalln(err)
		}
		var friends []int
		if request["friends"] != nil {
			friends = convertSlice(request["friends"].([]interface{}))
		}

		userID := s.AddUser(name, age, friends)
		response := fmt.Sprintf("User %s with ID %d was created\n", name, userID)

		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write([]byte(response)); err != nil {
			log.Fatalln(err)
		}
	}
}

// MakeFriends - http хендлер создания дружеских связей
func MakeFriends(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(r.Body)

		request := map[string]string{}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Fatalln(err)
		}

		id1, err := strconv.Atoi(request["source_id"])
		if err != nil {
			log.Fatalln(err)
		}
		id2, err := strconv.Atoi(request["target_id"])
		if err != nil {
			log.Fatalln(err)
		}

		s.MakeFriends(id1, id2)
		response := fmt.Sprintf("%s and %s are friends now\n", s.GetUser(id1).GetName(), s.GetUser(id2).GetName())
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(response)); err != nil {
			log.Fatalln(err)
		}

		return
	}
}

// DeleteUser - http хендлер удаления пользователя
func DeleteUser(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(r.Body)

		request := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Fatalln(err)
		}

		id, err := strconv.Atoi(request["target_id"])
		if err != nil {
			log.Fatalln(err)
		}

		response := fmt.Sprintf("User %s is deleted\n", s.GetUser(id).GetName())

		s.DeleteFromFriends(id)
		s.DeleteUser(id)

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(response)); err != nil {
			log.Fatalln(err)
		}
	}
}

// GetAllFriends -  http хендлер получения друзей пользователя
func GetAllFriends(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		val := chi.URLParam(r, "user_id")

		userID, err := strconv.Atoi(val)

		if err != nil {
			body := MakeBody()
			w.WriteHeader(http.StatusNotFound)
			body.SetMessage("Идентификатор не корректен")

			if err := json.NewEncoder(w).Encode(body); err != nil {
				log.Fatalln(err)
			}
			return
		}

		friendsID := s.GetFriendsID(userID)

		if len(friendsID) < 1 {
			body := MakeBody()
			w.WriteHeader(http.StatusOK)
			body.SetMessage("Нет друзей")

			if err := json.NewEncoder(w).Encode(body); err != nil {
				log.Fatalln(err)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(friendsID); err != nil {
			log.Fatalln(err)
		}

	}
}

// UpdateAge - http хендлер обновления возраста пользователя
func UpdateAge(s *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}(r.Body)

		val := chi.URLParam(r, "user_id")

		userID, err := strconv.Atoi(val)
		if err != nil {
			body := MakeBody()
			w.WriteHeader(http.StatusNotFound)
			body.SetMessage("Идентификатор не корректен")

			if err := json.NewEncoder(w).Encode(body); err != nil {
				log.Fatalln(err)
			}
		}

		request := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Fatalln(err)
		}

		age, err := strconv.Atoi(request["new age"])
		if err != nil {
			body := MakeBody()
			w.WriteHeader(http.StatusOK)
			body.SetMessage("Возраст не корректен")

			if err := json.NewEncoder(w).Encode(body); err != nil {
				log.Fatalln(err)
			}
			return
		}

		s.UpdateAge(userID, age)
		response := fmt.Sprintf("User %d's age has been updated to %d\n", userID, age)
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte(response)); err != nil {
			log.Fatalln(err)
		}

	}
}

// convertSlice - helper funct
func convertSlice(in []interface{}) (out []int) {
	b := make([]int, len(in))
	for i := range in {
		b[i] = in[i].(int)
	}
	return
}
