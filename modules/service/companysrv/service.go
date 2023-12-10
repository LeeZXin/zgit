package companysrv

import (
	"context"
	"zgit/modules/model/companymd"
)

func GetCompanyInfoByCompanyId(ctx context.Context, id string) (companymd.CompanyInfo, bool, error) {
	return companymd.CompanyInfo{}, true, nil
}
