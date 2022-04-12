package uuid

import "github.com/google/uuid"

//
// GetUUid
//  @Description: 一定会获得uuid 如果获取失败会重新获得
//  @return uuid.UUID
//
func GetUUid() uuid.UUID {
	for {
		serviceID, err := uuid.NewRandom()
		if err == nil {
			return serviceID
		}
	}
}
