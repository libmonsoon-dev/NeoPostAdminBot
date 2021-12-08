package logger

type Factory interface {
	New(componentName string) Logger
}
