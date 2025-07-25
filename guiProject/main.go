package main

import (
    "context"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
    "encoding/json"
    "log"
)

type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Name     string `json:"name"`
    Age      int    `json:"age"`
    Color    string `json:"color"`
}

var (
    ctx = context.Background()
    rdb *redis.Client // Declare Redis client at package level
    app = fiber.New()
)

func main() {
    initRedis()
	defer rdb.Close()

    setRoutes()

       // Renders our homepage
   app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("./index.html", fiber.Map{
            "Name": "User",
        })
    })


    // ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
    log.Fatal(app.Listen("0.0.0.0:3000"))

}


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

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~


func setRoutes() {

    app.Post("/signup", func(c *fiber.Ctx) error {

        // Get the name from the form
        username := c.FormValue("nusername")
        password := c.FormValue("npassword")

        // Print it server side
        fmt.Printf("signing up %s with password %s\n", username, password)

        
        if (key_exists(username)) {
            fmt.Printf("user exists already")
            return c.Render("./index.html", fiber.Map{
                "warning": "User already exists! Log in instead!",
                "Name": "User",
            })
        } else {
        
            // Add the name to the database
            if err := create_new_user(username, password); err != nil {
                fmt.Printf("adding user failed")
                panic(err)
            }

            // go to login screen
            return mainPage(username, password)(c)
        }

    })



    app.Post("/login", func(c *fiber.Ctx) error {

        fmt.Println("DEBUG: /login route triggered") // Check if this appears in logs

        username := c.FormValue("username")
        password := c.FormValue("password")

        fmt.Printf("DEBUG: username=%s, password=%s\n", username, password) // Verify form data

        if !key_exists(username) {

            fmt.Println("DEBUG: User doesn't exist")

            return c.Render("./index.html", fiber.Map{
                "warning": "Username not found! Sign up!",
                "Name": "User",
            })

        } else {
            value, _ := GetUserData(username, "password");

            if (value != password){
                fmt.Printf("wrong password")

                return c.Render("./index.html", fiber.Map{
                    "warning": "Wrong password, check your spelling",
                    "Name": "User",

                })
            } else {
                fmt.Println("DEBUG: Login successful") // Check if this prints
                return mainPage(username, password)(c)
            }
        }
        

    })


    app.Post("/userInfo", func(c *fiber.Ctx) error {
        username := c.FormValue("username") 
        password := c.FormValue("password") 
        name := c.FormValue("name")
        age := c.FormValue("age")
        color := c.FormValue("color")

        if err := SetUserData(username, "name", name); err != nil {
            fmt.Printf("adding name failed")
            panic(err)
        }
        if err := SetUserData(username, "age", age); err != nil {
            fmt.Printf("adding age failed")
            panic(err)
        }
        if err := SetUserData(username, "color", color); err != nil {
            fmt.Printf("adding color failed")
            panic(err)
        }

        return c.Render("./logged-in.html", fiber.Map{
            "username": username,
            "password": password,
            "name": name,
            "age": age,
            "color": color,
        })
    })

    //
    fmt.Println("Finished route setup")

}

func mainPage(username string, password string) fiber.Handler {
    return func(c *fiber.Ctx) error {

        var userName, _ = GetUserData(username, "name")
        var userAge, _ = GetUserData(username, "age")
        var userColor, _ = GetUserData(username, "color")

        return c.Render("./logged-in.html", fiber.Map{
            "username": username,
            "password": password,
            "name": userName,
            "age": userAge,
            "color": userColor,
        })


        return nil

    }
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