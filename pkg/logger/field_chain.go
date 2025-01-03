package logger

type FieldChain []Field

func WithField[T any](k string, v T) FieldChain {
	return FieldChain{NewField(k, v)}
}

func (c FieldChain) WithField(k string, v interface{}) FieldChain {
	return append(c, Field{k, v})
}

func (c FieldChain) Info(msg string, fields ...Field) {
	Info(msg, append(c, fields...)...)
}

func (c FieldChain) Debug(msg string, fields ...Field) {
	Debug(msg, append(c, fields...)...)
}

func (c FieldChain) Fatal(msg string, fields ...Field) {
	Fatal(msg, append(c, fields...)...)
}

func (c FieldChain) Error(msg string, fields ...Field) {
	Error(msg, append(c, fields...)...)
}
