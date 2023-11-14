package room

import (
	"math/rand"
	"tetris/proto/proto"
)

func fill0() *proto.Array2 {
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

func createPiece() *proto.Array2 {
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

func collide(s *proto.UpdateState) bool {
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

func merge(s *proto.UpdateState) {
	player, arena := s.Player, s.Arena
	for y, row := range player.Matrix.Rows {
		for x, value := range row.Values {
			if value != 0 {
				arena.Matrix.Rows[y+int(player.Pos.Y)].Values[x+int(player.Pos.X)] = value
			}
		}
	}
}

func reset(s *proto.UpdateState) {
	player, arena := s.Player, s.Arena
	player.Matrix = createPiece()
	player.Pos = &proto.Pos{
		X: uint32(len(arena.Matrix.Rows[0].Values)/2 - len(player.Matrix.Rows[0].Values)/2),
		Y: 0,
	}
	if collide(s) {
		s.End = true
	}
}
