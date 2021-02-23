package tools

func JudgePosition(newPosX, newPosY, oldPosX, oldPosY float32, heroSpeed float32) bool{
	if getAb(oldPosX - newPosX) > heroSpeed / 1000 * 5 {
		return false
	}

	if getAb(oldPosY - newPosY) > heroSpeed / 1000 * 5 {
		return false
	}
	return true
}
