package resolver

import t "github.com/ab-dek/golox/token"

type varInfo struct {
	token    t.Token
	resolved bool
	used     bool
}

type scope map[string]*varInfo
