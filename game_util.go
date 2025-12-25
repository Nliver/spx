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
	"math/rand"

	"github.com/goplus/spbase/mathf"
)

// -----------------------------------------------------------------------------
// Window and World Size Utilities

func (p *Game) getWindowSize() mathf.Vec2 {
	x, y := p.windowSize_()
	return mathf.NewVec2(float64(x), float64(y))
}

func (p *Game) windowSize_() (int, int) {
	if p.windowWidth_ == 0 {
		p.doWindowSize()
	}
	return p.windowWidth_, p.windowHeight_
}

func (p *Game) doWindowSize() {
	if p.windowWidth_ == 0 {
		c := p.costumes[p.costumeIndex_]
		p.windowWidth_, p.windowHeight_ = c.getSize()
	}
}

func (p *Game) worldSize_() (int, int) {
	if p.worldWidth_ == 0 {
		p.doWorldSize()
	}
	return p.worldWidth_, p.worldHeight_
}

func (p *Game) doWorldSize() {
	if p.worldWidth_ == 0 {
		c := p.costumes[p.costumeIndex_]
		p.worldWidth_, p.worldHeight_ = c.getSize()
	}
}

// -----------------------------------------------------------------------------
// Touch and Collision Utilities

func (p *Game) touchingPoint(dst *SpriteImpl, x, y float64) bool {
	return dst.touchPoint(x, y)
}

func (p *Game) touchingSpriteBy(dst *SpriteImpl, name string) *SpriteImpl {
	if dst == nil {
		return nil
	}

	for _, item := range p.spriteMgr.items {
		if sp, ok := item.(*SpriteImpl); ok && sp != dst {
			if sp.name == name && (sp.isVisible && !sp.isDying) {
				if sp.touchingSprite(dst) {
					return sp
				}

			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// Object Position Utilities

func (p *Game) objectPos(obj any) (float64, float64) {
	switch v := obj.(type) {
	case SpriteName:
		if sp := p.spriteMgr.findSprite(v); sp != nil {
			return sp.getXY()
		}
		panic("objectPos: sprite not found - " + v)
	case specialObj:
		if v == Mouse {
			return p.getMousePos()
		}
	case Pos:
		if v == Random {
			worldW, worldH := p.worldSize_()
			mx, my := rand.Intn(worldW), rand.Intn(worldH)
			return float64(mx - (worldW >> 1)), float64((worldH >> 1) - my)
		}
	case Sprite:
		return spriteOf(v).getXY()
	}
	panic("objectPos: unexpected input")
}

// -----------------------------------------------------------------------------
// Pen Utilities

func (p *Game) EraseAll() {
	penMgr.DestroyAllPens()
}

// -----------------------------------------------------------------------------
// Shape Management Utilities

func (p *Game) getItems() []Shape {
	return p.spriteMgr.all()
}

func (p *Game) addShape(child Shape) {
	p.spriteMgr.addShape(child)
}

func (p *Game) addClonedShape(src, clone Shape) {
	p.spriteMgr.addClonedShape(src, clone)
}

func (p *Game) removeShape(child Shape) {
	p.spriteMgr.removeShape(child)
}

func (p *Game) activateShape(child Shape) {
	p.spriteMgr.activateShape(child)
}

func (p *Game) findSprite(name SpriteName) *SpriteImpl {
	return p.spriteMgr.findSprite(name)
}

func (p *Game) getAllShapes() []Shape {
	return p.spriteMgr.all()
}

func (p *Game) getTempShapes() []Shape {
	return p.spriteMgr.getTempShapes()
}
