package volta

// AddExchanges adds the given exchanges to the application
// If an exchange with the same name already exists, it will be overwritten
func (a *App) AddExchanges(exchange ...Exchange) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.exchanges == nil {
		a.exchanges = make(map[string]Exchange)
	}

	for _, e := range exchange {
		a.exchanges[e.Name] = e
	}
}

// declareExchange declares the given exchange to RabbitMQ
// Internal use only
func (a *App) declareExchange(exchange Exchange) error {
	channel, err := a.baseConnection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.ExchangeDeclare(exchange.Name, exchange.Type, exchange.Durable, exchange.AutoDelete, exchange.Internal, exchange.NoWait, nil)
}

// PurgeExchange purges the given exchange
// If force is true, the exchange will be deleted even if it is in use
// If force is false, the exchange will be deleted only if it is not in use
func (a *App) PurgeExchange(name string, force bool) error {
	channel, err := a.baseConnection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.ExchangeDelete(name, !force, false)
}
