package gelly

// v1
// type Key[T any] struct {
// 	id  int
// 	gen uint64
// }

// type cell struct {
// 	nextFree int
// 	gen      uint64
// 	itv      int
// 	vti      int
// }

// type Pool[T any] struct {
// 	gen          uint64
// 	nextFreeCell int
// 	lenCells     int
// 	capValues    int
// 	nextValue    int
// 	cells        []cell
// 	Values       []T
// }

// func (p *Pool[T]) Create(val T) Key[T] {
// 	if p.nextFreeCell >= p.lenCells {
// 		p.lenCells++
// 		p.cells = append(p.cells, cell{nextFree: p.lenCells})
// 	}

// 	current := p.nextFreeCell
// 	p.cells[current].gen = p.gen

// 	if p.nextValue == p.capValues {
// 		p.Values = append(p.Values, val)
// 		p.capValues++
// 	} else {
// 		p.Values[p.nextValue] = val
// 	}
// 	p.cells[current].itv = p.nextValue
// 	p.cells[p.nextValue].vti = current

// 	p.nextValue++
// 	p.nextFreeCell = p.cells[current].nextFree

// 	return Key[T]{id: current, gen: p.gen}
// }

// func (p *Pool[T]) Set(key Key[T], val T) bool {
// 	if key.id >= p.lenCells {
// 		return false
// 	}

// 	if key.gen != p.cells[key.id].gen {
// 		return false
// 	}

// 	valueIndex := p.cells[key.id].itv
// 	p.Values[valueIndex] = val
// 	return true
// }

// func (p *Pool[T]) Get(key Key[T]) (*T, bool) {
// 	if key.id >= p.lenCells {
// 		return nil, false
// 	}

// 	if key.gen != p.cells[key.id].gen {
// 		return nil, false
// 	}

// 	valueIndex := p.cells[key.id].itv
// 	return &p.Values[valueIndex], true
// }

// func (p *Pool[T]) Destroy(key Key[T]) bool {
// 	if key.id >= p.lenCells {
// 		return false
// 	}

// 	if key.gen != p.cells[key.id].gen {
// 		return false
// 	}

// 	p.gen++
// 	p.cells[key.id].gen = p.gen
// 	p.cells[key.id].nextFree = p.nextFreeCell
// 	p.nextFreeCell = key.id

// 	valueIndexToSwap := p.cells[key.id].itv
// 	p.nextValue--
// 	p.Values[valueIndexToSwap] = p.Values[p.nextValue]

// 	valueIndex := p.cells[p.nextValue].vti
// 	p.cells[valueIndex].itv = valueIndexToSwap
// 	p.cells[valueIndexToSwap].vti = valueIndex

// 	p.Values = p.Values[:p.nextValue]
// 	return true
// }

// signature int64
// blueprint EntityBlueprint
// children...

// type Entity interface {
// 	Init(cmd *Command, k Key)
// 	Dispose(cmd *Command, k Key)
// }

// // AddComponent(cmd *Command, k Key)
// // GetComponent(cmd *Command, k Key)
// // RemoveComponent(cmd *Command, k Key)
// // GetSigature(cmd *Command, k Key)

// type EntityBase struct {
// }

// type Event interface {
// 	IsEvent()
// }

// ARENA START
/*
type Key struct {
	id  int
	gen uint64
}

type cell[T any] struct {
	next int
	gen  uint64
	val  T
}

type Arena[T any] struct {
	next  int
	len   int
	gen   uint64
	cells []cell[T]
	Data  []T
}

func (a *Arena[T]) Insert(val T) Key {
	if a.next >= a.len {
		a.len++
		a.cells = append(a.cells, cell[T]{next: a.len})
	}

	a.cells[a.next].val = val
	a.cells[a.next].gen = a.gen

	next := a.next
	key := Key{id: next, gen: a.gen}
	a.next = a.cells[next].next

	a.cells[next].next = -1

	return key
}

func (a *Arena[T]) Set(k Key, v T) bool {
	if k.id >= a.len {
		return false
	}

	cell := a.cells[k.id]
	if cell.gen != k.gen {
		return false
	}

	a.cells[k.id].val = v

	return true
}

func (a *Arena[T]) Get(k Key) (*T, bool) {
	// var v T

	if k.id >= a.len {
		// return v, false
		return nil, false
	}

	cell := a.cells[k.id]
	if cell.gen != k.gen {
		// return v, false
		return nil, false
	}

	return &cell.val, true
}

func (a *Arena[T]) Remove(k Key) bool {
	if k.id >= a.len {
		return false
	}

	if a.cells[k.id].gen != k.gen {
		return false
	}

	a.cells[k.id].next = a.next
	a.next = k.id

	a.gen++
	a.cells[k.id].gen = a.gen

	return true
}

func (a *Arena[T]) ForEach(fn func(k Key, v T)) {
	for i := 0; i < len(a.cells); i++ {
		if a.cells[i].next < 0 {
			fn(Key{id: i, gen: a.gen}, a.cells[i].val)
		}
	}
}
*/
// END

// type Rigidbody struct {
// }

// func (c Rigidbody) Type() ComponentType {
// 	return RIGIDBODY_COMPONENT
// }

// type Velocity struct {
// 	X int
// 	Y int
// }

// func (v Velocity) Type() ComponentType {
// 	return VELOCITY_COMPONENT
// }

// type Component interface {
// 	Type() ComponentType
// }

// type actionType int

// const (
// 	createEntity actionType = iota
// 	addComponent
// )

// type commandAction struct {
// 	t actionType
// 	k Key
// 	e Entity
// 	c Component
// }

// type Command struct {
// 	actions []commandAction
// }

// func (cmd *Command) CreateEntity(e Entity) {
// 	// cmd.actions = append(cmd.actions, commandAction{
// 	// 	t: createEntity,
// 	// 	e: e,
// 	// })
// }

// // func DestroyEntity(cmd *Command, entity Entity) {
// // }

// func (cmd *Command) AddComponent(k Key, c Component) {
// 	// cmd.actions = append(cmd.actions, commandAction{
// 	// 	t: addComponent,
// 	// 	k: k,
// 	// 	c: c,
// 	// })
// }

// // func GetComponent[C Component](cmd *Command, entity Entity) *C {
// // 	return nil
// // }

// // func RemoveComponent[C Component](cmd *Command, entity Entity) {
// // }

// // // func AddChildren(cmd *Command, parent Entity, child EntityBlueprint) Entity {
// // // 	return Entity{}
// // // }

// // // func Send[E Event](cmd *Command, from Entity, to Entity, event E) {
// // // }

// // func Publish[E Event](cmd *Command, entity Entity, event E) {
// // }

// // func Subscribe[E Event](cmd *Command, entity Entity, h func(cmd *Command, entity Entity, event E)) {
// // }

// // func Unsubcribe[E Event](cmd *Command, entity Entity) {
// // }

// type World struct {
// 	arena arena
// 	// pools    [64]any

// 	frontCmd *Command
// 	backCmd  *Command
// }

// func (w *World) Init(setup func(cmd *Command)) {
// 	w.frontCmd = &Command{}
// 	w.backCmd = &Command{}

// 	setup(w.frontCmd)

// 	for _, action := range w.frontCmd.actions {
// 		switch action.t {

// 		case createEntity:
// 			// TODO: maybe need to store this entity
// 			// TODO: maybe need two queues for commands ?
// 			key := w.arena.insert(action.e)
// 			action.e.Init(w.backCmd, key)

// 			// TODO: may not need to store components instead pool but just on the entity
// 		case addComponent:

// 			log.Println("Should try to add the component")
// 		}
// 	}

// 	log.Println("BACK", w.backCmd)
// }

// // TODO: maybe convert Event to Message for same thign
// func (w *World) Write(event Event) {
// }

// // TODO: maybe send and event to all of them and switch to message instead
// func (w *World) Update(c *Client, dt time.Duration) {
// }

// func (w *World) Draw(r *ebiten.Image, dt time.Duration) {
// }

// func (w *World) Dispose() {
// }

// // SYSTEM

// // type System interface {
// // 	Init() Access
// // 	Join(e Key)
// // 	Leave(e Key)
// // 	Execute(cmd *Command, entities []Key)
// // 	Dispose()
// // }

// //	type Access struct {
// //		Include   []Component
// //		Exclude   []Component
// //		Resources []Resource
// //	}

// // type Read[T any] struct {
// // }

// // type Write[T any] struct {
// // }

// // func (r Read[T]) Mode()  {}
// // func (w Write[T]) Mode() {}

// // type Mode[T any] interface {
// // 	Read[T] | Write[T]
// // }

// // func TryInsert[T Mode[T]](m T) {
// // }

// // internal
// // type entity struct {
// // 	flags uint64
// // }

// // type pool[T any] struct {
// // 	next int
// // 	itd  []int
// // 	dti  []int
// // 	data []T
// // }

// // func newPool[T any]() *pool[T] {
// // 	return &pool[T]{
// // 		next: 0,
// // 		itd:  []int{},
// // 		dti:  []int{},
// // 		data: []T{},
// // 	}
// // }

// // // TODO: doesn't check if data already in there
// // func (p *pool[T]) add(index int, value T) {
// // 	length := len(p.itd)

// // 	if index >= length {
// // 		diff := (index + 1) - length

// // 		for i := 0; i < diff; i++ {
// // 			p.itd = append(p.itd, 0)
// // 		}
// // 	}

// // 	p.itd[index] = p.next

// // 	if p.next >= len(p.data) {
// // 		p.data = append(p.data, value)
// // 		p.dti = append(p.dti, p.next)
// // 	} else {
// // 		p.data[p.next] = value
// // 		p.dti[p.next] = index
// // 	}

// // 	p.next++
// // }

// // func (p *pool[T]) get(index int) *T {
// // 	// TODO: maybe check out of bound
// // 	return &p.data[p.itd[index]]
// // }

// // func (p *pool[T]) rem(index int) {
// // 	// TODO: maybe check out of bound
// // 	// TODO: check trying to remove empty

// // 	last := p.next - 1
// // 	indirect := p.itd[index]
// // 	previous := p.dti[last]

// // 	p.dti[indirect] = previous
// // 	p.itd[previous] = last
// // 	p.data[indirect] = p.data[last]

// // 	p.next--
// // }

// // func createWorld() *World {
// // 	return &World{
// // 		entities: NewArena[entity](),
// // 		pools:    [64]any{},
// // 	}
// // }

// // func createentity(w *World) Key {
// // 	return w.entities.Insert(entity{Flags: 0})
// // }

// // func destroyentity(w *World, entity Key) {
// // 	w.entities.Remove(entity)
// // 	// TODO: also remove components
// // }

// // func RegisterComponent[C Component](w *World) {
// // 	var c C
// // 	w.pools[c.Type()] = newpool[C]()
// // }

// // func addComponent[C Component](w *World, entity Key, component C) {
// // 	// TODO: maybe do not need checking
// // 	e, ok := w.entities.Get(entity)
// // 	if !ok {
// // 		log.Fatal("entity does not exist")
// // 	}

// // 	t := component.Type()

// // 	// TODO: maybe get back a pointer instead
// // 	e.Flags |= (1 << t)
// // 	w.entities.Set(entity, e)

// // 	pool := w.pools[t].(*pool[C])
// // 	pool.add(entity.id, component)
// // }

// // func handler[M Message](w *World, fn func(e Key, m M)) Key {
// // 	return Key{}
// // }

// // func subscribe(w *World, entity Key, handler Key) {
// // }

// // func getComponent[C Component](w *World, entity Key) *C {
// // 	_, ok := w.entities.Get(entity)
// // 	if !ok {
// // 		log.Fatal("entity does not exist")
// // 	}

// // 	// TODO: find a way to replace that
// // 	var c C
// // 	pool := w.pools[c.Type()].(*pool[C])
// // 	return pool.get(entity.id)
// // }

// // type World struct {
// // }

// // func CreateWorld() *World {
// // 	return &World{}
// // }

// // func Createentity(w *World) Key {
// // 	return Key{}
// // }

// // type ComponentType int

// // type Component interface {
// // 	Type() ComponentType
// // }

// // func AddComponent(w *World, entity Key, component Component) {
// // }

// // func Destroyentity(w *World, entity Key) bool {
// // 	return false
// // }

// // type Key struct {
// // }

// // type entity struct {
// // }
