package distributed_lock

type Locker interface {
	Lock()
	UnLock()
}
