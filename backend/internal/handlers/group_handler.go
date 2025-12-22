package handlers

import (
	"readagain/internal/models"
	"readagain/internal/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type GroupHandler struct {
	service *services.GroupService
}

func NewGroupHandler(service *services.GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}

func (h *GroupHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	groups, total, err := h.service.GetAll(page, limit, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch groups"})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(fiber.Map{
		"groups":       groups,
		"total":        total,
		"page":         page,
		"limit":        limit,
		"total_pages":  totalPages,
	})
}

func (h *GroupHandler) GetByID(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	group, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Group not found"})
	}
	return c.JSON(group)
}

func (h *GroupHandler) Create(c *fiber.Ctx) error {
	var input struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	userIDInterface := c.Locals("user_id")
	if userIDInterface == nil {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	
	userID, ok := userIDInterface.(uint)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	group := &models.Group{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   userID,
	}

	if err := h.service.Create(group); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create group"})
	}

	return c.Status(201).JSON(group)
}

func (h *GroupHandler) Update(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	updates := make(map[string]interface{})
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}

	if err := h.service.Update(uint(id), updates); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update group"})
	}

	return c.JSON(fiber.Map{"message": "Group updated successfully"})
}

func (h *GroupHandler) Delete(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete group"})
	}

	return c.JSON(fiber.Map{"message": "Group deleted successfully"})
}

func (h *GroupHandler) GetMembers(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	members, err := h.service.GetMembers(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch members"})
	}
	return c.JSON(members)
}

func (h *GroupHandler) AddMember(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var input struct {
		UserID uint `json:"user_id" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.AddMember(uint(id), input.UserID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Member added successfully"})
}

func (h *GroupHandler) RemoveMember(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID, _ := strconv.Atoi(c.Params("userId"))

	if err := h.service.RemoveMember(uint(id), uint(userID)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to remove member"})
	}

	return c.JSON(fiber.Map{"message": "Member removed successfully"})
}

func (h *GroupHandler) AddMembers(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var input struct {
		UserIDs []uint `json:"user_ids" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.AddMembers(uint(id), input.UserIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add members"})
	}

	return c.JSON(fiber.Map{"message": "Members added successfully"})
}
