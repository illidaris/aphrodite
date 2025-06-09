package idsnow

import (
	"time"
)

var (
	defTime = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) // 起始时间
)

// MachineID: NodeId 2^8 FrameId 2^2 Gene2 2^2 Gene 2^4
// func (i *IdGenerate) NewItem() {
// 	for gene := 0; gene < 32; gene++ {
// 		for gene2 := 0; gene2 < 4; gene2++ {
// 			for frm := 0; frm < 4; frm++ {
// 				sf, err := New(Settings{
// 					BitsMachineID: 16,
// 					BitsSequence:  8,
// 					StartTime:     defTime,
// 					TimeUnit:      1e7, // nsec, i.e. 10 msec
// 					MachineID: func() (int, error) {
// 						return int(uint32(i.NodeId)<<8 |
// 							uint32(frm)<<2 |
// 							uint32(gene2)<<2 |
// 							uint32(gene)), nil
// 					},
// 				})
// 				if err != nil {
// 					println(err)
// 				}
// 				i.SonyflakeMap[uint8(gene)][uint8(frm)] = sf
// 			}
// 		}
// 	}
// }

// func (i *IdGenerate) NewIDX(ctx context.Context, key string) uint64 {

// }
// func (i *IdGenerate) NewID(ctx context.Context, key string) (uint64, error) {

// }
// func (i *IdGenerate) NewIDIterate(ctx context.Context, iterate func(uint64), key string, opts ...dep.Option) error {

// }
