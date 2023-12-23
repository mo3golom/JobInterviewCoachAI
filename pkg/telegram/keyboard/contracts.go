package keyboard

const (
	ButtonData ButtonType = iota
	ButtonUrl
)

type (
	ButtonType int64

	Button struct {
		Value string
		Type  ButtonType
	}

	InlineButton struct {
		Value string
		Data  []string
		Type  ButtonType
	}

	BuildInlineKeyboardIn struct {
		Command *string
		Buttons []InlineButton
	}

	BuildKeyboardIn struct {
		Buttons []Button
	}

	BuildKeyboardCustomIn struct {
		Buttons [][]Button
	}
)
