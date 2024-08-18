package gelly

type Key struct {
	id  int
	gen uint
	ok  bool
}

type Value[T any] struct {
	K Key
	V T
}

type cell struct {
	nfr int
	itv int
	gen uint
}

type Pool[T any] struct {
	nfr      int
	gen      uint
	len      int
	creates  []Value[T]
	destroys []Key
	cells    []cell
	Values   []Value[T]
}

func (p *Pool[T]) Create(v T) Key {
	k := Key{id: p.nfr, gen: p.gen, ok: true}
	p.creates = append(p.creates, Value[T]{
		K: k,
		V: v,
	})

	if p.nfr >= p.len {
		p.cells = append(p.cells, cell{itv: -1, gen: p.gen})
		p.len++
		p.nfr++
		return k
	}

	p.cells[p.nfr].itv = -1
	p.cells[p.nfr].gen = p.gen
	p.nfr = p.cells[p.nfr].nfr
	return k
}

func (p *Pool[T]) Get(k Key) (*T, bool) {
	if !p.valid(k) {
		return nil, false
	}

	return &p.Values[p.cells[k.id].itv].V, true
}

func (p *Pool[T]) Destroy(k Key) bool {
	if !p.valid(k) {
		return false
	}

	p.gen++
	p.destroys = append(p.destroys, k)
	return true
}

func (p *Pool[T]) For(h func(i int, k Key, v *T) bool) {
	for i := range p.Values {
		br := h(i, p.Values[i].K, &p.Values[i].V)
		if br {
			break
		}
	}
}

func (p *Pool[T]) Apply() {
	for _, k := range p.destroys {
		p.cells[k.id].nfr = p.nfr
		p.cells[k.id].gen = p.gen
		p.nfr = k.id

		last := p.Values[len(p.Values)-1]
		p.Values[p.cells[k.id].itv] = last
		p.cells[last.K.id].itv = p.cells[k.id].itv
		p.Values = p.Values[:len(p.Values)-1]
	}
	p.destroys = p.destroys[:0]

	for _, create := range p.creates {
		p.cells[create.K.id].itv = len(p.Values)
		p.Values = append(p.Values, create)
	}
	p.creates = p.creates[:0]
}

func (p *Pool[T]) valid(k Key) bool {
	if !k.ok {
		return false
	}

	if k.id >= p.len {
		return false
	}

	if p.cells[k.id].gen != k.gen {
		return false
	}

	if p.cells[k.id].itv < 0 {
		return false
	}

	return true
}

// // Pool
// type Key[T any] struct {
// 	id  int
// 	gen uint
// 	ok  bool
// }

// type Value[T any] struct {
// 	K Key
// 	V T
// 	T Transform
// }

// type cell struct {
// 	nfr int
// 	itv int
// 	gen uint
// }

// type Pool[T any] struct {
// 	nfr      int
// 	gen      uint
// 	len      int
// 	creates  []Value[T]
// 	destroys []Key
// 	cells    []cell
// 	Values   []Value[T]
// }

// // func (p *Pool[T]) SetTransform(k Key, t Transform) {
// // }

// func (p *Pool[T]) Create(v T) Key {
// 	k := Key{id: p.nfr, gen: p.gen, ok: true}
// 	p.creates = append(p.creates, Value[T]{
// 		K: k,
// 		V: v,
// 	})

// 	if p.nfr >= p.len {
// 		p.cells = append(p.cells, cell{itv: -1, gen: p.gen})
// 		p.len++
// 		p.nfr++
// 		return k
// 	}

// 	p.cells[p.nfr].itv = -1
// 	p.cells[p.nfr].gen = p.gen
// 	p.nfr = p.cells[p.nfr].nfr
// 	return k
// }

// func (p *Pool[T]) Get(k Key) (*T, bool) {
// 	if !p.valid(k) {
// 		return nil, false
// 	}

// 	return &p.Values[p.cells[k.id].itv].V, true
// }

// func (p *Pool[T]) Destroy(k Key) bool {
// 	if !p.valid(k) {
// 		return false
// 	}

// 	p.gen++
// 	p.destroys = append(p.destroys, k)
// 	return true
// }

// func (p *Pool[T]) For(h func(k Key, v *T) bool) {
// 	for i := range p.Values {
// 		br := h(p.Values[i].K, &p.Values[i].V)
// 		if br {
// 			break
// 		}
// 	}
// }

// func (p *Pool[T]) Apply() {
// 	for _, k := range p.destroys {
// 		p.cells[k.id].nfr = p.nfr
// 		p.cells[k.id].gen = p.gen
// 		p.nfr = k.id

// 		last := p.Values[len(p.Values)-1]
// 		p.Values[p.cells[k.id].itv] = last
// 		p.cells[last.K.id].itv = p.cells[k.id].itv
// 		p.Values = p.Values[:len(p.Values)-1]
// 	}
// 	p.destroys = p.destroys[:0]

// 	for _, create := range p.creates {
// 		p.cells[create.K.id].itv = len(p.Values)
// 		p.Values = append(p.Values, create)
// 	}
// 	p.creates = p.creates[:0]
// }

// func (p *Pool[T]) valid(k Key) bool {
// 	if !k.ok {
// 		return false
// 	}

// 	if k.id >= p.len {
// 		return false
// 	}

// 	if p.cells[k.id].gen != k.gen {
// 		return false
// 	}

// 	if p.cells[k.id].itv < 0 {
// 		return false
// 	}

// 	return true
// }
