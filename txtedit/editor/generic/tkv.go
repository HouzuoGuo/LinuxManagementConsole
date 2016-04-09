package generic

const (
	TKV_NIL_TITLE = "" // denote lacking of a title in key-value pair
)

type GenericTitleKeyValueConf struct {
	TitlePrefix, TitleSuffix string
	QuoteValues              bool
	QuoteKeys                bool
	KVSeparator              string
}

type GenericTitleKeyValue struct {
	Config GenericTitleKeyValueConf
	Values map[GenericNode]map[GenericNode]GenericNode // mapping from title to key-value pairs
}
