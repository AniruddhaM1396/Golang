// // main.go
// package main

// import (
//  "mflix/handlers"

//  "github.com/gin-gonic/gin"
// )

// func main() {
//  // Connect to MongoDB
//  //models.ConnectToDB()

//  // Set up Gin router
//  router := gin.Default()

//  // Define routes
//  router.GET("/users", handlers.GetUsers)
//  router.POST("/users", handlers.CreateUser)
//  // Add routes for groups, expenses, comments, settlements

//  // Start the server
//  router.Run(":8080")
// }

// Declare the entry point into our application
package main

// Add our dependencies from the standard library, Gin, and MongoDB
import (
    "context"
    "fmt"
    //"log"
    "mflix/models"
    "net/http"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Define your MongoDB connection string
const uri = "mongodb+srv://aniruddha1396:NDB0vq3GIexHhkKS@cluster1.04oimto.mongodb.net/?retryWrites=true&w=majority"

// Create a global variable to hold our MongoDB connection
var mongoClient *mongo.Client

// This function runs before we call our main function and connects to our MongoDB database. If it cannot connect, the application stops.
func init() {
    connect_to_mongodb()
}

// Our entry point into our application
func main() {
    // The simplest way to start a Gin application using the frameworks defaults
    r := gin.Default()

    // Our route definitions
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Hello World",
        })
    })
    r.GET("/users", getUsers)
    r.POST("/addusers", createUser)
    r.GET("/movies/:id", getMovieByID)
    r.POST("/movies/aggregations", aggregateMovies)

    // The Run() method starts our Gin server
    r.Run()
}

// Implemention of the /movies route that returns all of the movies from our movies collection.
func getUsers(c *gin.Context) {
    // Find movies
    cursor, err := mongoClient.Database("User_data").Collection("Users").Find(context.TODO(), bson.D{{}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Map results
    var movies []bson.M
    if err = cursor.All(context.TODO(), &movies); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return movies
    c.JSON(http.StatusOK, movies)
}

func createUser(c *gin.Context) {
    var newUser models.User
    //ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newUser.ID = primitive.NewObjectID()
    fmt.Println(newUser)
    _, err := mongoClient.Database("User_data").Collection("Users").InsertOne(context.Background(), newUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, newUser)
}

// The implementation of our /movies/{id} endpoint that returns a single movie based on the provided ID
func getMovieByID(c *gin.Context) {

    // Get movie ID from URL
    idStr := c.Param("id")

    // Convert id string to ObjectId
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Find movie by ObjectId
    var movie bson.M
    err = mongoClient.Database("sample_mflix").Collection("movies").FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&movie)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return movie
    c.JSON(http.StatusOK, movie)
}

// The implementation of our /movies/aggregations endpoint that allows a user to pass in an aggregation to run our the movies collection.
func aggregateMovies(c *gin.Context) {
    // Get aggregation pipeline from request body
    var pipeline interface{}
    if err := c.ShouldBindJSON(&pipeline); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Run aggregations
    cursor, err := mongoClient.Database("sample_mflix").Collection("movies").Aggregate(context.TODO(), pipeline)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Map results
    var result []bson.M
    if err = cursor.All(context.TODO(), &result); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return result
    c.JSON(http.StatusOK, result)
}

// Our implementation code to connect to MongoDB at startup
func connect_to_mongodb() error {
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

    client, err := mongo.Connect(context.TODO(), opts)
    if err != nil {
        panic(err)
    }
    err = client.Ping(context.TODO(), nil)
    mongoClient = client
    return err
}