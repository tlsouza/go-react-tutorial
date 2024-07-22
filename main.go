package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tlsouza/go-react-tutorial/src/model"
)

const prefix = "/api/todos"

func main() {
	app := fiber.New()
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("failed to load .env file!")
	}
	PORT := os.Getenv("PORT")

	todos := model.TodoList{}
	//Create a Todo
	app.Post(prefix, func(c *fiber.Ctx) error {
		todo := &model.Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "invalid body for todo task"})
		}

		todo.Id = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(200).JSON(todos)
	})

	//Get todos
	app.Get(prefix, func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	//Update Todos
	app.Patch(prefix+"/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {
				todos[i].Completed = !todos[i].Completed
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "not found"})

	})

	//Delete Todos
	app.Delete(prefix+"/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.Id) == id {

				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "not found"})

	})
	log.Fatal(app.Listen(":" + PORT))
}
