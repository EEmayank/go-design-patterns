package proxy

import (
	"fmt"
)

// wraps an object to hide some of its characterstics.
// provides an abstraction latyer that is easy to work with and can be changed easilty

type UserFinder interface {
	FinderUser(id uint32) (User, error)
}

type User struct {
	ID int32
}

type UserList []User

func (u *UserList) FindUser(id int32) (User, error) {
	for _, user := range *u {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("users %d could not be found", id)
}

func (u *UserList) addUser(newUser User) {
	*u = append(*u, newUser)
}

type UserListProxy struct {
	SomeDatabase           UserList
	StackCache             UserList
	StackCapacity          int
	DidLastSearchUsedCache bool
}

func (u *UserListProxy) FinderUser(id int32) (User, error) {
	user, err := u.StackCache.FindUser(id)
	if err == nil {
		fmt.Println("Returning user from cache")
		u.DidLastSearchUsedCache = true
		return user, nil
	}
	user, err = u.SomeDatabase.FindUser(id)
	if err != nil {
		return User{}, err
	}
	u.addUserToStack(user)
	fmt.Println("Returning user from database")
	u.DidLastSearchUsedCache = false
	return user, nil
}

func (u *UserListProxy) addUserToStack(user User) {
	if len(u.StackCache) >= u.StackCapacity {
		u.StackCache = append(u.StackCache[1:], user)
	} else {
		u.StackCache.addUser(user)
	}
}
