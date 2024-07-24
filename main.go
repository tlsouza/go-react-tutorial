package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tlsouza/go-react-tutorial/repository"
	"github.com/tlsouza/go-react-tutorial/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const prefix = "/api/todos"

func main() {
	app := fiber.New()
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("failed to load .env file!")
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	client := repository.Connect()
	defer client.Disconnect(context.Background())

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodos)

	log.Fatal(app.Listen("0.0.0.0:" + port))

}

func getTodos(c *fiber.Ctx) error {
	var todos model.TodoList

	cursor, err := repository.Collection.Find(context.Background(), bson.M{})

	if err != nil {
		return nil
	}

	for cursor.Next(context.Background()) {
		var todo model.Todo

		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)

	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(model.Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Body"})
	}

	insertedObject, err := repository.Collection.InsertOne(context.Background(), todo)

	if err != nil {
		return err
	}

	todo.Id = insertedObject.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(todo)

}

func updateTodo(c *fiber.Ctx) error {

	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Error": "Invalid ID"})
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"completed": true}}
	_, err = repository.Collection.UpdateOne(c.Context(), filter, update)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})

}

func deleteTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"Error": "Invalid ID"})
	}

	filter := bson.M{"_id": objectId}
	_, err = repository.Collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
