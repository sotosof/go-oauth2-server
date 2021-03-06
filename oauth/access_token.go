package oauth

import (
	"time"

	"github.com/RichardKnop/go-oauth2-server/util"
)

// GrantAccessToken deletes old tokens and grants a new access token
func (s *Service) GrantAccessToken(client *Client, user *User, scope string) (*AccessToken, error) {
	// Delete expired access tokens
	s.deleteExpiredAccessTokens(client, user)

	// Create a new access token
	accessToken := newAccessToken(
		s.cnf.Oauth.AccessTokenLifetime, // expires in
		client, // client
		user,   // user
		scope,  // scope
	)
	if err := s.db.Create(accessToken).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}

// deleteExpiredAccessTokens deletes expired access tokens
func (s *Service) deleteExpiredAccessTokens(client *Client, user *User) {
	s.db.Unscoped().Where(AccessToken{
		ClientID: util.IntOrNull(int64(client.ID)),
		UserID:   util.IntOrNull(int64(user.ID)),
	}).Where("expires_at <= ?", time.Now()).Delete(new(AccessToken))
}
