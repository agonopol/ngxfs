package ngxfs

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}