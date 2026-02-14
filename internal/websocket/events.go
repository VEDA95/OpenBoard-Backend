package websocket

// Event types for board-related WebSocket events
// Topic format: "board:{boardId}"
const (
	EventBoardUpdated = "board.updated"
	EventBoardDeleted = "board.deleted"

	EventListCreated   = "list.created"
	EventListUpdated   = "list.updated"
	EventListDeleted   = "list.deleted"
	EventListReordered = "list.reordered"

	EventCardCreated = "card.created"
	EventCardUpdated = "card.updated"
	EventCardDeleted = "card.deleted"
	EventCardMoved   = "card.moved"

	EventCommentCreated = "comment.created"
	EventCommentUpdated = "comment.updated"
	EventCommentDeleted = "comment.deleted"

	EventLabelCreated = "label.created"
	EventLabelUpdated = "label.updated"
	EventLabelDeleted = "label.deleted"

	EventChecklistItemCreated = "checklist.item.created"
	EventChecklistItemUpdated = "checklist.item.updated"
	EventChecklistItemDeleted = "checklist.item.deleted"
)

// BoardTopic generates the topic string for a board
func BoardTopic(boardID string) string {
	return "board:" + boardID
}
