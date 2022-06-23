package threadPostgres

func prepareFlat(since string, desc bool, limit int64, params []interface{}) (string, []interface{}) {
	var query string
	if desc && since == "" {
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1
				  ORDER BY created DESC, id DESC
				  LIMIT $2`
	} else if desc && since != "" {
		params = append(params, since)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 AND id < $2
				  ORDER BY created DESC, id DESC
				  LIMIT $3`
	} else if since == "" {
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1
				  ORDER BY created, id
				  LIMIT $2`
	} else {
		params = append(params, since)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 AND id > $2
				  ORDER BY created, id
				  LIMIT $3`
	}
	params = append(params, limit)
	return query, params
}

//GetPostsOnThreadTreeCommand                    = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND parent_path > (SELECT parent_path FROM Posts WHERE id = $2) ORDER BY parent_path, id LIMIT $3;"
//GetPostsOnThreadTreeDescCommand                = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM Posts WHERE thread = $1 AND parent_path < (SELECT parent_path FROM Posts WHERE id = $2) ORDER BY parent_path DESC LIMIT $3;"

func prepareTree(since string, desc bool, limit int64, params []interface{}) (string, []interface{}) {
	var query string
	if desc && since == "" {
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1
				  ORDER BY pathTree DESC
				  LIMIT $2;`
	} else if desc && since != "" {
		params = append(params, since)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 and pathTree < (select pathTree from post where id = $2)
				  ORDER BY pathTree DESC
				  LIMIT $3;`
	} else if since == "" {
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1
				  ORDER BY pathTree, id
				  LIMIT $2;`
	} else {
		params = append(params, since)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 and pathTree > (select pathTree from post where id = $2)
				  ORDER BY pathTree, id
				  LIMIT $3;`
	}
	params = append(params, limit)
	return query, params
}

func prepareParentTree(since string, desc bool, limit int64, params []interface{}) (string, []interface{}) {
	var query string
	if desc && since == "" {
		params = append(params, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE pathTree[1] in (
				  	select id from post where parent = 0 and thread = $1 order by id DESC limit $2
				  )
				  ORDER BY pathTree[1] DESC, pathTree, id`
	} else if desc && since != "" {
		params = append(params, since, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE pathTree[1] in (
				  	select id from post where thread = $1 and parent = 0 and id < (select pathTree[1] from post where id = $2) order by id DESC limit $3
				  )
				  ORDER BY pathTree[1] DESC, pathTree, id`
	} else if since == "" {
		params = append(params, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created
				  FROM Post
				  WHERE pathTree[1] in (
				  	select id from post where parent = 0 and thread = $1 order by id limit $2
				  )
				  ORDER BY pathTree, id`
	} else {
		params = append(params, since, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE pathTree[1] in (
				  	select id from post where thread = $1 and parent = 0 and id > (select pathTree[1] from post where id = $2) order by id limit $3
				  )
				  ORDER BY pathTree, id`
	}
	return query, params
}
