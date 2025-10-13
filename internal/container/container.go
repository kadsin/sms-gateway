package container

import (
	"fmt"
	"reflect"
	"sync"
)

type lifetime int

const (
	lifetimeSingleton lifetime = iota
	lifetimeTransient
)

type registeredService struct {
	lifetime lifetime
	service  any
}

var (
	services = map[reflect.Type]*registeredService{}
	mu       sync.RWMutex
)

func RegisterTransient[T any](factory func() T) {
	register(lifetimeTransient, getTypeOf[T](), factory)
}

func RegisterSingleton[T any](factory func() T) {
	register(lifetimeSingleton, getTypeOf[T](), factory())
}

func RegisterTransientAs[I any, T any](factory func() T) {
	interfaceType := getTypeOf[I]()
	panicIfNotImplemented[T](interfaceType)

	register(lifetimeTransient, interfaceType, factory)
}

func RegisterSingletonAs[I any, T any](factory func() T) {
	interfaceType := getTypeOf[I]()
	panicIfNotImplemented[T](interfaceType)

	register(lifetimeSingleton, interfaceType, factory())
}

func panicIfNotImplemented[T any](interfaceType reflect.Type) {
	implementedType := getTypeOf[T]()
	if !implementedType.Implements(interfaceType) {
		panic(fmt.Sprintf("container: type %v does not implement %v", implementedType, interfaceType))
	}
}

func register(l lifetime, t reflect.Type, service any) {
	mu.Lock()
	defer mu.Unlock()

	services[t] = &registeredService{
		lifetime: l,
		service:  service,
	}
}

func Resolve[T any]() T {
	t := getTypeOf[T]()

	mu.RLock()
	svc, ok := services[t]
	mu.RUnlock()

	if !ok {
		panic(fmt.Sprintf("container: service not registered for type %v", t))
	}

	switch svc.lifetime {
	case lifetimeTransient:
		instance := reflect.ValueOf(svc.service).Call(nil)[0].Interface().(T)
		if !ok {
			panic(fmt.Sprintf("container: factory did not return %v for %v", t, t))
		}

		return instance

	case lifetimeSingleton:
		instance, ok := svc.service.(T)
		if !ok {
			panic(fmt.Sprintf("container: invalid singleton type for %v", t))
		}

		return instance

	default:
		panic("container: unknown service lifetime")
	}
}

func getTypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func Flush() {
	mu.Lock()
	defer mu.Unlock()

	services = map[reflect.Type]*registeredService{}
}
