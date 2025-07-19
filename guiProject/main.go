package main

import (
    "context"
    "fmt"

    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
    app := fiber.New()

    // Renders our homepage
    app.Get("/", func(c *fiber.Ctx) error {
        return c.Render("./index.html", fiber.Map{
            "Name": "Ok",
        })
    })

    // Post route for form submission
    app.Post("/submit", func(c *fiber.Ctx) error {

        // Get the name from the form
        name := c.FormValue("name")

        // Print it server side
        fmt.Printf("some idiot said her name was %s\n", name)

        // Add the name to the database
        if err := add_to_db(name); err != nil {
            panic(err)
        }

        // Get the name from the database
        if val, err := get_from_db(); err != nil {
            panic(err)
        } else {
            // Return a document with our data
            return c.Render("./index.html", fiber.Map{
                "Name": val,
            })
        }

    })

    app.Listen("0.0.0.0:3000")

}

// Add data from a form to the database
func add_to_db(data string) error {

    // Make a new redis connection
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    // Insert the user's data into the database
    err := rdb.Set(ctx, "key", data, 0).Err()
    if err != nil {
        panic(err)
    }

    return nil

}

// Get data from the database and display on a webpage
func get_from_db() (string, error) {

    // Make a new redis connection
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    // Get a datum from the database
    if val, err := rdb.Get(ctx, "key").Result(); err != nil {
        panic(err)
    } else {
        return val, nil
    }
}