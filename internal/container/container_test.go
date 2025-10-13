package container_test

import (
	"testing"

	"github.com/kadsin/sms-gateway/internal/container"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	ID int
}

type Bar interface {
	GetName() string
}

type barImpl struct {
	name string
}

func (b *barImpl) GetName() string {
	return b.name
}

func Test_Container_RegisterAndResolveSingleton(t *testing.T) {
	container.Flush()

	counter := 0
	container.RegisterSingleton(func() *Foo {
		counter++
		return &Foo{ID: counter}
	})

	f1 := container.Resolve[*Foo]()
	f2 := container.Resolve[*Foo]()

	require.Equal(t, f1, f2)
	require.Equal(t, 1, f1.ID)
	require.Equal(t, 1, f2.ID)
}

func Test_Container_RegisterAndResolveTransient(t *testing.T) {
	container.Flush()

	counter := 0
	container.RegisterTransient(func() *Foo {
		counter++
		return &Foo{ID: counter}
	})

	f1 := container.Resolve[*Foo]()
	f2 := container.Resolve[*Foo]()

	require.NotEqual(t, f1, f2)
	require.Equal(t, 1, f1.ID)
	require.Equal(t, 2, f2.ID)
}

func Test_Container_RegisterSingletonAsInterface(t *testing.T) {
	container.Flush()

	container.RegisterSingletonAs[Bar](func() *barImpl {
		return &barImpl{name: "bar-singleton"}
	})

	b1 := container.Resolve[Bar]()
	b2 := container.Resolve[Bar]()

	require.Equal(t, "bar-singleton", b1.GetName())
	require.Equal(t, b1, b2)
}

func Test_Container_RegisterTransientAsInterface(t *testing.T) {
	container.Flush()

	container.RegisterTransientAs[Bar](func() *barImpl {
		return &barImpl{name: "bar-transient"}
	})

	b1 := container.Resolve[Bar]()
	b2 := container.Resolve[Bar]()

	require.Equal(t, "bar-transient", b1.GetName())
	require.NotSame(t, b1, b2)
}

func Test_Container_ResolveUnregisteredServicePanics(t *testing.T) {
	container.Flush()

	require.Panics(t, func() {
		container.Resolve[*Foo]()
	})
}

func Test_Container_PanicIfNotImplemented(t *testing.T) {
	container.Flush()

	type NotImpl struct{}

	require.Panics(t, func() {
		container.RegisterSingletonAs[Bar](func() *NotImpl { return &NotImpl{} })
	})
}
