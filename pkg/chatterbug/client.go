package chatterbug

import (
	"context"
	"net/http"

	graphql "github.com/hasura/go-graphql-client"
)

type Chatterbug struct {
	client *graphql.Client
}

type OrganizationMember struct {
	ID        string
	Name      string
	Login     string
	AvatarURL string
	Timezone  string
}

// GetOrganizationMembers fetches the list of organization members using the
// Chatterbug GraphQL query inteface.
func (c Chatterbug) GetOranizationMembers() ([]OrganizationMember, error) {
	var q struct {
		CurrentUser struct {
			ID           graphql.ID
			Name         graphql.String
			Organization struct {
				Id      graphql.String
				Name    graphql.String
				Members []OrganizationMember
			}
		}
	}

	if err := c.client.Query(context.Background(), &q, nil); err != nil {
		logger.WithError(err).Error("failed to lookup organization members")
		return nil, err
	}

	return q.CurrentUser.Organization.Members, nil
}

type Leader struct {
	User struct {
		Login string
	}

	StudyTime float32
	Rank      int
}

type Leaderboard struct {
	Name           string
	CountOfMembers int
	TopPositions   []Leader `graphql:"topPositions(limit: $limit)"`
}

// GetLeaderboard gets the current leaderboard for the authenticated user's
// organization using the Chatterbug GraphQL query interface.
func (c Chatterbug) GetLeaderboard() (Leaderboard, error) {
	var q struct {
		LeaderBoard Leaderboard
	}

	vars := map[string]interface{}{
		"limit": graphql.Int(10),
	}

	if err := c.client.Query(context.Background(), &q, vars); err != nil {
		logger.WithError(err).Error("failed to lookup organization members")
		return Leaderboard{}, err
	}

	return q.LeaderBoard, nil
}

func New(tok string) *Chatterbug {
	return &Chatterbug{
		client: graphql.NewClient("https://app.chatterbug.com/api/graphql", &http.Client{
			Transport: newChatterbugRoundTripper(tok),
		}),
	}
}
