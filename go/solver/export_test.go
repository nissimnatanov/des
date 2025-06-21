package solver

func DisableNLog(disable bool) bool {
	old := disableNLog
	disableNLog = disable
	return old
}

const TrialAndErrorStepName = trialAndErrorStepName
