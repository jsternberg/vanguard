package vanguard

type Inventory struct {
	hosts map[string]Host
}

type Host interface {
	Name() string
	Addr() string
}

type staticHost struct {
	name string
	addr string
}

func StaticHost(name string, addr string) Host {
	return &staticHost{name: name, addr: addr}
}

func (h *staticHost) Name() string {
	return h.name
}

func (h *staticHost) Addr() string {
	return h.addr
}

func NewInventory() *Inventory {
	return &Inventory{
		hosts: make(map[string]Host),
	}
}

func (inventory *Inventory) AddHost(host Host) {
	inventory.hosts[host.Name()] = host
}

func (inventory *Inventory) GetHost(name string) (Host, bool) {
	host, ok := inventory.hosts[name]
	return host, ok
}

func (inventory *Inventory) Hosts() []string {
	hosts := make([]string, 0, len(inventory.hosts))
	for name := range inventory.hosts {
		hosts = append(hosts, name)
	}
	return hosts
}
