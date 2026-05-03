package models

type Response struct {
	Success bool
}

type ResponseBody struct {
	Success bool
	Body    interface{}
}

type PaginatedBody struct {
	Success bool
	Body    interface{}
	Next    int
	Prev    int
}
