package camera

type Camera interface {
	Look()

	MoveUp()
	MoveDown()
	MoveLeft()
	MoveRight()
	MoveForward()
	MoveBackward()

	ChangeYaw(float32)
	ChangePitch(float32)

	GetFirstViewIndex() *float32

	SetPosition(float32, float32, float32)
	GetPosition() (float32, float32, float32)
	SetSpeed(float32)

	ProcessMouse(float64, float64, float64, float64)
}
