package projectmd

type ListProjectByCorpIdReqDTO struct {
	CorpId string
	Cursor string
	Limit  int
}

type ListProjectByCorpIdRespDTO struct {
	Data []Project
	Next string
}
