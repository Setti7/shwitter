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

// Map of the sender name (entity struct name) and its signalConfig, which contains the callbacks
var signalMapping = map[string]signalConfig{}

type signalCallback func(instance interface{}, args ...interface{})

// Connect a callback to a signal sender.
func (s signal) Connect(sender interface{}, f signalCallback) {
	name := reflect.TypeOf(sender).String()

	// Set the signal sender and append the callback
	callbacks := append(signalMapping[name].callbacks, f)
	signalMapping[name] = signalConfig{sender: sender, callbacks: callbacks}
}

// Emit a signal, calling all of its registered callbacks.
func (s signal) Emit(sender interface{}, instance interface{}, args ...interface{}) {
	name := reflect.TypeOf(sender).String()
	c := signalMapping[name]

	// Call all callbacks asynchronously
	for _, f := range c.callbacks {
		go f(instance, args)
	}
}

// Clear the signal config of the given sender.
//
// Useful for testing: disconnect all signals for easier unit testing.
func Clear(s signal, sender interface{}) {
	name := reflect.TypeOf(sender).String()
	signalMapping[name] = signalConfig{}
}
