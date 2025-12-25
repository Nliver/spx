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
	"github.com/goplus/spbase/mathf"
)

// ======================== Pen Component ========================
// This file contains pen-related functionality for sprites,
// including pen drawing, color control, and size management.

// -----------------------------------------------------------------------------
// Pen Color Parameter Types
// -----------------------------------------------------------------------------

type PenColorParam int

const (
	PenHue PenColorParam = iota
	PenSaturation
	PenBrightness
	PenTransparency
)

// -----------------------------------------------------------------------------
// Pen Control
// -----------------------------------------------------------------------------

func (p *SpriteImpl) PenUp() {
	p.checkOrCreatePen()
	p.isPenDown = false
	penMgr.PenUp(*p.penObj)
}

func (p *SpriteImpl) PenDown() {
	p.checkOrCreatePen()
	p.isPenDown = true
	p.movePen(p.x, p.y)
	penMgr.PenDown(*p.penObj, false)
}

func (p *SpriteImpl) Stamp() {
	p.checkOrCreatePen()
	penMgr.SetPenStampTexture(*p.penObj, p.getCostumePath())
	penMgr.PenStamp(*p.penObj)
}

// -----------------------------------------------------------------------------
// Pen Color Control
// -----------------------------------------------------------------------------

func (p *SpriteImpl) SetPenColor__0(color Color) {
	p.checkOrCreatePen()
	p.penColor = toMathfColor(color)
	p.applyPenColorProperty()
}

func (p *SpriteImpl) SetPenColor__1(kind PenColorParam, value float64) {
	switch kind {
	case PenHue:
		p.setPenHue(value)
	case PenSaturation:
		p.setPenSaturation(value)
	case PenBrightness:
		p.setPenBrightness(value)
	case PenTransparency:
		p.setPenTransparency(value)
	}
}

func (p *SpriteImpl) ChangePenColor(kind PenColorParam, delta float64) {
	switch kind {
	case PenHue:
		p.changePenHue(delta)
	case PenSaturation:
		p.changePenSaturation(delta)
	case PenBrightness:
		p.changePenBrightness(delta)
	case PenTransparency:
		p.changePenTransparency(delta)
	}
}

// -----------------------------------------------------------------------------
// Pen HSV Color Components
// -----------------------------------------------------------------------------

func (p *SpriteImpl) setPenHue(value float64) {
	p.checkOrCreatePen()
	p.penHue = mathf.Clamp(value, 0, 100)
	p.applyPenHsvProperty()
}

func (p *SpriteImpl) changePenHue(delta float64) {
	p.setPenHue(p.penHue + delta)
}

func (p *SpriteImpl) setPenSaturation(value float64) {
	p.checkOrCreatePen()
	p.penSaturation = mathf.Clamp(value, 0, 100)
	p.applyPenHsvProperty()
}

func (p *SpriteImpl) changePenSaturation(delta float64) {
	p.setPenSaturation(p.penSaturation + delta)
}

func (p *SpriteImpl) setPenBrightness(value float64) {
	p.checkOrCreatePen()
	p.penBrightness = mathf.Clamp(value, 0, 100)
	p.applyPenHsvProperty()
}

func (p *SpriteImpl) changePenBrightness(delta float64) {
	p.setPenBrightness(p.penBrightness + delta)
}

func (p *SpriteImpl) setPenTransparency(value float64) {
	p.checkOrCreatePen()
	p.penTransparency = mathf.Clamp(value, 0, 100)
	p.applyPenHsvProperty()
}

func (p *SpriteImpl) changePenTransparency(delta float64) {
	p.setPenTransparency(p.penTransparency + delta)
}

// -----------------------------------------------------------------------------
// Pen Size Control
// -----------------------------------------------------------------------------

func (p *SpriteImpl) SetPenSize(size float64) {
	p.checkOrCreatePen()
	p.penWidth = size
	penMgr.SetPenSizeTo(*p.penObj, size)
}

func (p *SpriteImpl) ChangePenSize(delta float64) {
	p.checkOrCreatePen()
	p.SetPenSize(p.penWidth + delta)
}

// -----------------------------------------------------------------------------
// Internal Pen Management
// -----------------------------------------------------------------------------

func (p *SpriteImpl) checkOrCreatePen() {
	if p.penObj == nil {
		obj := penMgr.CreatePen()
		p.penObj = &obj
		p.penTransparency = p.penColor.A * 100
	}
}

func (p *SpriteImpl) destroyPen() {
	if p.penObj != nil {
		penMgr.DestroyPen(*p.penObj)
		p.penObj = nil
	}
}

func (p *SpriteImpl) movePen(x, y float64) {
	if p.penObj == nil {
		return
	}
	applyRenderOffset(p, &x, &y)
	penMgr.MovePenTo(*p.penObj, mathf.NewVec2(x, -y))
}

func (p *SpriteImpl) applyPenColorProperty() {
	p.checkOrCreatePen()
	h, s, v := p.penColor.ToHSV()
	p.penHue = (h / 360) * 100
	p.penSaturation = s * 100
	p.penBrightness = v * 100
	p.penTransparency = p.penColor.A * 100
	penMgr.SetPenColorTo(*p.penObj, p.penColor)
}

func (p *SpriteImpl) applyPenHsvProperty() {
	color := mathf.NewColorHSV((p.penHue/100)*360, p.penSaturation/100, p.penBrightness/100)
	p.penColor = color
	p.penColor.A = p.penTransparency / 100
	penMgr.SetPenColorTo(*p.penObj, p.penColor)
}
