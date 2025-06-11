package idsnow

type IMachineManager interface {
	Register(id string)
	GetMacgineIds() map[string]int
}
