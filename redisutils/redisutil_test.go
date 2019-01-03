package redisutils

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func Test_first(t *testing.T) {

	InitService(nil)
	rd := GetRedisHdl()

	m := make(map[string]interface{})
	m["k1"] = "v1"

	rd.HSetAll("myhash", m)
	x := rd.HGetAll("myhash")

	for k, v := range x {
		fmt.Printf("key=[%s], value=[%s]\n", k, v)
	}
}

// ping tests connectivity for redisutils (PONG should be returned)
func ping(c redis.Conn) error {
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redisutils.String to convert the interface type to string
	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	set(c)
	get(c)
	setStruct(c)
	getStruct(c)
	return nil
}

// set executes the redisutils SET command
func set(c redis.Conn) error {
	_, err := c.Do("SET", "Favorite Movie", "Repo Man")
	if err != nil {
		fmt.Printf("Error")
		return nil
	}
	_, err = c.Do("SET", "Release Year", 1984)
	if err != nil {
		fmt.Printf("Error")
		return nil
	}
	return nil
}

// get executes the redisutils GET command
func get(c redis.Conn) error {

	// Simple GET example with String helper
	key := "Favorite Movie"
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %s\n", key, s)

	// Simple GET example with Int helper
	key = "Release Year"
	i, err := redis.Int(c.Do("GET", key))
	if err != nil {
		return (err)
	}
	fmt.Printf("%s = %d\n", key, i)

	// Example where GET returns no results
	key = "Nonexistent Key"
	s, err = redis.String(c.Do("GET", key))
	if err == redis.ErrNil {
		fmt.Printf("%s does not exist\n", key)
	} else if err != nil {
		return err
	} else {
		fmt.Printf("%s = %s\n", key, s)
	}

	return nil
}

type User struct {
	Username  string
	MobileID  int
	Email     string
	FirstName string
	LastName  string
}

func setStruct(c redis.Conn) error {

	const objectPrefix string = "user:"

	usr := User{
		Username:  "otto",
		MobileID:  1234567890,
		Email:     "ottoM@repoman.com",
		FirstName: "Otto",
		LastName:  "Maddox",
	}

	// serialize User object to JSON
	json, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	// SET object
	_, err = c.Do("SET", objectPrefix+usr.Username, json)
	if err != nil {
		return err
	}

	return nil
}

func getStruct(c redis.Conn) error {

	const objectPrefix string = "user:"

	username := "otto"
	s, err := redis.String(c.Do("GET", objectPrefix+username))
	if err == redis.ErrNil {
		fmt.Println("User does not exist")
	} else if err != nil {
		return err
	}

	usr := User{}
	err = json.Unmarshal([]byte(s), &usr)

	fmt.Printf("%+v\n", usr)

	return nil

}


func Test_three(t *testing.T) {

	InitService(nil)

	hdl := GetRedisHdl()

	iter := hdl.GetListIterator("x:jt:L_c2")

	for iter.HasNext() {
		key := iter.Next()
		fmt.Printf("KEY: [%s]\n", key)
	}

	miter := hdl.GetMapIterator("a")

	for miter.HasNext() {
		key, value := miter.Next()
		fmt.Printf("KEY: [%s], Value=[%s]\n", key, value)
	}
	{

		miter := hdl.GetMapIterator("x:jt:L_c2")

		for miter.HasNext() {
			key, value := miter.Next()
			fmt.Printf("KEY: [%s], Value=[%s]\n", key, value)
		}
	}
}

func Test_four(t *testing.T) {
	InitService(nil)

	hdl := GetRedisHdl()

	v := hdl.HGet("a", "b")
	fmt.Printf("[%v]\n", v)

}
