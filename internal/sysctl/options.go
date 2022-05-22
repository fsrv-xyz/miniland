package sysctl

type PropertyOption func(s *Property)

func OptionPathFormatSlash(ps string) PropertyOption {
	return func(property *Property) {
		property.path = encodePathSlashes(ps)
	}
}

func OptionPathFormatDots(ps string) PropertyOption {
	return func(property *Property) {
		property.path = encodePathDots(ps)
	}
}
