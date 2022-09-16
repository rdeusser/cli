package cli

func convertInt64ToInt(i []int64) []int {
	var values []int
	for _, v := range i {
		values = append(values, int(v))
	}
	return values
}

func convertInt64ToInt8(i []int64) []int8 {
	var values []int8
	for _, v := range i {
		values = append(values, int8(v))
	}
	return values
}

func convertInt64ToInt16(i []int64) []int16 {
	var values []int16
	for _, v := range i {
		values = append(values, int16(v))
	}
	return values
}

func convertInt64ToInt32(i []int64) []int32 {
	var values []int32
	for _, v := range i {
		values = append(values, int32(v))
	}
	return values
}

func convertUint64ToUint(i []uint64) []uint {
	var values []uint
	for _, v := range i {
		values = append(values, uint(v))
	}
	return values
}

func convertUint64ToUint8(i []uint64) []uint8 {
	var values []uint8
	for _, v := range i {
		values = append(values, uint8(v))
	}
	return values
}

func convertUint64ToUint16(i []uint64) []uint16 {
	var values []uint16
	for _, v := range i {
		values = append(values, uint16(v))
	}
	return values
}

func convertUint64ToUint32(i []uint64) []uint32 {
	var values []uint32
	for _, v := range i {
		values = append(values, uint32(v))
	}
	return values
}

func convertUint64ToUintptr(i []uint64) []uintptr {
	var values []uintptr
	for _, v := range i {
		values = append(values, uintptr(v))
	}
	return values
}

func convertFloat64ToFloat32(f []float64) []float32 {
	var values []float32
	for _, v := range f {
		values = append(values, float32(v))
	}
	return values
}

func convertComplex128ToComplex64(c []complex128) []complex64 {
	var values []complex64
	for _, v := range c {
		values = append(values, complex64(v))
	}
	return values
}
