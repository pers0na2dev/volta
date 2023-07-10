package volta

func (m *App) OnMessage(handler ...OnMessage) {
	m.onMessage = append(m.onMessage, handler...)
}
