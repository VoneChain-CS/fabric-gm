package global

import (
	"sync"

	"github.com/VoneChain-CS/fabric-gm/common/flogging"
	"github.com/VoneChain-CS/fabric-gm/core/chaincode"
	"github.com/VoneChain-CS/fabric-gm/core/chaincode/persistence"
)

var (
	G      Global
	logger = flogging.MustGetLogger("global")
)

type Global struct {
	CCManage CCManage
}

type CCManage struct {
	laucher *chaincode.RuntimeLauncher
	store   *persistence.Store
	mu      sync.Mutex
}

func (ccm *CCManage) SetLaucher(laucher *chaincode.RuntimeLauncher) {
	ccm.mu.Lock()
	defer ccm.mu.Unlock()
	ccm.laucher = laucher
}

func (ccm *CCManage) SetStore(store *persistence.Store) {
	ccm.mu.Lock()
	defer ccm.mu.Unlock()
	ccm.store = store
}

func (ccm *CCManage) Uninstall(packageID string) error {
	ccm.mu.Lock()
	defer ccm.mu.Unlock()
	if err := ccm.store.Delete(packageID); err != nil {
		return err
	}
	return ccm.laucher.Stop(packageID)
}

func (ccm *CCManage) LoadPackage(packageID string) ([]byte, error) {
	ccm.mu.Lock()
	defer ccm.mu.Unlock()
	return ccm.store.Load(packageID)
}
