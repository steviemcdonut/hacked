package compression

// EncodedFrame contains the streams of one compressed frame.
type EncodedFrame struct {
	Bitstream  []byte
	Maskstream []byte
}

type tileDelta [PixelPerTile]byte

type frameDelta struct {
	tiles []tileDelta
}

// SceneEncoder encodes an entire scene of bitmaps sharing the same palette.
type SceneEncoder struct {
	hTiles int
	vTiles int
	stride int

	lastFrame []byte
	deltas    []frameDelta
}

// NewSceneEncoder returns a new instance.
func NewSceneEncoder(width, height int) *SceneEncoder {
	e := &SceneEncoder{
		hTiles: width / TileSideLength,
		vTiles: height / TileSideLength,
		stride: width,
	}
	e.lastFrame = make([]byte, e.vTiles*TileSideLength*e.stride)
	return e
}

// AddFrame registers a further frame to the scene.
func (e *SceneEncoder) AddFrame(frame []byte) error {
	var delta frameDelta
	for vTile := 0; vTile < e.vTiles; vTile++ {
		vStart := vTile * e.stride * TileSideLength
		for hTile := 0; hTile < e.hTiles; hTile++ {
			delta.tiles = append(delta.tiles, e.deltaTile(vStart+(hTile*TileSideLength), frame))
		}
	}
	e.deltas = append(e.deltas, delta)

	e.copyLastFrame(frame)
	return nil
}

func (e *SceneEncoder) deltaTile(offset int, frame []byte) tileDelta {
	var delta tileDelta
	for y := 0; y < TileSideLength; y++ {
		start := offset + (y * e.stride)
		for x := 0; x < TileSideLength; x++ {
			pixel := frame[start+x]
			if pixel != e.lastFrame[start+x] {
				delta[y*TileSideLength+x] = pixel
			}
		}
	}
	return delta
}

func (e *SceneEncoder) copyLastFrame(frame []byte) {
	hPixel := e.hTiles * TileSideLength
	for y := 0; y < e.vTiles*TileSideLength; y++ {
		start := y * e.stride
		copy(e.lastFrame[start:start+hPixel], frame[start:])
	}
}

// Encode processes all the previously registered frames and creates the necessary components for decoding.
func (e *SceneEncoder) Encode() (words []ControlWord, paletteLookup []byte, frames []EncodedFrame, err error) {
	var paletteLookupWriter PaletteLookupWriter
	frames = make([]EncodedFrame, len(e.deltas))
	for frameIndex := 0; frameIndex < len(e.deltas); frameIndex++ {
		outFrame := &frames[frameIndex]
		delta := e.deltas[frameIndex]
		var bitstream BitstreamWriter

		controlIndex := len(words)
		paletteIndex := paletteLookupWriter.Write(delta.tiles[0][:])

		words = append(words, ControlWordOf(12, CtrlColorTile16ColorsMasked, paletteIndex))
		bitstream.Write(12, uint32(controlIndex))

		outFrame.Bitstream = bitstream.Buffer()
		outFrame.Maskstream = []byte{0x10, 0x32, 0x54, 0x76, 0x98, 0xBA, 0xDC, 0xFE}
	}
	paletteLookup = paletteLookupWriter.Buffer

	return
}
