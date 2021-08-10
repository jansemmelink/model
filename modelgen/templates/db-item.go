package system

type amount float64

type DbItem struct {
	ID int64 `json:"id"`
	// Created  time.Time `json:"created"`
	// Modified time.Time `json:"modified"`
}

/*

relationship is either:
	parent-child: if parent is deleted, all child items with this parent are also deleted, and their children, recursively
		e.g. product belongs to a client, client is parent of product
		e.g. order item belongs to an order, when order is deleted, those items are also deleted
		get of parent does not necessarily get children:
		- otherwise getClient will load all client.products, then all product[].variants, and all client.sources etc... huge lot of stuff
		- but getOrder should normally all order items - because it just makes sense ... but support with/without
			load children only when included in parent struct
			so if Client has a []Product field, it will load all products, but it does not have, so it won't
			so if Order has a []OrderItem field, it will load order items, and OrderItem[] must have parent Order to know which to load

		any type can only have one parent
		if no parent (e.g. client) then it is a standalone item that can be created without a parent

		to create a child, the parent must be used, i.e. Client.AddProduct(details), or AddProduct(client, details)
		to create a top level parent, no parent is specified, i.e. AddClient(details)

	reference:
		referenced item cannot be deleted while referred to
		e.g. product refers to an image that may be used by several products or variants
			image cannot be deleted while any product or variant refers to it

			image may still have a parent, e.g. the client
			if has common parent (image->client, product->client, variant->product->client), then only products/variants also belonging to that client, may refer to it
			if no common parent (e.g. images are standalone or in generic image groups, then they can be referred to without belonging to same client)
*/
