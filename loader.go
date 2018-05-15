package tempura

import (
	"bytes"
	"image"
	"io"

	"io/ioutil"

	"github.com/explodes/tempura/tinge"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/audio/wav"
	"github.com/pkg/errors"
	"golang.org/x/image/font"
)

type AssetFunc func(name string) ([]byte, error)

type Loader interface {
	Reader(name string) (io.Reader, error)
	ReadCloser(name string) (audio.ReadSeekCloser, error)
	Image(name string, transforms ...tinge.Transform) (image.Image, error)
	EbitenImage(name string, filter ebiten.Filter, transforms ...tinge.Transform) (*ebiten.Image, error)
	Font(name string) (*truetype.Font, error)
	Face(name string, size float64) (font.Face, error)
	SFX(context *audio.Context, fmt, name string) (AudioPlayer, error)
	AudioLoop(context *audio.Context, fmt, name string) (*audio.Player, error)
}

var _ Loader = (*loaderImpl)(nil)

type loaderImpl struct {
	bytes AssetFunc
	debug bool
}

func NewLoader(assets AssetFunc) Loader {
	return NewLoaderDebug(assets, false)
}

func NewLoaderDebug(assets AssetFunc, debug bool) Loader {
	return &loaderImpl{
		bytes: assets,
		debug: debug,
	}
}

func (l *loaderImpl) Reader(name string) (io.Reader, error) {
	b, err := l.bytes(name)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

type readSeekCloserImpl struct {
	*bytes.Reader
}

func (readSeekCloserImpl) Close() error {
	return nil
}

func (l *loaderImpl) ReadCloser(name string) (audio.ReadSeekCloser, error) {
	b, err := l.bytes(name)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(b)
	return &readSeekCloserImpl{Reader: reader}, nil
}

func (l *loaderImpl) getImage(name string, transforms ...tinge.Transform) (image.Image, error) {
	r, err := l.Reader(name)
	if err != nil {
		return nil, err
	}
	src, _, err := image.Decode(r)
	for _, transform := range transforms {
		src, err = transform(src)
		if err != nil {
			return nil, err
		}
	}
	return src, err
}

func (l *loaderImpl) Image(name string, transforms ...tinge.Transform) (image.Image, error) {
	if l.debug {
		defer LogDuration("Image for %s", name).End()
	}
	return l.getImage(name, transforms...)
}

func (l *loaderImpl) EbitenImage(name string, filter ebiten.Filter, transforms ...tinge.Transform) (*ebiten.Image, error) {
	if l.debug {
		defer LogDuration("EbitenImage for %s", name).End()
	}
	src, err := l.getImage(name, transforms...)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(src, filter)
}

func (l *loaderImpl) font(name string) (*truetype.Font, error) {
	b, err := l.bytes(name)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(b)
}

func (l *loaderImpl) Font(name string) (*truetype.Font, error) {
	if l.debug {
		defer LogDuration("Font for %s", name).End()
	}
	return l.font(name)
}

func (l *loaderImpl) Face(name string, size float64) (font.Face, error) {
	if l.debug {
		defer LogDuration("Face for %s", name).End()
	}
	f, err := l.font(name)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})
	return face, nil
}

type readSeekCloseLengther interface {
	audio.ReadSeekCloser
	Length() int64
}

func (l *loaderImpl) audioStream(context *audio.Context, fmt, name string) (readSeekCloseLengther, error) {
	rsc, err := l.ReadCloser(name)
	if err != nil {
		return nil, err
	}
	switch fmt {
	case "mp3":
		return mp3.Decode(context, rsc)
	case "wav":
		return wav.Decode(context, rsc)
	default:
		return nil, errors.Errorf("format not supported: %s", fmt)
	}
}

type AudioPlayer interface {
	Play()
	Pause()
	SetVolume(v float64)
}

type audioReplayer struct {
	stream     []byte
	context    *audio.Context
	lastPlayer *audio.Player
	volume     float64
}

func newPlayer(context *audio.Context, stream []byte) AudioPlayer {
	return &audioReplayer{
		stream:  stream,
		context: context,
		volume:  1,
	}
}

func (a *audioReplayer) Play() {
	p, err := audio.NewPlayerFromBytes(a.context, a.stream)
	if err != nil {
		return
	}
	p.SetVolume(a.volume)
	p.Play()
	a.lastPlayer = p
}

func (a *audioReplayer) Pause() {
	if a.lastPlayer != nil {
		a.lastPlayer.Pause()
	}
}

func (a *audioReplayer) SetVolume(v float64) {
	if v > 1 {
		v = 1
	} else if v < 0 {
		v = 0
	}
	a.volume = v
}

func (l *loaderImpl) SFX(context *audio.Context, fmt, name string) (AudioPlayer, error) {
	if l.debug {
		defer LogDuration("SFX for %s", name).End()
	}
	stream, err := l.audioStream(context, fmt, name)
	if err != nil {
		return nil, err
	}
	streamBytes, err := ioutil.ReadAll(stream)
	if err != nil {
		return nil, err
	}
	return newPlayer(context, streamBytes), nil
}

func (l *loaderImpl) AudioLoop(context *audio.Context, fmt, name string) (*audio.Player, error) {
	if l.debug {
		defer LogDuration("AudioLoop for %s", name).End()
	}
	stream, err := l.audioStream(context, fmt, name)
	if err != nil {
		return nil, err
	}
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	return audio.NewPlayer(context, loop)
}
