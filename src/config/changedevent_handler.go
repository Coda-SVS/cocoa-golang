package config

type ChangedEventCallback func(key string, value any)

type ChangedEventHandler struct {
	name       string
	invokePath string
	callback   ChangedEventCallback
}

func NewChangedEventHandler(name string, invokePath string, callback ChangedEventCallback) *ChangedEventHandler {
	return &ChangedEventHandler{
		name:       name,
		invokePath: invokePath,
		callback:   callback,
	}
}

func (ceh *ChangedEventHandler) Name() string {
	return ceh.name
}

func (ceh *ChangedEventHandler) InvokePath() string {
	return ceh.invokePath
}

func (ceh *ChangedEventHandler) Callback() ChangedEventCallback {
	return ceh.callback
}
