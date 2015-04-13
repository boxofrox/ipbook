package registry

import ()

type Registry struct {
	register map[string]string
}

func New() *Registry {
	return &Registry{
		register: make(map[string]string),
	}
}

func (r *Registry) Contains(name string) bool {
	_, exists := r.register[name]
	return exists
}

func (r *Registry) Put(name, ip string) error {
	if "" == name {
		return INVALID_NAME
	}

	if "" == ip {
		delete(r.register, name)
	} else {
		r.register[name] = ip
	}

	return nil
}

func (r *Registry) Get(name string) (string, error) {
	if "" == name {
		return "", INVALID_NAME
	}

	if _, exists := r.register[name]; false == exists {
		return "", UNKNOWN_NAME
	}

	return r.register[name], nil
}

type RegistryError int

const (
	INVALID_NAME RegistryError = iota
	UNKNOWN_NAME
)

func (err RegistryError) Error() string {
	switch err {
	case INVALID_NAME:
		return "invalid name"

	case UNKNOWN_NAME:
		return "unknown name"
	}

	return "unknown error"
}
