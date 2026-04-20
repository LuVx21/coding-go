package runtime_x

import (
	"log"
	"runtime"
)

func PrintCallerChain() {
	const depth = 32
	pcs := make([]uintptr, depth)

	frames := runtime.CallersFrames(pcs[:runtime.Callers(2, pcs)])

	log.Println("📞 调用链：")
	for {
		frame, more := frames.Next()
		log.Printf("  -> %s\n", frame.Function)
		if !more {
			break
		}
	}
}
