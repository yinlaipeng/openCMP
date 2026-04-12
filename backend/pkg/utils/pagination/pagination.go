package pagination

// Pagination represents pagination information
type Pagination struct {
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
	Total   int64       `json:"total"`
	Pages   int         `json:"pages"`
	HasNext bool        `json:"has_next"`
	Items   interface{} `json:"items"`
}

// GetOffset calculates the offset for database queries
func (p *Pagination) GetOffset() int64 {
	if p.Page <= 1 {
		return 0
	}
	return int64((p.Page - 1) * p.Limit)
}

// CalculatePages calculates the total number of pages
func (p *Pagination) CalculatePages() {
	if p.Limit <= 0 {
		p.Pages = 0
		return
	}
	p.Pages = int((p.Total + int64(p.Limit) - 1) / int64(p.Limit))
	p.HasNext = p.Page < p.Pages
}
