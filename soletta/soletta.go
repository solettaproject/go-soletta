//Package soletta provides Go bindings for Soletta library.
package soletta

/*
#cgo pkg-config: soletta

#include "soletta.h"
*/
import "C"
import "errors"

//Soletta API version
const MainloopSourceAPIVersion uint16 = C.SOL_MAINLOOP_SOURCE_TYPE_API_VERSION

//Initializes the Soletta library.
//
//This function setup all needed infrastructure.
//It should be called prior the use of any Soletta API.
func Init() error {
	if C.sol_init() < 0 {
		return errors.New("Failed to initialize Soletta")
	}
	return nil
}

//Runs the main loop.
//
//This function executes the main loop and it will return only after Quit() is called
func Run() error {
	if C.sol_run() != C.EXIT_SUCCESS {
		return errors.New("Soletta mainloop exited with non success value")
	}
	return nil
}

//Terminates the main loop.
func Quit() {
	C.sol_quit()
}

//Shutdown Soletta library.
//
//This function shuts down Soletta and once it's called, no other Soletta API should be used.
func Shutdown() {
	C.sol_shutdown()
}
