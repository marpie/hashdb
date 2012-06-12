package hashdb

type GetResponse struct {
	hash     string
	password string
	err      error
}

type GetRequest struct {
	request  string
	response chan *GetResponse
}

type PutResponse error

type PutRequest struct {
	hash     string
	password string
	response chan PutResponse
}
