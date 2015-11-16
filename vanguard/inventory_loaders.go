package vanguard

import (
	"encoding/json"
	"errors"
	"os/exec"
)

func DynamicInventory(name string, args ...string) (*Inventory, error) {
	cmd := exec.Command(name, append(args, "--list")...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, err
	}

	inventory := NewInventory()
	for _, host := range data {
		switch host := host.(type) {
		case []interface{}:
			for _, name := range host {
				_, ok := inventory.GetHost(name.(string))
				if ok {
					continue
				}
				inventory.AddHost(StaticHost(name.(string), ""))
			}
		default:
			return nil, errors.New("invalid dynamic inventory format")
		}
	}
	return inventory, nil
}
