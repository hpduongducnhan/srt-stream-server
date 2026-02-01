package app

type Level string

const (
	OK       Level = "🟢 OK"
	WARNING  Level = "🟡 WARNING"
	CRITICAL Level = "🔴 CRITICAL"
)

func levelByPercent(v float64) Level {
	switch {
	case v >= 85:
		return CRITICAL
	case v >= 70:
		return WARNING
	default:
		return OK
	}
}

func levelByLoad(load float64, cpu int) Level {
	ratio := load / float64(cpu)
	switch {
	case ratio > 2:
		return CRITICAL
	case ratio > 1:
		return WARNING
	default:
		return OK
	}
}
