export interface ComputingResource {
  index: number,
  accountId: number[],
  peerId: string,
  config: VmConfig,
  rentalStatistics: RentalStatistics,
  rentalInfo: RentalInfo,
  status: ResourceStatus,
}

export interface ResourceStatus{
  isInuse: boolean,
  isLocked: boolean,
  isUnused: boolean,
  isOffline: boolean,
}

export interface VmConfig{
  cpu: number,
  memory: number,
  system: string,
  cpuModel: string,
}

export interface RentalStatistics {
  rentalCount: number,
  rentalDuration: number,
  faultCount: number,
  faultDuration: number,
}

export interface RentalInfo {
  rentUnitPrice: number,
  rentDuration: number,
  endOfRent: number,
}

export interface UnitPriceParam {
  unitPrice: number,
}
