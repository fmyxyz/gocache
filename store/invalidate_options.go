package store

// InvalidateOptions represents the cache invalidation available Options
type InvalidateOptions struct {
	// Tags allows to specify associated tags to the current value
	Tags []string
}

// TagsValue returns the tags option value
func (o InvalidateOptions) TagsValue() []string {
	return o.Tags
}

type InvalidateOption func(o *InvalidateOptions)

func InvalidateTags(tags ...string) InvalidateOption {
	return func(o *InvalidateOptions) {
		o.Tags = tags
	}
}
