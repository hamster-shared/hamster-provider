export interface VmConfig {
  cpu: number;
  mem: number;
  disk: number;
  system: string;
  image: string;
  accessPort: number;
  // 虚拟化类型，docker/kvm
  type: string;
}

export interface ProviderConfig {
  chainApi: string;
  seedOrPhrase: string;
  vm: VmConfig | Recordable;
  bootstraps: string[];
  publicIP: string;
}
