//go:build profiler
// +build profiler

package profiler

func init() {
	Enabled = true
	println("Profiler Enabled (macro_on)")
}
