package plugin

import (
	"errors"
	"fmt"
	"log"
	"syscall/js"

	"github.com/goplus/builder/tools/ai"
	"github.com/goplus/builder/tools/ai/wasmtrans"
	"github.com/goplus/ixgo"
)

type AIPlugin struct {
	aiDescription                 string
	aiInteractionAPIEndpoint      string
	aiInteractionAPITokenProvider func() string
}

func init() {
	Register("ai", &AIPlugin{})
}

func (p *AIPlugin) RegisterJS() {
	js.Global().Set("setAIDescription", js.FuncOf(p.setDescription))
	js.Global().Set("setAIInteractionAPIEndpoint", js.FuncOf(p.setEndpoint))
	js.Global().Set("setAIInteractionAPITokenProvider", js.FuncOf(p.setTokenProvider))
}

// RegisterPatch registers any required ixgo patch for AI plugin.
func (p *AIPlugin) RegisterPatch(ctx *ixgo.Context) error {
	patch := `
package ai

import . "github.com/goplus/builder/tools/ai"

// Generic helper for Player command handling
func Gopt_Player_Gopx_OnCmd[T any](p *Player, handler func(cmd T) error) {
	var cmd T
	PlayerOnCmd_(p, cmd, handler)
}
`
	return ctx.RegisterPatch("github.com/goplus/builder/tools/ai", patch)
}

func (p *AIPlugin) Init() {
	ai.SetDefaultTransport(wasmtrans.New(
		wasmtrans.WithEndpoint(p.aiInteractionAPIEndpoint),
		wasmtrans.WithTokenProvider(p.aiInteractionAPITokenProvider),
	))
	ai.SetDefaultKnowledgeBase(map[string]any{
		"AI-generated descriptive summary of the game world": p.aiDescription,
	})
}

// --- JS Function Handlers ---

func (p *AIPlugin) setDescription(this js.Value, args []js.Value) any {
	if len(args) > 0 {
		p.aiDescription = args[0].String()
	}
	return nil
}

func (p *AIPlugin) setEndpoint(this js.Value, args []js.Value) any {
	if len(args) > 0 {
		p.aiInteractionAPIEndpoint = args[0].String()
	}
	return nil
}

func (p *AIPlugin) setTokenProvider(this js.Value, args []js.Value) any {
	if len(args) > 0 && args[0].Type() == js.TypeFunction {
		tokenProviderFunc := args[0]
		p.aiInteractionAPITokenProvider = func() string {
			result := tokenProviderFunc.Invoke()
			if result.Type() != js.TypeObject || result.Get("then").IsUndefined() {
				return result.String()
			}

			resultChan := make(chan string, 1)
			then := js.FuncOf(func(this js.Value, args []js.Value) any {
				var result string
				if len(args) > 0 {
					result = args[0].String()
				}
				resultChan <- result
				return nil
			})
			defer then.Release()

			errChan := make(chan error, 1)
			catch := js.FuncOf(func(this js.Value, args []js.Value) any {
				errMsg := "promise rejected"
				if len(args) > 0 {
					errVal := args[0]
					if errVal.Type() == js.TypeObject && errVal.Get("message").Type() == js.TypeString {
						errMsg = fmt.Sprintf("promise rejected: %s", errVal.Get("message"))
					} else if errVal.Type() == js.TypeString {
						errMsg = fmt.Sprintf("promise rejected: %s", errVal)
					} else {
						errMsg = fmt.Sprintf("promise rejected: %v", errVal)
					}
				}
				errChan <- errors.New(errMsg)
				return nil
			})
			defer catch.Release()

			result.Call("then", then).Call("catch", catch)
			select {
			case result := <-resultChan:
				return result
			case err := <-errChan:
				log.Printf("failed to get token: %v", err)
				return ""
			}
		}
	}
	return nil
}
