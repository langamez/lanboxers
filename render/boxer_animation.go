package render

import (
	"lanBox/domain"
)

func PlayPunch(punch domain.PunchDetail) {
	var (
		boxer          = punch.Attacker
		mCopyBoxer     = Snapshot(punch.Attacker.BaseBoxer)
		mAffectedParts = map[domain.BodyPart][]domain.Direction{domain.Shoulder: {punch.Direction}, domain.Arm: {punch.Direction}}
	)
	// Set boxer situation direction
	mCopyBoxer.SituationDir = punch.Direction
	// Lock the boxer
	boxer.Lock.Lock()
	defer boxer.Lock.Unlock()
	// Init effect
	mCopyBoxer.Situation = domain.PunchInit
	BoxerFrame(boxer, mCopyBoxer, mAffectedParts)
	Frame(1)
	// Main effect
	mCopyBoxer.Situation = domain.Punch
	BoxerFrame(boxer, mCopyBoxer, mAffectedParts)
	// todo check is it required
	Frame(1)
	// Form back to Idle
	mCopyBoxer.Situation = domain.Idle
	BoxerFrame(boxer, mCopyBoxer, mAffectedParts)
}
