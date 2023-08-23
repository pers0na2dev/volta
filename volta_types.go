package volta

type Handler func(*Ctx) error
type OnBindError func(*Ctx, error) error

type Map map[string]interface{}

type Exchange struct {
	Name       string
	Type       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

type Queue struct {
	Name       string
	RoutingKey string
	Exchange   string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}
