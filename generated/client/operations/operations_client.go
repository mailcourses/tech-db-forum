package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// New creates a new operations API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Client {
	return &Client{transport: transport, formats: formats}
}

/*
Client for operations API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

/*
Clear очисткаs всех данных в базе

Безвозвратное удаление всей пользовательской информации из базы данных.

*/
func (a *Client) Clear(params *ClearParams) (*ClearOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewClearParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "clear",
		Method:             "POST",
		PathPattern:        "/service/clear",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ClearReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ClearOK), nil

}

/*
ForumCreate созданиеs форума

Создание нового форума.

*/
func (a *Client) ForumCreate(params *ForumCreateParams) (*ForumCreateCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewForumCreateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "forumCreate",
		Method:             "POST",
		PathPattern:        "/forum/create",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ForumCreateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ForumCreateCreated), nil

}

/*
ForumGetOne получениеs информации о форуме

Получение информации о форуме по его идентификаторе.

*/
func (a *Client) ForumGetOne(params *ForumGetOneParams) (*ForumGetOneOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewForumGetOneParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "forumGetOne",
		Method:             "GET",
		PathPattern:        "/forum/{slug}/details",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ForumGetOneReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ForumGetOneOK), nil

}

/*
ForumGetThreads списокs ветвей обсужления форума

Получение списка ветвей обсужления данного форума.

Ветви обсуждения выводятся отсортированные по дате создания.

*/
func (a *Client) ForumGetThreads(params *ForumGetThreadsParams) (*ForumGetThreadsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewForumGetThreadsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "forumGetThreads",
		Method:             "GET",
		PathPattern:        "/forum/{slug}/threads",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ForumGetThreadsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ForumGetThreadsOK), nil

}

/*
ForumGetUsers пользователиs данного форума

Получение списка пользователей, у которых есть пост или ветка обсуждения в данном форуме.

Пользователи выводятся отсортированные по nickname в порядке возрастания.

*/
func (a *Client) ForumGetUsers(params *ForumGetUsersParams) (*ForumGetUsersOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewForumGetUsersParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "forumGetUsers",
		Method:             "GET",
		PathPattern:        "/forum/{slug}/users",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ForumGetUsersReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ForumGetUsersOK), nil

}

/*
PostGetOne получениеs информации о ветке обсуждения

Получение информации о ветке обсуждения по его имени.

*/
func (a *Client) PostGetOne(params *PostGetOneParams) (*PostGetOneOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostGetOneParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "postGetOne",
		Method:             "GET",
		PathPattern:        "/post/{id}/details",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &PostGetOneReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PostGetOneOK), nil

}

/*
PostUpdate изменениеs сообщения

Изменение сообщения на форуме.

Если сообщение поменяло текст, то оно должно получить отметку `isEdited`.

*/
func (a *Client) PostUpdate(params *PostUpdateParams) (*PostUpdateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostUpdateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "postUpdate",
		Method:             "POST",
		PathPattern:        "/post/{id}/details",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &PostUpdateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PostUpdateOK), nil

}

/*
PostsCreate созданиеs новых постов

Добавление новых постов в ветку обсуждения на форум.

*/
func (a *Client) PostsCreate(params *PostsCreateParams) (*PostsCreateCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostsCreateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "postsCreate",
		Method:             "POST",
		PathPattern:        "/thread/{slug_or_id}/create",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &PostsCreateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*PostsCreateCreated), nil

}

/*
Status получениеs инфомарции о базе данных

Получение инфомарции о базе данных.

*/
func (a *Client) Status(params *StatusParams) (*StatusOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewStatusParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "status",
		Method:             "GET",
		PathPattern:        "/service/status",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &StatusReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*StatusOK), nil

}

/*
ThreadCreate созданиеs ветки

Добавление новой ветки обсуждения на форум.

*/
func (a *Client) ThreadCreate(params *ThreadCreateParams) (*ThreadCreateCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewThreadCreateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "threadCreate",
		Method:             "POST",
		PathPattern:        "/forum/{slug}/create",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ThreadCreateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ThreadCreateCreated), nil

}

/*
ThreadGetOne получениеs информации о ветке обсуждения

Получение информации о ветке обсуждения по его имени.

*/
func (a *Client) ThreadGetOne(params *ThreadGetOneParams) (*ThreadGetOneOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewThreadGetOneParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "threadGetOne",
		Method:             "GET",
		PathPattern:        "/thread/{slug_or_id}/details",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ThreadGetOneReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ThreadGetOneOK), nil

}

/*
ThreadGetPosts сообщенияs данной ветви обсуждения

Получение списка сообщений в данной ветке форуме.

Сообщения выводятся отсортированные по дате создания.

*/
func (a *Client) ThreadGetPosts(params *ThreadGetPostsParams) (*ThreadGetPostsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewThreadGetPostsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "threadGetPosts",
		Method:             "GET",
		PathPattern:        "/thread/{slug_or_id}/posts",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ThreadGetPostsReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ThreadGetPostsOK), nil

}

/*
ThreadUpdate обновлениеs ветки

Обновление ветки обсуждения на форуме.

*/
func (a *Client) ThreadUpdate(params *ThreadUpdateParams) (*ThreadUpdateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewThreadUpdateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "threadUpdate",
		Method:             "POST",
		PathPattern:        "/thread/{slug_or_id}/details",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ThreadUpdateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ThreadUpdateOK), nil

}

/*
ThreadVote проголосоватьs за ветвь обсуждения

Изменение голоса за ветвь обсуждения.

Один пользователь учитывается только один раз и может изменить своё
мнение.

*/
func (a *Client) ThreadVote(params *ThreadVoteParams) (*ThreadVoteOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewThreadVoteParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "threadVote",
		Method:             "POST",
		PathPattern:        "/thread/{slug_or_id}/vote",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &ThreadVoteReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*ThreadVoteOK), nil

}

/*
UserCreate созданиеs нового пользователя

Создание нового пользователя в базе данных.

*/
func (a *Client) UserCreate(params *UserCreateParams) (*UserCreateCreated, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUserCreateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "userCreate",
		Method:             "POST",
		PathPattern:        "/user/{nickname}/create",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &UserCreateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UserCreateCreated), nil

}

/*
UserGetOne получениеs информации о пользователе

Получение информации о пользователе форума по его имени.

*/
func (a *Client) UserGetOne(params *UserGetOneParams) (*UserGetOneOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUserGetOneParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "userGetOne",
		Method:             "GET",
		PathPattern:        "/user/{nickname}/profile",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &UserGetOneReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UserGetOneOK), nil

}

/*
UserUpdate изменениеs данных о пользователе

Изменение информации в профиле пользователя.

*/
func (a *Client) UserUpdate(params *UserUpdateParams) (*UserUpdateOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewUserUpdateParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "userUpdate",
		Method:             "POST",
		PathPattern:        "/user/{nickname}/profile",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"http", "https"},
		Params:             params,
		Reader:             &UserUpdateReader{formats: a.formats},
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	return result.(*UserUpdateOK), nil

}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
