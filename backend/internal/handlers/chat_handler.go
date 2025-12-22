package handlers

import (
	"log"
	"strconv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"readagain/internal/models"
	"readagain/internal/services"
	ws "readagain/internal/websocket"
)

type ChatHandler struct {
	chatService *services.ChatService
	hub         *ws.Hub
}

func NewChatHandler(chatService *services.ChatService, hub *ws.Hub) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		hub:         hub,
	}
}

// WebSocket Connection Handler
func (h *ChatHandler) HandleWebSocket(c *websocket.Conn) {
	roomID, _ := strconv.ParseUint(c.Params("roomId"), 10, 32)
	userID := c.Locals("userID").(uint)
	username := c.Locals("username").(string)

	// Verify user is member of room
	isMember, err := h.chatService.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		log.Printf("User %d is not a member of room %d", userID, roomID)
		c.Close()
		return
	}

	client := &ws.Client{
		Hub:      h.hub,
		Conn:     c,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		RoomID:   uint(roomID),
		Username: username,
	}

	h.hub.Register <- client

	go client.WritePump()
	client.ReadPump()
}

// REST API Handlers

// CreateRoom creates a new chat room
func (h *ChatHandler) CreateRoom(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req struct {
		Type        string `json:"type" validate:"required,oneof=group direct book_discussion"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		GroupID     *uint  `json:"group_id"`
		BookID      *uint  `json:"book_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	room := &models.ChatRoom{
		Type:        req.Type,
		Name:        req.Name,
		Description: req.Description,
		GroupID:     req.GroupID,
		BookID:      req.BookID,
		CreatedBy:   userID,
		IsActive:    true,
	}

	if err := h.chatService.CreateRoom(room); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create room"})
	}

	// Add creator as admin member
	member := &models.ChatMember{
		RoomID: room.ID,
		UserID: userID,
		Role:   "admin",
	}
	if err := h.chatService.AddMember(member); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add creator as member"})
	}

	return c.Status(fiber.StatusCreated).JSON(room)
}

// GetRoom retrieves a specific room
func (h *ChatHandler) GetRoom(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid room ID"})
	}

	userID := c.Locals("userID").(uint)

	// Verify user is member
	isMember, err := h.chatService.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	room, err := h.chatService.GetRoomByID(uint(roomID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Room not found"})
	}

	return c.JSON(room)
}

// GetUserRooms retrieves all rooms for the authenticated user
func (h *ChatHandler) GetUserRooms(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	rooms, total, err := h.chatService.GetUserRooms(userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch rooms"})
	}

	return c.JSON(fiber.Map{
		"rooms": rooms,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdateRoom updates room details
func (h *ChatHandler) UpdateRoom(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid room ID"})
	}

	userID := c.Locals("userID").(uint)

	// Check if user is admin of the room
	members, err := h.chatService.GetRoomMembers(uint(roomID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify permissions"})
	}

	isAdmin := false
	for _, member := range members {
		if member.UserID == userID && member.Role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins can update room"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.chatService.UpdateRoom(uint(roomID), updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update room"})
	}

	return c.JSON(fiber.Map{"message": "Room updated successfully"})
}

// DeleteRoom deletes a room
func (h *ChatHandler) DeleteRoom(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid room ID"})
	}

	userID := c.Locals("userID").(uint)

	room, err := h.chatService.GetRoomByID(uint(roomID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Room not found"})
	}

	if room.CreatedBy != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only creator can delete room"})
	}

	if err := h.chatService.DeleteRoom(uint(roomID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete room"})
	}

	return c.JSON(fiber.Map{"message": "Room deleted successfully"})
}

// AddMembers adds members to a room
func (h *ChatHandler) AddMembers(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid room ID"})
	}

	userID := c.Locals("userID").(uint)

	// Check if user is admin
	members, err := h.chatService.GetRoomMembers(uint(roomID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify permissions"})
	}

	isAdmin := false
	for _, member := range members {
		if member.UserID == userID && (member.Role == "admin" || member.Role == "moderator") {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Only admins/moderators can add members"})
	}

	var req struct {
		UserIDs []uint `json:"user_ids" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	newMembers := make([]models.ChatMember, len(req.UserIDs))
	for i, uid := range req.UserIDs {
		newMembers[i] = models.ChatMember{
			RoomID: uint(roomID),
			UserID: uid,
			Role:   "member",
		}
	}

	if err := h.chatService.AddMembers(newMembers); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add members"})
	}

	return c.JSON(fiber.Map{"message": "Members added successfully"})
}

// RemoveMember removes a member from a room
func (h *ChatHandler) RemoveMember(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid room ID"})
	}

	memberID, err := strconv.ParseUint(c.Params("memberId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid member ID"})
	}

	userID := c.Locals("userID").(uint)

	// Check if user is admin or removing themselves
	if userID != uint(memberID) {
		members, err := h.chatService.GetRoomMembers(uint(roomID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify permissions"})
		}

		isAdmin := false
		for _, member := range members {
			if member.UserID == userID && (member.Role == "admin" || member.Role == "moderator") {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
		}
	}

	if err := h.chatService.RemoveMember(uint(roomID), uint(memberID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove member"})
	}

	return c.JSON(fiber.Map{"message": "Member removed successfully"})
}

// GetMembers retrieves all members of a room
func (h *ChatHandler) GetMembers(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid room ID"})
	}

	userID := c.Locals("userID").(uint)

	// Verify user is member
	isMember, err := h.chatService.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	members, err := h.chatService.GetRoomMembers(uint(roomID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch members"})
	}

	return c.JSON(members)
}

// GetMessages retrieves messages from a room
func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid room ID"})
	}

	userID := c.Locals("userID").(uint)

	// Verify user is member
	isMember, err := h.chatService.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))

	messages, total, err := h.chatService.GetMessages(uint(roomID), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}

	// Update last read
	h.chatService.UpdateLastRead(uint(roomID), userID)

	return c.JSON(fiber.Map{
		"messages": messages,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// SendMessage sends a message to a room
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	roomID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid room ID"})
	}

	userID := c.Locals("userID").(uint)

	// Verify user is member
	isMember, err := h.chatService.IsMember(uint(roomID), userID)
	if err != nil || !isMember {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	var req struct {
		Message     string `json:"message" validate:"required"`
		MessageType string `json:"message_type"`
		FileURL     string `json:"file_url"`
		FileName    string `json:"file_name"`
		ReplyToID   *uint  `json:"reply_to_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.MessageType == "" {
		req.MessageType = "text"
	}

	message := &models.ChatMessage{
		RoomID:      uint(roomID),
		UserID:      userID,
		Message:     req.Message,
		MessageType: req.MessageType,
		FileURL:     req.FileURL,
		FileName:    req.FileName,
		ReplyToID:   req.ReplyToID,
	}

	if err := h.chatService.CreateMessage(message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send message"})
	}

	// Broadcast via WebSocket
	msg, _ := h.chatService.GetMessageByID(message.ID)
	wsMessage := &ws.Message{
		Type:      "message",
		RoomID:    uint(roomID),
		UserID:    userID,
		MessageID: message.ID,
		Content:   req.Message,
		Data:      msg,
	}
	h.hub.Broadcast <- wsMessage

	return c.Status(fiber.StatusCreated).JSON(message)
}

// UpdateMessage updates a message
func (h *ChatHandler) UpdateMessage(c *fiber.Ctx) error {
	messageID, err := strconv.ParseUint(c.Params("messageId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid message ID"})
	}

	userID := c.Locals("userID").(uint)

	message, err := h.chatService.GetMessageByID(uint(messageID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Message not found"})
	}

	if message.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Can only edit your own messages"})
	}

	var req struct {
		Message string `json:"message" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.chatService.UpdateMessage(uint(messageID), req.Message); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update message"})
	}

	return c.JSON(fiber.Map{"message": "Message updated successfully"})
}

// DeleteMessage deletes a message
func (h *ChatHandler) DeleteMessage(c *fiber.Ctx) error {
	messageID, err := strconv.ParseUint(c.Params("messageId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid message ID"})
	}

	userID := c.Locals("userID").(uint)

	message, err := h.chatService.GetMessageByID(uint(messageID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Message not found"})
	}

	if message.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Can only delete your own messages"})
	}

	if err := h.chatService.DeleteMessage(uint(messageID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete message"})
	}

	return c.JSON(fiber.Map{"message": "Message deleted successfully"})
}

// AddReaction adds a reaction to a message
func (h *ChatHandler) AddReaction(c *fiber.Ctx) error {
	messageID, err := strconv.ParseUint(c.Params("messageId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid message ID"})
	}

	userID := c.Locals("userID").(uint)

	var req struct {
		Emoji string `json:"emoji" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	reaction := &models.ChatReaction{
		MessageID: uint(messageID),
		UserID:    userID,
		Emoji:     req.Emoji,
	}

	if err := h.chatService.AddReaction(reaction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(reaction)
}

// RemoveReaction removes a reaction from a message
func (h *ChatHandler) RemoveReaction(c *fiber.Ctx) error {
	messageID, err := strconv.ParseUint(c.Params("messageId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid message ID"})
	}

	userID := c.Locals("userID").(uint)
	emoji := c.Query("emoji")

	if emoji == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Emoji is required"})
	}

	if err := h.chatService.RemoveReaction(uint(messageID), userID, emoji); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove reaction"})
	}

	return c.JSON(fiber.Map{"message": "Reaction removed successfully"})
}

// GetUnreadCount gets total unread message count for user
func (h *ChatHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	count, err := h.chatService.GetUnreadCount(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get unread count"})
	}

	return c.JSON(fiber.Map{"unread_count": count})
}

// GetOnlineUsers gets list of online users
func (h *ChatHandler) GetOnlineUsers(c *fiber.Ctx) error {
	users := h.hub.GetOnlineUsers()
	return c.JSON(fiber.Map{"online_users": users})
}
