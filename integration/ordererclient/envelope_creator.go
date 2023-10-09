/*
Copyright IBM Corp All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package ordererclient

import (
	"github.com/hyperledger/fabric/protoutil"
	. "github.com/onsi/gomega"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/integration/nwo"
)

func CreateBroadcastEnvelope(n *nwo.Network, entity interface{}, channel string, data []byte) *common.Envelope {
	var signer *nwo.SigningIdentity
	switch creator := entity.(type) {
	case *nwo.Peer:
		signer = n.PeerUserSigner(creator, "Admin")
	case *nwo.Orderer:
		signer = n.OrdererUserSigner(creator, "Admin")
	}
	Expect(signer).NotTo(BeNil())

	env, err := protoutil.CreateSignedEnvelope(
		//common.HeaderType_MESSAGE,
		common.HeaderType_ENDORSER_TRANSACTION,
		channel,
		signer,
		&common.Envelope{Payload: data},
		0,
		0,
	)
	Expect(err).NotTo(HaveOccurred())

	return env
}
