package handlers

func RegisterAll(r *Registry) {
	r.Register("send_email", sendEmail)
	r.Register("resize_photo", resizePhoto)
}
