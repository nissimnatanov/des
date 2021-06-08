package board

func GetIntegrityChecks() bool {
	return integrityChecks
}

func SetIntegrityChecks(enabled bool) bool {
	prev := integrityChecks
	integrityChecks = enabled
	return prev
}

var integrityChecks bool
