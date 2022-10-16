package storage

import (
	"user_service/pkg/user"
)

var userId int = 1

// Storage - хранилище пользователй
type Storage struct {
	store map[int]*user.User
}

// MakeStorage - создает экземпляр хранилище пользователей
func MakeStorage() *Storage {

	return &Storage{
		make(map[int]*user.User),
	}
}

// AddUser - добавление пользователя
func (s *Storage) AddUser(name string, age int, friends []int) int {

	if name != "" {

		newUser := user.MakeUser()
		newUser.SetName(name)
		newUser.SetAge(age)
		newUser.SetFriends(friends)
		userId := len(s.store) + 1

		s.store[userId] = newUser
	}
	return userId
}

// UpdateAge - обновление возраста
func (s *Storage) UpdateAge(id, age int) {
	s.store[id].SetAge(age)
}

// DeleteUser - удаление пользователя
func (s *Storage) DeleteUser(id int) {
	delete(s.store, id)
}

// GetAllUsers - вывод всех пользователей
func (s *Storage) GetAllUsers() map[int]*user.User {
	return s.store
}

// GetUser - вывод пользователя
func (s *Storage) GetUser(id int) *user.User {
	return s.store[id]
}

// GetFriendsID - получение id друзей
func (s *Storage) GetFriendsID(id int) []int {

	targetUser, ok := s.store[id]

	if ok {
		return targetUser.GetFriends()
	}

	return make([]int, 0)
}

// DeleteFromFriends - удаляет дружеские связи
func (s *Storage) DeleteFromFriends(id int) {
	for userID, userObj := range s.store {
		if userID != id {
			userObj.DeleteFriend(id)
		}
	}
}

// MakeFriends - создает дружеские связи
func (s *Storage) MakeFriends(a int, b int) {
	s.store[a].AddFriend(b)
	s.store[b].AddFriend(a)
}
