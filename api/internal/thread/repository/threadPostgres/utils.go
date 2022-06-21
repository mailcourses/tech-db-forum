package threadPostgres

func prepareFlat(since string, desc bool, limit int64, params []interface{}) (string, []interface{}) {
	var query string
	if desc && since == "" {
		params = append(params, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1
				  ORDER BY created DESC, id DESC
				  LIMIT $2`
	} else if desc && since != "" {
		params = append(params, since, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 AND id < $2
				  ORDER BY created DESC, id DESC
				  LIMIT $3`
	} else if since == "" {
		params = append(params, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1
				  ORDER BY created, id
				  LIMIT $2`
	} else {
		params = append(params, since, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 AND id > $2
				  ORDER BY created, id
				  LIMIT $3`
	}
	return query, params
}

func prepareTree(since string, desc bool, limit int64, params []interface{}) (string, []interface{}) {
	var query string
	if desc && since == "" {
		params = append(params, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1
				  ORDER BY pathTree DESC, id
				  LIMIT $2`
	} else if desc && since != "" {
		params = append(params, since, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 and pathTree < (select pathTree from post where id = $2)
				  ORDER BY pathTree DESC, id
				  LIMIT $3`
	} else if since == "" {
		params = append(params, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1
				  ORDER BY pathTree, id
				  LIMIT $2`
	} else {
		params = append(params, since, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 and pathTree > (select pathTree from post where id = $2)
				  ORDER BY pathTree
				  LIMIT $3`
	}

	return query, params
}

func prepareParentTree(since string, desc bool, limit int64, params []interface{}) (string, []interface{}) {
	var query string
	if desc && since == "" {
		params = append(params, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 and pathTree[1] in (
				  	select pathTree[1] from post where parent = 0 and thread = $1 order by pathTree[1] DESC limit $2
				  )
				  ORDER BY pathTree[1] DESC, pathTree`
	} else if desc && since != "" {
		params = append(params, since, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 and pathTree[1] in (
				  	select pathTree[1] from post where thread = $1 and pathTree[1] < (select pathTree[1] from post where id = $2) limit $3
				  )
				  ORDER BY pathTree[1] DESC, pathTree`
	} else if since == "" {
		params = append(params, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created
				  FROM Post
				  WHERE thread = $1 and pathTree[1] in (
				  	select pathTree[1] from post where parent = 0 and thread = $1 limit $2
				  )
				  ORDER BY pathTree`
	} else {
		params = append(params, since, limit)
		query = `SELECT id, parent, author, message, is_edited, forum, thread, created FROM Post
				  WHERE thread = $1 and pathTree[1] in (
				  	select pathTree[1] from post where thread = $1 and pathTree[1] > (select pathTree[1] from post where id = $2) limit $3
				  )
				  ORDER BY pathTree`
	}
	return query, params
}
