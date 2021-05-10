package http

const usersAPI = "/users"
const contentsAPI = "/contents"
const ipfsAPI = "/ipfs"

func (h *Handler) SetupRoutes() {
	h.ge.POST(usersAPI, h.RegisterHandler)
	h.ge.GET(usersAPI, h.AuthorizeJWT(), h.GetUserInfo)
	h.ge.PUT(usersAPI, h.AuthorizeJWT(), h.UserUpdateHandler)
	h.ge.DELETE(usersAPI, h.AuthorizeJWT(), h.DeleteUserHandler)

	h.ge.POST(usersAPI+"/login", h.LoginHandler)

	h.ge.POST(contentsAPI, h.AuthorizeJWT(), h.NewContentHandler)
	h.ge.GET(contentsAPI, h.AuthorizeJWT(), h.GetContentHandler)
	h.ge.PUT(contentsAPI, h.AuthorizeJWT(), h.ContentUpdateHandler)
	h.ge.DELETE(contentsAPI, h.AuthorizeJWT(), h.DeleteContentHandler)

	h.ge.POST(contentsAPI+"/rate", h.AuthorizeJWT(), h.RateContentHandler)
	h.ge.POST(contentsAPI+"/comment", h.AuthorizeJWT(), h.CommentHandler)
	h.ge.GET(contentsAPI+"/comment", h.GetCommentsHandler)

	h.ge.GET(contentsAPI+"/all", h.GetAllContentsHandler)
	h.ge.GET(contentsAPI+"/user", h.AuthorizeJWT(), h.GetUserContentsHandler)
	h.ge.POST(contentsAPI+"/search", h.TextSearchHandler)
	// TODO: decide if this need Auth
	h.ge.GET(ipfsAPI, h.IPFSinfoHandler)
}
