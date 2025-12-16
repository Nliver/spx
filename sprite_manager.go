/*
 * Copyright (c) 2021 The XGo Authors (xgo.dev). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package spx

import (
	"log"

	"github.com/goplus/spx/v2/internal/engine"
	gtime "github.com/goplus/spx/v2/internal/time"
)

// shape is your existing Shape definition.
// type shape any

// spriteManager manages the lifecycle of all sprites/shapes.
// It is responsible for:
//   - activation (delayed add)
//   - destruction (delayed remove)
//   - render layer grouping
//   - minimizing per-frame allocations
type spriteManager struct {
	// active shapes
	items []Shape

	// shapes waiting to be activated
	tempItems []Shape

	// shapes waiting to be destroyed
	destroyItems []Shape
}

// newSpriteManager creates a spriteManager with preallocated buffers.
func newSpriteManager() *spriteManager {
	return &spriteManager{
		items:        make([]Shape, 0, 64),
		tempItems:    make([]Shape, 0, 16),
		destroyItems: make([]Shape, 0, 16),
	}
}

// reset clears all internal state while keeping allocated memory.
// It is safe to call between scenes or rounds.
func (sm *spriteManager) reset() {
	// clear active items
	sm.items = sm.items[:0]

	// clear pending activation and destruction
	sm.tempItems = sm.tempItems[:0]
	sm.destroyItems = sm.destroyItems[:0]
}

// clear releases all internal references.
// After calling clear, the manager should be discarded and replaced
// by a newly created spriteManager.
func (sm *spriteManager) clear() {
	sm.items = nil
	sm.tempItems = nil
	sm.destroyItems = nil
}

//
// ========== lifecycle operations ==========
//

// add immediately adds a shape to the active list.
// Equivalent to the original addShape.
func (sm *spriteManager) add(s Shape) {
	sm.items = append(sm.items, s)
}

// remove schedules a shape for destruction at the end of the frame.
// Equivalent to the original removeShape.
func (sm *spriteManager) remove(s Shape) {
	sm.destroyItems = append(sm.destroyItems, s)
}

// todo(refactor): old logic
func (sm *spriteManager) addShape(child Shape) {
	sm.add(child)
}

func (sm *spriteManager) addClonedShape(src, clone Shape) {
	items := sm.items
	idx := sm.doFindSprite(src)
	if idx < 0 {
		log.Println("addClonedShape: clone a deleted sprite")
		gco.Abort()
	}
	n := len(items)
	newItems := make([]Shape, n+1)
	copy(newItems[:idx], items)
	copy(newItems[idx+2:], items[idx+1:])
	newItems[idx] = clone
	newItems[idx+1] = src
	sm.items = newItems
	sm.updateRenderLayers()
}

func (sm *spriteManager) removeShape(child Shape) {
	items := sm.items
	for i, item := range items {
		if item == child {
			newItems := make([]Shape, len(items)-1)
			copy(newItems, items[:i])
			copy(newItems[i:], items[i+1:])
			sm.remove(item)
			sm.items = newItems
			return
		}
	}
	sm.updateRenderLayers()
}

func (sm *spriteManager) activateShape(child Shape) {
	items := sm.items
	for i, item := range items {
		if item == child {
			if i == 0 {
				return
			}
			newItems := make([]Shape, len(items))
			copy(newItems, items[:i])
			copy(newItems[i:], items[i+1:])
			newItems[len(items)-1] = child
			sm.items = newItems
			return
		}
	}
	sm.updateRenderLayers()
}

func (sm *spriteManager) goBackLayers(spr *SpriteImpl, n int) {
	if engine.HasLayerSortMethod() {
		log.Println("Cannot manually set sprite layer when a layer sort mode is active.")
		return
	}

	idx := sm.doFindSprite(spr)
	if idx < 0 {
		return
	}
	items := sm.items
	// go back
	if n > 0 {
		newIdx := idx
		for newIdx > 0 {
			newIdx--
			item := items[newIdx]
			if _, ok := item.(*SpriteImpl); ok {
				n--
				if n == 0 {
					break
				}
			}
		}
		// should consider that backdrop is always at the bottom
		if newIdx != idx {
			// p.getItems() requires immutable items, so we need copy before modify
			newItems := make([]Shape, len(items))
			copy(newItems, items[:newIdx])
			copy(newItems[newIdx+1:], items[newIdx:idx])
			copy(newItems[idx+1:], items[idx+1:])
			newItems[newIdx] = spr
			sm.items = newItems
		}
	} else if n < 0 { // go front
		newIdx := idx
		lastIdx := len(items) - 1
		if newIdx < lastIdx {
			for {
				newIdx++
				if newIdx >= lastIdx {
					break
				}
				item := items[newIdx]
				if _, ok := item.(*SpriteImpl); ok {
					n++
					if n == 0 {
						break
					}
				}
			}
		}
		if newIdx != idx {
			// p.getItems() requires immutable items, so we need copy before modify
			newItems := make([]Shape, len(items))
			copy(newItems, items[:idx])
			copy(newItems[idx:newIdx], items[idx+1:])
			copy(newItems[newIdx+1:], items[newIdx+1:])
			newItems[newIdx] = spr
			sm.items = newItems
		}
	}
	sm.updateRenderLayers()
}

//
// ========== flush phase ==========
//

// flushActivate moves all pending shapes into the active list.
func (sm *spriteManager) flushActivate() {
	if len(sm.tempItems) == 0 {
		return
	}

	for _, item := range sm.tempItems {
		if result, ok := item.(interface{ onUpdate(float64) }); ok {
			result.onUpdate(gtime.DeltaTime())
		}
	}
}

// flushDestroy removes all scheduled shapes from the active list.
// Removal is deferred to preserve stable iteration during updates.
func (sm *spriteManager) flushDestroy() {
	if len(sm.destroyItems) == 0 {
		return
	}

	for _, item := range sm.destroyItems {
		sprite, ok := item.(*SpriteImpl)
		if ok && sprite.syncSprite != nil {
			sprite.syncSprite.Destroy()
			sprite.syncSprite = nil
		}
	}

	sm.destroyItems = sm.destroyItems[:0]
}

//
// ========== render layer management ==========
//

// updateRenderLayers set only when marked dirty.
func (sm *spriteManager) updateRenderLayers() {
	if engine.HasLayerSortMethod() {
		return
	}

	layer := 0
	for _, item := range sm.items {
		if sp, ok := item.(*SpriteImpl); ok {
			layer++
			sp.setLayer(layer)
		}
	}
}

//
// ========== query helpers ==========
//

// all returns all active shapes.
// all requires immutable items, so need copy before modify
func (sm *spriteManager) all() []Shape {
	return sm.items
}

func (sm *spriteManager) getTempShapes() []Shape {
	sm.tempItems = getTempShapes(sm.tempItems, sm.items)
	return sm.tempItems
}

// count returns the number of active shapes.
func (sm *spriteManager) count() int {
	return len(sm.items)
}

//
// ========== internal helpers ==========
//

func (sm *spriteManager) doFindSprite(src Shape) int {
	for idx, item := range sm.items {
		if item == src {
			return idx
		}
	}
	return -1
}

func (sm *spriteManager) findSprite(name SpriteName) *SpriteImpl {
	for _, item := range sm.items {
		if sp, ok := item.(*SpriteImpl); ok {
			if !sp.isCloned_ && sp.name == name {
				return sp
			}
		}
	}
	return nil
}

func getTempShapes(dst []Shape, src []Shape) []Shape {
	if dst == nil {
		dst = make([]Shape, 50)
	}
	dst = dst[:0]
	if cap(dst) < len(src) {
		dst = make([]Shape, len(src))
	} else {
		dst = dst[:len(src)]
	}
	copy(dst, src)
	return dst
}
