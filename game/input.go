package game

import "github.com/Bredgren/geo"

type input struct {
	Left      bool
	Right     bool
	Move      float64
	Jump      bool
	Punch     bool
	PunchAxis geo.Vec
	Launch    bool
	Slam      bool
}

func (i *input) handleLeft(down bool) bool {
	i.Left = down
	return false
}

func (i *input) handleRight(down bool) bool {
	i.Right = down
	return false
}

func (i *input) handleMove(val float64) bool {
	i.Move = val
	return false
}

func (i *input) handleJump(down bool) bool {
	i.Jump = down
	return false
}

func (i *input) handlePunch(down bool) bool {
	i.Punch = down
	return false
}

func (i *input) handlePunchH(val float64) bool {
	i.PunchAxis.X = val
	return false
}

func (i *input) handlePunchV(val float64) bool {
	i.PunchAxis.Y = -val
	return false
}

func (i *input) handleLaunch(down bool) bool {
	i.Launch = down
	return false
}

func (i *input) handleSlam(down bool) bool {
	i.Slam = down
	return false
}
