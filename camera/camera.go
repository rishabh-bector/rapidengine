package camera

type Camera interface {
	Look()

	MoveUp()
	MoveDown()
	MoveLeft()
	MoveRight()

	GetFirstViewIndex() *float32

	SetPosition(float32, float32)
	GetPosition() (float32, float32)
	SetSpeed(float32)
}
