package system

import "context"

//Context keeps track of what was already read and does not read the same item again
func NewContext(base context.Context) *Context {
	return &Context{
		Context:         base,
		itemByTypeAndID: map[string]map[ID]IItem{},
	}
}

type Context struct {
	context.Context
	itemByTypeAndID map[string]map[ID]IItem //index on ext type name and id then store pointer to item struct, e.g. ["client"][21] = &Client{DbItem:DbItem{ID:21}, ...}
}

func (c *Context) addItem(extTypeName string, item IItem) {
	if extTypeName != "" && item != nil {
		itemByID, ok := c.itemByTypeAndID[extTypeName]
		if !ok {
			itemByID = map[ID]IItem{}
			c.itemByTypeAndID[extTypeName] = itemByID
		}
		itemByID[item.GetID()] = item
	}
}

func (c *Context) getItem(extTypeName string, id ID) (IItem, bool) {
	itemByID, ok := c.itemByTypeAndID[extTypeName]
	if !ok {
		return nil, false
	}
	item, ok := itemByID[id]
	if !ok {
		return nil, false
	}
	return item, true
}
