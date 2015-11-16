package vanguard

type Inventory struct {
	hosts map[string]Host
}

type Host struct {
	Name string
	Addr string
}

func NewInventory() *Inventory {
	return &Inventory{
		hosts: make(map[string]Host),
	}
}

func (inventory *Inventory) AddHost(host Host) {
	inventory.hosts[host.Name] = host
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
