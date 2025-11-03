export interface ServerConfig {
  id: string;
  name: string;
  type: "vanilla" | "paper" | "fabric";
  version: string;
  port: number;
  memoryMb: number;
  path: string;
  eula: boolean;
  jarUrl?: string;
}

export type ServerState = "stopped" | "running" | "starting" | "crashed";

export interface PlayerInfo {
  current: number;
  max: number;
}

export interface ServerInfo {
  config: ServerConfig;
  state: ServerState;
  pid: number;
  uptimeSec: number;
  lastExitErr: string;
  players?: PlayerInfo;
}

export interface CreateServerRequest {
  name: string;
  type: "vanilla" | "paper" | "fabric";
  version: string;
  port: number;
  memoryMb: number;
  eula: boolean;
}

export interface LogEvent {
  type: string;
  serverId: string;
  data?: {
    stream: "stdout" | "stderr";
    line: string;
  };
}

export interface APIResponse<T = any> {
  error?: string;
  data?: T;
}
