package factory

import (
	"fmt"
	"github.com/hamster-shared/hamster-provider/core/modules/provider"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/aptos"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/avalanche"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/bsc"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/ethereum"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/near"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/optimism"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/polygon"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/starkware"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/sui"
	"github.com/hamster-shared/hamster-provider/core/modules/provider/thegraph"
)

func GetChain(deployType uint32) (provider.Chain, error) {
	switch deployType {
	case 0:
		return thegraph.New(), nil
	case 1:
		return aptos.New(), nil
	case 2:
		return sui.New(), nil
	case 3:
		return ethereum.New(), nil
	case 4:
		return bsc.New(), nil
	case 5:
		return polygon.New(), nil
	case 6:
		return avalanche.New(), nil
	case 7:
		return optimism.New(), nil
	case 8:
		return nil, fmt.Errorf("not support deployType %d", deployType)
	case 9:
		return starkware.New(), nil
	case 10:
		return near.New(), nil
	case 11:
		return nil, fmt.Errorf("not support deployType %d", deployType)
	default:
		return nil, fmt.Errorf("not support deployType %d", deployType)
	}
}
