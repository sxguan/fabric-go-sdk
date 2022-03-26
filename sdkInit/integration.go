package sdkInit

import (
	"encoding/hex"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
)

func DiscoverLocalPeers(ctxProvider contextAPI.ClientProvider, expectedPeers int) ([]fabAPI.Peer, error) {
	ctx, err := contextImpl.NewLocal(ctxProvider)
	if err != nil {
		return nil, fmt.Errorf("error creating local context: %v", err)
	}

	discoveredPeers, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			peers, serviceErr := ctx.LocalDiscoveryService().GetPeers()
			if serviceErr != nil {
				return nil, fmt.Errorf("getting peers for MSP [%s] error: %v", ctx.Identifier().MSPID, serviceErr)
			}
			if len(peers) < expectedPeers {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Expecting %d peers but got %d", expectedPeers, len(peers)), nil)
			}
			return peers, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return discoveredPeers.([]fabAPI.Peer), nil
}
func (t *SdkEnvInfo) InitService(chaincodeID, channelID string, org *OrgInfo, sdk *fabsdk.FabricSDK) error {
	handler := &SdkEnvInfo{
		ChaincodeID: chaincodeID,
	}
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(org.OrgUser), fabsdk.WithOrg(org.OrgName))
	// Channel client is used to query and execute transactions (Org1 is default org)
	var err error
	t.ChClient, err = channel.New(clientChannelContext)
	if err != nil {
		return err
	}
	t.EvClient, err = event.New(clientChannelContext, event.WithBlockEvents())
	if err != nil {
		return err
	}
	handler.ChClient = t.ChClient
	handler.EvClient = t.EvClient
	return nil
}

func regitserEvent(client *event.Client, chaincodeID string) (fabAPI.Registration, <-chan *fabAPI.CCEvent) {
	eventName := "chaincode-event"

	reg, notifier, err := client.RegisterChaincodeEvent(chaincodeID, eventName)
	if err != nil {
		fmt.Println("注册链码事件失败: %s", err)
	}

	return reg, notifier
}
func ChainCodeEventListener(c *event.Client, ccID string) fabAPI.Registration {

	reg, notifier := regitserEvent(c, ccID)

	// consume event
	go func() {
		for e := range notifier {
			log.Printf("Receive cc event, ccid: %v \neventName: %v\n"+
				"payload: %v \ntxid: %v \nblock: %v \nsourceURL: %v\n",
				e.ChaincodeID, e.EventName, string(e.Payload), e.TxID, e.BlockNumber, e.SourceURL)
		}
	}()

	return reg
}

func TxListener(c *event.Client, txIDCh chan string) {
	log.Println("Transaction listener start")
	defer log.Println("Transaction listener exit")

	for id := range txIDCh {
		// Register monitor transaction event
		log.Printf("Register transaction event for: %v", id)
		txReg, txCh, err := c.RegisterTxStatusEvent(id)
		if err != nil {
			log.Printf("Register transaction event error: %v", err)
			continue
		}
		defer c.Unregister(txReg)

		// Receive transaction event
		go func() {
			for e := range txCh {
				log.Printf("Receive transaction event: txid: %v, "+
					"validation code: %v, block number: %v",
					e.TxID,
					e.TxValidationCode,
					e.BlockNumber)
			}
		}()
	}
}

func BlockListener(ec *event.Client) fabAPI.Registration {
	// Register monitor block event
	beReg, beCh, err := ec.RegisterBlockEvent()
	if err != nil {
		log.Printf("Register block event error: %v", err)
	}
	log.Println("Registered block event")

	// Receive block event
	go func() {
		for e := range beCh {
			log.Printf("Receive block event:\nSourceURL: %v\nNumber: %v\nHash"+
				": %v\nPreviousHash: %v\n\n",
				e.SourceURL,
				e.Block.Header.Number,
				hex.EncodeToString(e.Block.Header.DataHash),
				hex.EncodeToString(e.Block.Header.PreviousHash))
		}
	}()

	return beReg
}
