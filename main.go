package main

import (
	"log"
	"net/http"
	"user_service/pkg/storage"
	"user_service/src/handler"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()
	currentStorage := storage.MakeStorage()

	log.Println("Run app")

	// хендлер создания пользователя
	router.Post("/create", handler.CreateUser(currentStorage))

	// хендлер делает друзей из двух пользователей
	router.Post("/make_friends", handler.MakeFriends(currentStorage))

	// хендлер удаляет пользователя
	router.Delete("/user", handler.DeleteUser(currentStorage))

	router.Get("/friends/{user_id}", handler.GetAllFriends(currentStorage))

	// хендлер обновляет возраст пользователя§
	router.Put("/user/{user_id}", handler.UpdateAge(currentStorage))

	// хендлер вывводит всех пользователей
	router.Get("/get_users", handler.GetAllUsers(currentStorage))

	// заглушка
	router.Get("/", handler.Get())

	log.Println(http.ListenAndServe("localhost:8080", router))
}
