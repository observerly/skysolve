/*****************************************************************************************************************/

//	@author		Michael Roberts <michael@observerly.com>
//	@package	@observerly/skysolve
//	@license	Copyright © 2021-2025 observerly

/*****************************************************************************************************************/

package index

/*****************************************************************************************************************/

import (
	"github.com/observerly/skysolve/pkg/catalog"
	"github.com/observerly/skysolve/pkg/healpix"
)

/*****************************************************************************************************************/

type Indexer struct {
	Catalog catalog.CatalogService
	HealPIX healpix.HealPIX
}

/*****************************************************************************************************************/

func NewIndexer(
	healpix healpix.HealPIX,
	catalog catalog.CatalogService,
) *Indexer {
	return &Indexer{
		Catalog: catalog,
		HealPIX: healpix,
	}
}

/*****************************************************************************************************************/
