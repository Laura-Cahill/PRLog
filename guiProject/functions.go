package main

import (
//    "context"
    "fmt"
//    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
    "encoding/json"
//    "github.com/gofiber/template/html/v2"
//    "github.com/gofiber/fiber/v2/middleware/session"
//    "log"
)

func key_exists(key string) (bool) {

    exists, err := rdb.Exists(ctx, key).Result()
    
    if err != nil {
        panic(err)
    }
    if exists > 0 {
        return true
    } else {
        return false
    }

}


func initRedis() {
    //initialize redis
    rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    //test redis
    if err := rdb.Ping(ctx).Err(); err != nil {
        fmt.Printf("Redis connection failed")
    }
    fmt.Println("Connected to Redis")
}

func create_new_user(username string, password string) error {
    // Create a user object
    newUser := User{
        Username: username,
        Password: password,
        Name:     "",
        Age:      0,
        Color:    "",
    }
    
    // Convert user to JSON string
    userJson, err := json.Marshal(newUser)
    if err != nil {
        fmt.Println("Failed to convert user to JSON:", err)
        return err
    }
    
    // Save to Redis
    err = rdb.Set(ctx, username, userJson, 0).Err()
    if err != nil {
        fmt.Println("Failed to save user to Redis:", err)
        return err
    }
    
    fmt.Println("User created successfully!")
    return nil
}


func GetUserData(username string, dataType string) (interface{}, error) {

    //get the user's data json from redis
    userJson, err := rdb.Get(ctx, username).Result()

    if err != nil {
        fmt.Printf("error retrieve json")
        panic(err)
    }

    //parse JSON into map called userData
    //interface{}: The values can be of any type (Go's way to handle dynamic JSON data)
    var userData map[string]interface{}

    //uses the unmarshal function to fill the userData map and check if error
    if err := json.Unmarshal([]byte(userJson), &userData); err != nil {
        fmt.Printf("json parse error")
        panic(err)
    }

    //return the field with original type
    if value, exists := userData[dataType]; exists {
        return value, nil
    }

    fmt.Printf("data '%s' not found", dataType)
    panic(err)
}


func SetUserData(username string, field string, value interface{}) error {
    //get existing user data from Redis
    userJson, err := rdb.Get(ctx, username).Result()
    if err != nil {
        fmt.Printf("error getting user data")
        panic(err)
    }

    //parse JSON into map
    var userData map[string]interface{}
    if err := json.Unmarshal([]byte(userJson), &userData); err != nil {
        fmt.Printf("error parsing json")
        panic(err)
    }

    //update the specific field
    userData[field] = value

    //convert back to JSON
    updatedJson, err := json.Marshal(userData)
    if err != nil {
        fmt.Printf("error making json")
        panic(err)
    }

    //save back to Redis
    if err := rdb.Set(ctx, username, updatedJson, 0).Err(); err != nil {
        fmt.Printf("error saving json")
        panic(err)
    }

    return nil
}