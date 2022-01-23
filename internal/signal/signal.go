package signal

import "reflect"

type signal string

const (
	PostCreate signal = "PostCreate"
	PostDelete signal = "PostDelete"
)

type signalConfig struct {
	sender    interface{}
	callbacks []signalCallback
}

var signalMapping = map[string]signalConfig{}

type signalCallback func(name string, instance interface{}, args ...interface{})

type Signal interface {
	Connect(s signal, sender interface{}, f signalCallback)
	Emit(s signal, sender interface{}, instance interface{})
}

func (s signal) Connect(sender interface{}, f signalCallback) {
	name := reflect.TypeOf(sender).String()

	// Set the signal sender and append the callback
	callbacks := append(signalMapping[name].callbacks, f)
	signalMapping[name] = signalConfig{sender: sender, callbacks: callbacks}
}

func (s signal) Emit(sender interface{}, instance interface{}, args ...interface{}) {
	name := reflect.TypeOf(sender).String()
	c := signalMapping[name]

	// Call all callbacks asynchronously
	for _, f := range c.callbacks {
		go f(name, instance, args)
	}
}
