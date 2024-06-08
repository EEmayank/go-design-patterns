package proxy

import (
	"math/rand"
	"testing"
)

func Test_UserListProxy(t *testing.T) {
	someDatabase := UserList{}

	// creating mock database with one million users
	for range 1000000 {
		n := rand.Int31()
		someDatabase = append(someDatabase, User{ID: n})
	}

	proxy := UserListProxy{
		SomeDatabase:  someDatabase,
		StackCache:    UserList{},
		StackCapacity: 2,
	}

	// getting 3 random user id from mock data
	knownIds := [3]int32{
		someDatabase[3].ID,
		someDatabase[4].ID,
		someDatabase[5].ID,
	}

	// trying to find user with an empty cache
	t.Run("FindUser = Empty Cache", func(t *testing.T) {
		user, err := proxy.FinderUser(knownIds[0])
		if err != nil {
			t.Fatal(err)
		}
		if user.ID != knownIds[0] {
			t.Error("Returned user doesn't match with expected")
		}
		if len(proxy.StackCache) != 1 {
			t.Error("After one successful search in an empty cache, the of stack should be one")
		}
		if proxy.DidLastSearchUsedCache {
			t.Error("No user can be returned from an empty cache")
		}
	})

	// asking for same user as before which now must be returned from the cache as the user must've been set in the cache
	// after the first search
	t.Run("FindUser = One User, ask for the same user", func(t *testing.T) {
		user, err := proxy.FinderUser(knownIds[0])
		if err != nil {
			t.Fatal(err)
		}
		if user.ID != knownIds[0] {
			t.Error("Returned user doesn't match with expected")
		}
		if len(proxy.StackCache) != 1 {
			t.Error("After one successful search in an empty cache, the of stack should be one")
		}
		if !proxy.DidLastSearchUsedCache {
			t.Error("The user should be returned from the cache")
		}
	})

	// overflowing the stackcache
	user1, err := proxy.FinderUser(knownIds[0])
	if err != nil {
		t.Fatal(err)
	}
	user2, err := proxy.FinderUser(knownIds[1])
	if err != nil {
		t.Fatal(err)
	}
	if proxy.DidLastSearchUsedCache {
		t.Error("The user2 wasn't stored on the proxy cache yet")
	}

	user3, err := proxy.FinderUser(knownIds[2])
	if err != nil {
		t.Fatal(err)
	}
	if proxy.DidLastSearchUsedCache {
		t.Error("The user3 wasn't stored on the proxy cache yet")
	}

	for i := range len(proxy.StackCache) {
		if proxy.StackCache[i].ID == user1.ID {
			t.Error("User that should be gone was found")
		}
		if len(proxy.StackCache) != proxy.StackCapacity {
			t.Errorf("After inserting %d users cache should not exceed the stack capacity of %d", len(proxy.StackCache), proxy.StackCapacity)
		}
	}

	for _, v := range proxy.StackCache {
		if v != user2 && v != user3 {
			t.Error("A non expected user was found on the cache")
		}
	}
}
