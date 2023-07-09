package volta

func (m *App) AddQueue(queue ...Queue) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.queues == nil {
		m.queues = make(map[string]Queue)
	}

	for _, q := range queue {
		m.queues[q.Name] = q
	}
}

func (m *App) declareQueue(q Queue) error {
	channel, err := m.baseConnection.Channel()
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

func (m *App) PurgeQueue(name string, noWait bool) error {
	channel, err := m.baseConnection.Channel()
	if err != nil {
		return err
	}

	_, err = channel.QueuePurge(name, noWait)
	if err != nil {
		return err
	}

	return nil
}
