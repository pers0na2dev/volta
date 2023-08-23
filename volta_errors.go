package volta

func (a *App) OnBindError(handler OnBindError) {
	a.onBindError = handler
}
