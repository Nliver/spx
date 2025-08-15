package enginewrap

var mainCallback func(call func())

func Init(call func(f func())) {
	mainCallback = call
}

func callInMainThread(call func()) {
	mainCallback(call)
}

const (
	MOUSE_BUTTON_LEFT   int64 = 1
	MOUSE_BUTTON_RIGHT  int64 = 2
	MOUSE_BUTTON_MIDDLE int64 = 3
)

// =============== input ===================
func (pself *inputMgrImpl) MousePressed() bool {
	return inputMgr.GetMouseState(MOUSE_BUTTON_LEFT) || inputMgr.GetMouseState(MOUSE_BUTTON_RIGHT)
}

// =============== window ===================

func (pself *platformMgrImpl) SetRunnableOnUnfocused(flag bool) {
	if !flag {
		println("TODO tanjp SetRunnableOnUnfocused")
	}
}
