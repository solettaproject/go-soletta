package soletta

import "sync"

var (
	pointerDb         = map[uintptr]interface{}{}
	key       uintptr = 0xdead
	dbMutex   sync.Mutex
)

func mapPointer(goPointer interface{}) uintptr {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	pointerDb[key] = goPointer
	current := key
	key++
	return current
}

func getPointerMapping(pointer uintptr) interface{} {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	if v, ok := pointerDb[pointer]; ok {
		return v
	}
	panic("Pointer not found in pointer database")
}

func removePointerMapping(pointer uintptr) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	delete(pointerDb, pointer)
}
