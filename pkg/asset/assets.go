package asset

// GetAssetKeyMap creates a map mapping asset IDs to quantities.
func (ctx *Context) GetAssetKeyMap() map[string]*Asset {
	m := map[string]*Asset{}

	for i := range ctx.Assets {
		m[ctx.Assets[i].Id] = ctx.Assets[i]
	}

	return m
}
