package proc

func jumpUserMode(entry, stack uintptr)

func (t *Task) Run() {
	t.State = StateRunning
	jumpUserMode(t.UserEntry, t.UserStack)
}
