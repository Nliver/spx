package gdspx

import (
	inengine "github.com/goplus/spx/v2/pkg/gdspx/internal/engine"
	ffi "github.com/goplus/spx/v2/pkg/gdspx/internal/ffi"
	engine "github.com/goplus/spx/v2/pkg/gdspx/pkg/engine"
)

func IsWebIntepreterMode() bool {
	return inengine.IsWebIntepreterMode()
}

func LinkEngine(callback engine.EngineCallbackInfo) {
	inengine.Link(engine.EngineCallbackInfo(callback))
}

func ToGdArray(slice interface{}) ffi.GdArray {
	return ffi.ToGdArray(slice)
}
