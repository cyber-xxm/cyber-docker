package di

import "sync"

type Get func(serviceName string) interface{}

type ServiceConstructor func(get Get) interface{}

type ServiceConstructorMap map[string]ServiceConstructor

type service struct {
	constructor ServiceConstructor
	instance    interface{}
}

type Container struct {
	serviceMap map[string]service
	mutex      sync.RWMutex
}

func NewContainer(serviceConstructors ServiceConstructorMap) *Container {
	c := &Container{
		serviceMap: make(map[string]service),
		mutex:      sync.RWMutex{},
	}
	if serviceConstructors != nil {
		c.Update(serviceConstructors)
	}
	return c
}

func (c *Container) Update(serviceConstructors ServiceConstructorMap) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for serviceName, constructor := range serviceConstructors {
		c.serviceMap[serviceName] = service{
			constructor: constructor,
			instance:    nil,
		}
	}
}

func (c *Container) get(serviceName string) interface{} {
	service, ok := c.serviceMap[serviceName]
	if !ok {
		return nil
	}
	if service.instance == nil {
		service.instance = service.constructor(c.get)
		c.serviceMap[serviceName] = service
	}
	return service.instance
}

func (c *Container) Get(serviceName string) interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.get(serviceName)
}
