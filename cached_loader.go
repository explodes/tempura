package tempura

import (
	"image"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"golang.org/x/image/font"
)

var _ Loader = (*cachedLoader)(nil)

type ebitenImage struct {
	name   string
	filter ebiten.Filter
}

type face struct {
	name string
	size float64
}

type sfx struct {
	context   *audio.Context
	fmt, name string
}

type cachedLoader struct {
	Loader
	imageCache       map[string]image.Image
	ebitenImageCache map[ebitenImage]*ebiten.Image
	fontCache        map[string]*truetype.Font
	faceCache        map[face]font.Face
	sfxCache         map[sfx]AudioPlayer
	loopCache        map[sfx]*audio.Player
}

func NewCachedLoader(loader Loader) Loader {
	return &cachedLoader{
		Loader: loader,
	}
}

func (l *cachedLoader) Image(name string) (image.Image, error) {
	if l.imageCache != nil {
		v, ok := l.imageCache[name]
		if ok {
			return v, nil
		}
	}
	v, err := l.Loader.Image(name)
	if err != nil {
		return nil, err
	}
	if l.imageCache == nil {
		l.imageCache = make(map[string]image.Image)
	}
	l.imageCache[name] = v
	return v, nil
}

func (l *cachedLoader) EbitenImage(name string, filter ebiten.Filter) (*ebiten.Image, error) {
	key := ebitenImage{name: name, filter: filter}
	if l.ebitenImageCache != nil {
		v, ok := l.ebitenImageCache[key]
		if ok {
			return v, nil
		}
	}
	v, err := l.Loader.EbitenImage(name, filter)
	if err != nil {
		return nil, err
	}
	if l.ebitenImageCache == nil {
		l.ebitenImageCache = make(map[ebitenImage]*ebiten.Image)
	}
	l.ebitenImageCache[key] = v
	return v, nil
}

func (l *cachedLoader) Font(name string) (*truetype.Font, error) {
	if l.fontCache != nil {
		v, ok := l.fontCache[name]
		if ok {
			return v, nil
		}
	}
	v, err := l.Loader.Font(name)
	if err != nil {
		return nil, err
	}
	if l.fontCache == nil {
		l.fontCache = make(map[string]*truetype.Font)
	}
	l.fontCache[name] = v
	return v, nil
}

func (l *cachedLoader) Face(name string, size float64) (font.Face, error) {
	key := face{name: name, size: size}
	if l.faceCache != nil {
		v, ok := l.faceCache[key]
		if ok {
			return v, nil
		}
	}
	v, err := l.Loader.Face(name, size)
	if err != nil {
		return nil, err
	}
	if l.faceCache == nil {
		l.faceCache = make(map[face]font.Face)
	}
	l.faceCache[key] = v
	return v, nil
}

func (l *cachedLoader) SFX(context *audio.Context, fmt, name string) (AudioPlayer, error) {
	key := sfx{context: context, name: name, fmt: fmt}
	if l.sfxCache != nil {
		v, ok := l.sfxCache[key]
		if ok {
			return v, nil
		}
	}
	v, err := l.Loader.SFX(context, fmt, name)
	if err != nil {
		return nil, err
	}
	if l.sfxCache == nil {
		l.sfxCache = make(map[sfx]AudioPlayer)
	}
	l.sfxCache[key] = v
	return v, nil
}

func (l *cachedLoader) AudioLoop(context *audio.Context, fmt, name string) (*audio.Player, error) {
	key := sfx{context: context, name: name, fmt: fmt}
	if l.loopCache != nil {
		v, ok := l.loopCache[key]
		if ok {
			return v, nil
		}
	}
	v, err := l.Loader.AudioLoop(context, fmt, name)
	if err != nil {
		return nil, err
	}
	if l.loopCache == nil {
		l.loopCache = make(map[sfx]*audio.Player)
	}
	l.loopCache[key] = v
	return v, nil
}
