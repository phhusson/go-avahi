package avahi

import (
	"fmt"
	"github.com/godbus/dbus"
)

type AddressResolver struct {
	object       dbus.BusObject
	FoundChannel chan Address
}

func AddressResolverNew(conn *dbus.Conn, path dbus.ObjectPath) (*AddressResolver, error) {
	c := new(AddressResolver)

	c.object = conn.Object("org.freedesktop.Avahi.AddressResolver", path)
	c.FoundChannel = make(chan Address, 10)

	return c, nil
}

func (c *AddressResolver) interfaceForMember(method string) string {
	return fmt.Sprintf("%s.%s", "org.freedesktop.Avahi.AddressResolver", method)
}

func (c *AddressResolver) free() {
	c.object.Call(c.interfaceForMember("Free"), 0)
}

func (c *AddressResolver) GetObjectPath() dbus.ObjectPath {
	return c.object.Path()
}

func (c *AddressResolver) DispatchSignal(signal *dbus.Signal) error {
	if signal.Name == c.interfaceForMember("Found") {
		var address Address
		err := dbus.Store(signal.Body, &address.Interface, &address.Protocol,
			&address.Aprotocol, &address.Address, &address.Name,
			&address.Flags)
		if err != nil {
			return err
		}

		c.FoundChannel <- address
	}

	return nil
}