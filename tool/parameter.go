package tool

type CheckParameter interface {
	ForInter(interface{}) bool
}

type Check struct {
}

func (c *Check) ForInter(i interface{}) bool {
	if i == nil {
		return false
	}
	return true
}
