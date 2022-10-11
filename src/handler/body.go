package handler

type Body struct {
	Message string
}

// MakeBody - возвращает объект
func MakeBody() *Body {

	return &Body{}
}

// SetMessage - задает текс сообщения
func (b *Body) SetMessage(text string) {
	b.Message = text
}
