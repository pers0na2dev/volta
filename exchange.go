package volta

// AddExchanges adds the given exchanges to the application
// If an exchange with the same name already exists, it will be overwritten
func (m *App) AddExchanges(exchange ...Exchange) {
	m.exchangeMutex.Lock()
	defer m.exchangeMutex.Unlock()

	if m.exchanges == nil {
		m.exchanges = make(map[string]Exchange)
	}

	for _, e := range exchange {
		m.exchanges[e.Name] = e
	}
}

// declareExchange declares the given exchange to RabbitMQ
// Internal use only
func (m *App) declareExchange(exchange Exchange) error {
	channel, err := m.baseConnection.Channel()
	if err != nil {
		return err
	}

	err = channel.ExchangeDeclare(exchange.Name, exchange.Type, exchange.Durable, exchange.AutoDelete, exchange.Internal, exchange.NoWait, nil)
	if err != nil {
		return err
	}

	return nil
}

// PurgeExchange purges the given exchange
// If force is true, the exchange will be deleted even if it is in use
// If force is false, the exchange will be deleted only if it is not in use
func (m *App) PurgeExchange(name string, force bool) error {
	channel, err := m.baseConnection.Channel()
	if err != nil {
		return err
	}

	err = channel.ExchangeDelete(name, !force, false)
	if err != nil {
		return err
	}

	return nil
}
