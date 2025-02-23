package shortcodeid

type Repository interface {
	NextID() string
}
