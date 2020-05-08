package uc

type Paginator struct {
	defaultPageSize int
	defaultPage     int
}

func NewPaginator(defaultPageSize, defaultPage int) *Paginator {
	return &Paginator{defaultPageSize, defaultPage}
}

func (pg *Paginator) Paginate(totalSize, pageSize, page int) (offset, limit int) {
	if totalSize == 0 {
		return
	}
	if offset = pageSize * (page - 1); !(offset >= 0 && offset < totalSize) {
		offset, limit = pg.Paginate(totalSize, pg.defaultPageSize, pg.defaultPage)
		return
	}
	if limit = offset + pageSize; limit > totalSize {
		limit = totalSize
		return
	}
	return
}
