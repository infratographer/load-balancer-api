package api

// parseOriginID parses and validates an origin ID from the request path if the path param is found
// func (r *Router) parseOriginID(c *gin.Context) (string, error) {
// 	origin := struct {
// 		id string `uri:"origin_id" binding:"required"`
// 	}{}

// 	if err := r.parseUUIDFromURI(c, &origin); err != nil {
// 		return "", err
// 	}

// 	if err := r.parseUUID(origin.id); err != nil {
// 		return "", err
// 	}

// 	return origin.id, nil
// }
