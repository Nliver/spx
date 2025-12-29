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

// ======================== Sound Component ========================
// This file contains sound-related functionality for sprites,
// including sound playback, volume control, and sound effects.

// -----------------------------------------------------------------------------
// Sound Effect Types
// -----------------------------------------------------------------------------

type SoundEffectKind int

const (
	SoundPanEffect SoundEffectKind = iota
	SoundPitchEffect
)

// -----------------------------------------------------------------------------
// Internal Audio Management
// -----------------------------------------------------------------------------

func (p *SpriteImpl) playAudio(name SoundName, loop bool) soundId {
	p.checkAudioId()
	return p.g.playSound(p.syncSprite, p.audioId, name, loop, p.g.audioAttenuation, p.g.audioMaxDistance)
}

func (p *SpriteImpl) checkAudioId() {
	if p.audioId == 0 {
		p.audioId = p.g.sounds.allocAudio()
	}
}

// -----------------------------------------------------------------------------
// Sound Playback Control
// -----------------------------------------------------------------------------

func (p *SpriteImpl) Play__0(name SoundName, loop bool) {
	p.checkAudioId()
	p.g.playSound(p.syncSprite, p.audioId, name, loop, p.g.audioAttenuation, p.g.audioMaxDistance)
}

func (p *SpriteImpl) Play__1(name SoundName) {
	p.Play__0(name, false)
}

func (p *SpriteImpl) PlayAndWait(name SoundName) {
	p.checkAudioId()
	p.g.playSoundAndWait(p.syncSprite, p.audioId, name, p.g.audioAttenuation, p.g.audioMaxDistance)
}

func (p *SpriteImpl) doSoundAction(name SoundName, action func(name SoundName)) {
	action(name)
}

func (p *SpriteImpl) PausePlaying(name SoundName) {
	p.doSoundAction(name, p.g.pauseSound)
}

func (p *SpriteImpl) ResumePlaying(name SoundName) {
	p.doSoundAction(name, p.g.resumeSound)
}

func (p *SpriteImpl) StopPlaying(name SoundName) {
	p.doSoundAction(name, p.g.stopSound)
}

// -----------------------------------------------------------------------------
// Volume Control
// -----------------------------------------------------------------------------

func (p *SpriteImpl) Volume() float64 {
	p.checkAudioId()
	return p.g.sounds.getVolume(p.audioId)
}

func (p *SpriteImpl) SetVolume(volume float64) {
	p.checkAudioId()
	p.g.sounds.setVolume(p.audioId, volume)
}

func (p *SpriteImpl) ChangeVolume(delta float64) {
	p.checkAudioId()
	p.g.sounds.changeVolume(p.audioId, delta)
}

// -----------------------------------------------------------------------------
// Sound Effects Control
// -----------------------------------------------------------------------------

func (p *SpriteImpl) GetSoundEffect(kind SoundEffectKind) float64 {
	p.checkAudioId()
	return p.g.sounds.getEffect(p.audioId, kind)
}

func (p *SpriteImpl) SetSoundEffect(kind SoundEffectKind, value float64) {
	p.checkAudioId()
	p.g.sounds.setEffect(p.audioId, kind, value)
}

func (p *SpriteImpl) ChangeSoundEffect(kind SoundEffectKind, delta float64) {
	p.checkAudioId()
	p.g.sounds.changeEffect(p.audioId, kind, delta)
}
