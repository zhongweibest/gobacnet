package gobacnet

import (
	"github.com/zhongweibest/gobacnet/encoding"
	bactype "github.com/zhongweibest/gobacnet/types"
)

func (c *Client) iAm(dest bactype.Address) error {
	enc := encoding.NewEncoder()
	enc.NPDU(
		bactype.NPDU{
			Version:               bactype.ProtocolVersion,
			Destination:           &dest,
			IsNetworkLayerMessage: false,
			ExpectingReply:        false,
			Priority:              bactype.Normal,
			HopCount:              bactype.DefaultHopCount,
		})

	//	iams := []bactype.ObjectID{bactype.ObjectID{Instance: 1, Type: 5}}
	//	enc.IAm(iams)
	_, err := c.send(dest, enc.Bytes())
	return err
}
