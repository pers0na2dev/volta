package volta

func (a *App) AddQueue(queue ...Queue) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.queues == nil {
		a.queues = make(map[string]Queue)
	}

	for _, q := range queue {
		a.queues[q.Name] = q
	}
}

func (a *App) declareQueue(q Queue) error {
	channel, err := a.baseConnection.Channel()
	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(q.Name, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, nil)
	if err != nil {
		return err
	}

	err = channel.QueueBind(q.Name, q.RoutingKey, q.Exchange, q.NoWait, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) PurgeQueue(name string, noWait bool) error {
	channel, err := a.baseConnection.Channel()
	if err != nil {
		return err
	}

	_, err = channel.QueuePurge(name, noWait)
	if err != nil {
		return err
	}

	return nil
}
