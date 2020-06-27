package led

type EffectMode byte

// TODO Implement other modes
const (
	Raw EffectMode = iota
)

var Header = []byte{0x4C, 0x58}

func MakeLuxPayload(mode EffectMode, effectPayload []byte) []byte {
	result := make([]byte, 0)

	result = append(result, Header...)
	result = append(result, byte(mode))
	result = append(result, effectPayload...)

	return result
}

func MakeRawModeLuxPayload(ledCount uint8, grbData []byte) []byte {
	effectPayload := make([]byte, 0)

	effectPayload = append(effectPayload, ledCount)
	effectPayload = append(effectPayload, grbData...)

	return MakeLuxPayload(Raw, effectPayload)
}
