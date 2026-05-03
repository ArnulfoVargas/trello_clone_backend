package routes

import (
	"context"
	"math"
	"strconv"

	constants "github.com/ArnulfoVargas/trello_clone_backend.git/cmd/Constants"
	"github.com/ArnulfoVargas/trello_clone_backend.git/cmd/middlewares"
	models_board "github.com/ArnulfoVargas/trello_clone_backend.git/cmd/models/board"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func (r *Router) BindBoardsRoutes() {
	g := r.app.Group("/boards", middlewares.AuthMiddleware)

	g.Get("/", r.getBoards)
	g.Post("/create", r.createBoard)
	g.Put("/update", r.updateName)
	g.Delete("/:id", r.deleteBoard)
}

func (r *Router) getBoards(c fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	page, err := strconv.Atoi(c.Query("page", "1"))

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	countCtx := context.Background()
	maxCount, err := gorm.G[models_board.Board](r.GetConn()).
		Where("user_id = ?", userId).
		Distinct("id").
		Count(countCtx, "id")

	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	if maxCount == 0 {
		return c.JSON(fiber.Map{
			"status": true,
			"body":   []models_board.BoardOut{},
			"total":  0,
			"pages":  0,
		})
	}

	if page <= 0 {
		page = 1
	}

	findContext := context.Background()
	boards, err := gorm.G[models_board.Board](r.GetConn()).
		Where("user_id = ?", userId).
		Order("created_at desc").
		Limit(constants.PageSize).
		Offset((page - 1) * constants.PageSize).
		Find(findContext)

	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	outBoards := make([]models_board.BoardOut, len(boards))

	for i := 0; i < len(boards); i++ {
		b := boards[i]
		outBoards[i] = models_board.BoardOut{
			ID: b.ID,
			Name: b.Name,
		}
	}

	totalPages := int(math.Ceil(float64(maxCount) / float64(constants.PageSize)))

	return c.JSON(fiber.Map{
		"status": true,
		"total":  maxCount,
		"pages":  totalPages,
		"page":   page,
		"body":   outBoards,
	})
}

func (r *Router) createBoard(c fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	var boardIn models_board.BoardIn

	if err := c.Bind().Body(&boardIn); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	board := models_board.Board{
		UserID: userId,
		Name:   boardIn.Name,
	}

	createCtx := context.Background()
	if err := gorm.G[models_board.Board](r.GetConn()).Create(createCtx, &board); err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	out := models_board.BoardOut{
		ID:   board.ID,
		Name: board.Name,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"body":    out,
	})
}

func (r *Router) updateName(c fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	var boardIn models_board.BoardOut

	if err := c.Bind().Body(&boardIn); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	updateCtx := context.Background()
	_, err := gorm.G[models_board.Board](r.GetConn()).
		Where("user_id = ? AND id = ?", userId, boardIn.ID).
		Limit(1).
		Update(updateCtx, "name", boardIn.Name)

	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (r *Router) deleteBoard(c fiber.Ctx) error {
	userId := c.Locals("userID").(string)
	boardIn := c.Params("id", "") 

	if boardIn == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	delCtx := context.Background()

	_, err := gorm.G[models_board.Board](r.GetConn()).
		Where("user_id = ? AND id = ?", userId, boardIn).
		Limit(1).
		Delete(delCtx)
	
	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}

	return c.SendStatus(fiber.StatusOK)
}