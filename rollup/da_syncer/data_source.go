package da_syncer

import (
	"context"

	"github.com/tenderly/net-scroll-geth/core"
	"github.com/tenderly/net-scroll-geth/ethdb"
	"github.com/tenderly/net-scroll-geth/params"
	"github.com/tenderly/net-scroll-geth/rollup/da_syncer/blob_client"
	"github.com/tenderly/net-scroll-geth/rollup/da_syncer/da"
	"github.com/tenderly/net-scroll-geth/rollup/rollup_sync_service"
)

type DataSource interface {
	NextData() (da.Entries, error)
	L1Height() uint64
}

type DataSourceFactory struct {
	config        Config
	genesisConfig *params.ChainConfig
	l1Client      *rollup_sync_service.L1Client
	blobClient    blob_client.BlobClient
	db            ethdb.Database
}

func NewDataSourceFactory(blockchain *core.BlockChain, genesisConfig *params.ChainConfig, config Config, l1Client *rollup_sync_service.L1Client, blobClient blob_client.BlobClient, db ethdb.Database) *DataSourceFactory {
	return &DataSourceFactory{
		config:        config,
		genesisConfig: genesisConfig,
		l1Client:      l1Client,
		blobClient:    blobClient,
		db:            db,
	}
}

func (ds *DataSourceFactory) OpenDataSource(ctx context.Context, l1height uint64) (DataSource, error) {
	return da.NewCalldataBlobSource(ctx, l1height, ds.l1Client, ds.blobClient, ds.db)
}
