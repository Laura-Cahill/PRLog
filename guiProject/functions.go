package main

import (
//    "context"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/redis/go-redis/v9"
    "encoding/json"
    "time"
//    "github.com/gofiber/template/html/v2"
//    "github.com/gofiber/fiber/v2/middleware/session"
//    "log"
)

//-----------------------------------------------------

//initalize redis for the start of the program
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

//check if key exists, should likely be used for users only
func key_exists(key string) (bool) {

    //checks if key exists
    exists, err := rdb.Exists(ctx, key).Result()
    
    //if checking failed, panic
    if err != nil {
        panic(err)
    }
    //return if it exists or not
    if exists > 0 {
        return true
    } else {
        return false
    }

}

//create a new user and add them to the database
func create_new_user(username string, password string) error {

    currentTime := time.Now()
	var date = currentTime.Format("01-02-2006")

    // Create a user object
    newUser := User{
        Username: username,
        Password: password,
        Name:     "",
        Age:      "0",
        Color:    "",
        Logged: "0",
        Since: date,
        Email: "",
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

//get user data frrom their username, and the goal data type (like age or email)
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

//set user data with username, data type, and the data
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

//get the username of the current session. call with getUser(c) 
func getUser(c *fiber.Ctx) string {

    sess, err := store.Get(c)
    if err != nil {
        panic(err)
    }
        
    //get username from session and turn to string
    username := sess.Get("username")

    if username == nil {
        return ""
    } else {
        return username.(string)
    }


}

//adds workout to a user's data. int params are: sets, reps, pounds
//bug: only one workout with a a workout name (squats) can be entered at once
func AddWorkout(username string, date string, exerciseName string, params [3]int) error {

    //get full existing user data
    userData, err := rdb.Get(ctx, username).Result()

    if err != nil {
        panic(err)
    }

    //change into User struct
    var user User

    //try to unmarshall userData into user
    if err := json.Unmarshal([]byte(userData), &user); err != nil {
        panic(err)
    }

    //initialize empty Workouts item if nil
    if user.Workouts == nil {

        user.Workouts = make(map[string]map[string]interface{}) //Workouts: [map of ids paired with workout maps]

    }

    //create id with date
    workoutID := "workout_" + date

    //if this workout date id is not mapped to any value
    if user.Workouts[workoutID] == nil {

        user.Workouts[workoutID] = make(map[string]interface{}) //change value to be a map of interfaces values paired with string keys

        user.Workouts[workoutID]["date"] = date //the "date:" is mapped to a string of the current date. the rest is left empty.

    }

    //add exercise
    //bug to fix: if exercise name exists, it will be overridden
    user.Workouts[workoutID][exerciseName] = params

    //make it back into json
    userJson, err := json.Marshal(user)
    if err != nil {
        panic(err)
    }
    
    //set to database and return errors
    return rdb.Set(ctx, username, userJson, 0).Err()
}

//get a workout from a user with a specific date
func GetWorkout(username string, date string) (map[string]interface{}, error) {
    //get user data from Redis
    userData, err := rdb.Get(ctx, username).Result()

    //if key not found:
    if err == redis.Nil {

        return nil, fmt.Errorf("user not found")

    } else if err != nil { //if something else is the issue
        panic(err)
    }

    //unmarshall json into a map
    var user map[string]interface{}
    if err := json.Unmarshal([]byte(userData), &user); err != nil {
        panic(err)
    }

    //check if workouts exist
    if user["workouts"] == nil {
        return nil, fmt.Errorf("no workouts found for user")
    }

    workouts := user["workouts"].(map[string]interface{})

    //look up workout with date
    workoutID := "workout_" + date
    if workouts[workoutID] == nil {

        return nil, fmt.Errorf("no workout found on %s", date)

    }

    //create workout map
    workout := workouts[workoutID].(map[string]interface{})

    //return the workout map 
    return workout, nil
}


// ------------------------------------




func GetWorkoutsBetweenDates(username, startDate, endDate, exerciseName string) ([]WorkoutData, error) {

    current, _ := time.Parse("2006-01-02", startDate)
    end, _ := time.Parse("2006-01-02", endDate)

    var workouts []WorkoutData
    
    //oterate through each day in range
    for !current.After(end) {

        dateStr := current.Format("2006-01-02")

        workout, err := GetWorkout(username, dateStr)

        if err == nil && workout[exerciseName] != nil {

            params := workout[exerciseName].([]interface{})
            workouts = append(workouts, WorkoutData{

                Date:   dateStr,
                Sets:   int(params[0].(float64)),
                Reps:   int(params[1].(float64)),
                Weight: int(params[2].(float64)),

            })
        }

        current = current.AddDate(0, 0, 1) //next day
    }

    return workouts, nil
}

//workoutData represents a single workout record for graphing
type WorkoutData struct {
    Date   string `json:"date"`
    Sets   int    `json:"sets"`
    Reps   int    `json:"reps"`
    Weight int    `json:"weight"`
}


