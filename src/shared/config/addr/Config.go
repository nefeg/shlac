package addr

type Config struct{

	Protocol, Address string
}

func (c *Config)Network() string{
	return c.Protocol
}

func (c *Config)String() string{
	return c.Address
}
