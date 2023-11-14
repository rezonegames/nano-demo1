package room

import (
	"math/rand"
	"tetris/internal/game/util"
	"tetris/proto/proto"
	"time"
)

type Robot1 struct {
	client util.ClientEntity
	table  util.TableEntity
	id     int64
}

func NewRobot1(client util.ClientEntity, table util.TableEntity) *Robot1 {
	return &Robot1{
		client: client,
		table:  table,
		id:     client.GetId(),
	}
}

func (r *Robot1) fill0() *proto.Array2 {
	v := &proto.Array2{
		Rows: make([]*proto.Row, 0),
	}
	for i := 0; i < 20; i++ {
		values := make([]int32, 0)
		for j := 0; j < 12; j++ {
			values = append(values, 0)
		}
		v.Rows = append(v.Rows, &proto.Row{Values: values})
	}
	return v
}

func (r *Robot1) createPiece() *proto.Array2 {
	all := "ILJOTSZ"
	v := &proto.Array2{}
	t := all[rand.Intn(7)]
	switch t {
	case 'T':
		v.Rows = []*proto.Row{
			{Values: []int32{0, 0, 0}},
			{Values: []int32{1, 1, 1}},
			{Values: []int32{0, 1, 0}},
		}
	case 'O':
		v.Rows = []*proto.Row{
			{Values: []int32{2, 2}},
			{Values: []int32{2, 2}},
		}
	case 'L':
		v.Rows = []*proto.Row{
			{Values: []int32{0, 3, 0}},
			{Values: []int32{0, 3, 0}},
			{Values: []int32{0, 3, 3}},
		}
	case 'J':
		v.Rows = []*proto.Row{
			{Values: []int32{0, 4, 0}},
			{Values: []int32{0, 4, 0}},
			{Values: []int32{4, 4, 0}},
		}
	case 'I':
		v.Rows = []*proto.Row{
			{Values: []int32{0, 5, 0, 0}},
			{Values: []int32{0, 5, 0, 0}},
			{Values: []int32{0, 5, 0, 0}},
			{Values: []int32{0, 5, 0, 0}},
		}
	case 'S':
		v.Rows = []*proto.Row{
			{Values: []int32{0, 6, 6}},
			{Values: []int32{6, 6, 0}},
			{Values: []int32{0, 0, 0}},
		}
	case 'Z':
		v.Rows = []*proto.Row{
			{Values: []int32{7, 7, 0}},
			{Values: []int32{0, 7, 7}},
			{Values: []int32{0, 0, 0}},
		}
	}
	return v
}

func (r *Robot1) collide(s *proto.UpdateState) bool {
	player, arena := s.Player, s.Arena
	m, o := player.Matrix, player.Pos
	for y := 0; y < len(m.Rows); y++ {
		for x := 0; x < len(m.Rows[y].Values); x++ {
			oy, ox := y+int(o.Y), x+int(o.X)
			//log.Info("collide %d %d", oy, ox)
			if m.Rows[y].Values[x] != 0 &&
				(oy >= len(arena.Matrix.Rows) || arena.Matrix.Rows[oy].Values[ox] != 0) {
				return true
			}
		}
	}
	return false
}

func (r *Robot1) merge(s *proto.UpdateState) {
	player, arena := s.Player, s.Arena
	for y, row := range player.Matrix.Rows {
		for x, value := range row.Values {
			if value != 0 {
				arena.Matrix.Rows[y+int(player.Pos.Y)].Values[x+int(player.Pos.X)] = value
			}
		}
	}
}

func (r *Robot1) reset(s *proto.UpdateState) {
	player, arena := s.Player, s.Arena
	player.Matrix = r.createPiece()
	player.Pos = &proto.Pos{
		X: uint32(len(arena.Matrix.Rows[0].Values)/2 - len(player.Matrix.Rows[0].Values)/2),
		Y: 0,
	}
	if r.collide(s) {
		s.End = true
	}
}

func (r *Robot1) Update(t1 time.Time, t2 time.Time) error {
	d := t2.Sub(t1)
	if false {
		//if d > time.Second {
		p := r.client.GetPlayer()
		s := &proto.UpdateState{
			Fragment: "all",
			Player:   p.State.Player,
			Arena:    p.State.Arena,
			PlayerId: r.id,
			End:      p.End,
		}
		//
		// 如果已经结束了，返回
		if s.End {
			return nil
		}

		//
		// 如果刚上线就掉了， 初始化之
		if s.Arena.Matrix == nil {
			s.Arena.Matrix = r.fill0()
			r.reset(s)
		}

		//
		// 判断过了几秒，并且pos.y ++
		nSeconds := int(d.Seconds())
		for i := 0; i < nSeconds; i++ {
			s.Player.Pos.Y++
			//
			// 碰撞了，pos.y--
			if r.collide(s) {
				s.Player.Pos.Y--
				r.merge(s)
				r.reset(s)
				//
				// 如果结束了，就放弃循环
				if s.End {
					break
				}
			}
		}

		//
		// 广播之
		r.table.BroadcastPlayerState(r.id, s)
	}

	return nil
}
