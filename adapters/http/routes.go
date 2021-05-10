package http

const usersAPI = "/users"
const contentsAPI = "/contents"
const ipfsAPI = "/ipfs"

func (h *Handler) SetupRoutes() {
	h.ge.POST(usersAPI, h.RegisterHandler)
	h.ge.GET(usersAPI, h.AuthorizeUser(), h.GetUserInfo)
	h.ge.PUT(usersAPI, h.AuthorizeUser(), h.UserUpdateHandler)
	h.ge.DELETE(usersAPI, h.AuthorizeUser(), h.DeleteUserHandler)

	h.ge.POST(usersAPI+"/login", h.LoginHandler)
	h.ge.POST(usersAPI+"/logout", h.AuthorizeUser(), h.LogoutHandler)

	h.ge.POST(contentsAPI, h.AuthorizeUser(), h.NewContentHandler)
	h.ge.GET(contentsAPI, h.AuthorizeUser(), h.GetContentHandler)
	h.ge.PUT(contentsAPI, h.AuthorizeUser(), h.ContentUpdateHandler)
	h.ge.DELETE(contentsAPI, h.AuthorizeUser(), h.DeleteContentHandler)

	h.ge.POST(contentsAPI+"/rate", h.AuthorizeUser(), h.RateContentHandler)
	h.ge.POST(contentsAPI+"/comment", h.AuthorizeUser(), h.CommentHandler)
	h.ge.GET(contentsAPI+"/comment", h.GetCommentsHandler)

	h.ge.GET(contentsAPI+"/all", h.GetAllContentsHandler)
	h.ge.GET(contentsAPI+"/user", h.AuthorizeUser(), h.GetUserContentsHandler)
	h.ge.POST(contentsAPI+"/search", h.TextSearchHandler)
	// TODO: decide if this needs Auth
	h.ge.GET(ipfsAPI, h.IPFSinfoHandler)
}
