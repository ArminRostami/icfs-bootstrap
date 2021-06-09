package http

const usersAPI = "/users"
const contentsAPI = "/contents"
const ipfsAPI = "/ipfs"
const icfsAPI = "/icfs"

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
	h.ge.DELETE(contentsAPI+"/downloads", h.AuthorizeUser(), h.DeleteDownloadHandler)

	h.ge.POST(contentsAPI+"/review", h.AuthorizeUser(), h.ReviewContentHandler)
	h.ge.GET(contentsAPI+"/comment", h.GetCommentsHandler)

	h.ge.GET(contentsAPI+"/all", h.GetAllContentsHandler)
	h.ge.GET(contentsAPI+"/uploads", h.AuthorizeUser(), h.GetUserUploadsHandler)
	h.ge.GET(contentsAPI+"/downloads", h.AuthorizeUser(), h.GetUserDownloadsHandler)
	h.ge.POST(contentsAPI+"/search", h.TextSearchHandler)
	h.ge.GET(ipfsAPI, h.IPFSinfoHandler)

	h.ge.GET(icfsAPI, h.ICFSServer)

	h.ge.NoRoute(h.UIhandler)
}
