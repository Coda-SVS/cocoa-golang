package config

import "github.com/sasha-s/go-deadlock"

func init() {
	mtx = &deadlock.Mutex{}

	changedEventCallbackDict = make(map[string]map[string]*ChangedEventHandler)
	changedEventCallbackNameDict = make(map[string]string)
}

var (
	mtx *deadlock.Mutex

	changedEventCallbackDict     map[string]map[string]*ChangedEventHandler
	changedEventCallbackNameDict map[string]string
)

func AddChangedEventHandler(handler *ChangedEventHandler) bool {
	mtx.Lock()
	defer mtx.Unlock()

	h_invokePath := handler.InvokePath()

	nameHandlerDict, ok := changedEventCallbackDict[h_invokePath]
	if !ok {
		nameHandlerDict = make(map[string]*ChangedEventHandler)
		changedEventCallbackDict[h_invokePath] = nameHandlerDict
	}

	h_name := handler.Name()
	_, ok = nameHandlerDict[h_name]
	if ok {
		return false // 이미 등록된 EventHandler 존재
	} else {
		nameHandlerDict[h_name] = handler
		changedEventCallbackNameDict[h_name] = h_invokePath
		return true
	}
}

func DeleteChangedEventHandler(name string) bool {
	mtx.Lock()
	defer mtx.Unlock()

	h_invokePath, ok := changedEventCallbackNameDict[name]
	if !ok {
		return false
	}

	nameHandlerDict := changedEventCallbackDict[h_invokePath] // 무조건 존재해야 하므로 check pass

	delete(nameHandlerDict, name)
	delete(changedEventCallbackNameDict, name)
	return true
}

func invokeChangedEvent(path, key string, value any) {
	mtx.Lock()

	nameHandlerDict, ok := changedEventCallbackDict[path]
	if !ok {
		mtx.Unlock()
		return
	}

	tmp := make([]ChangedEventCallback, 0, len(nameHandlerDict))
	for _, handler := range nameHandlerDict {
		tmp = append(tmp, handler.Callback())
	}
	mtx.Unlock()

	for _, cb := range tmp {
		cb(key, value)
	}
}
