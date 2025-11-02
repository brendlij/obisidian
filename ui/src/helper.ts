/**
 * MCS Manager API Helper Functions
 * Simple TypeScript helper to interact with the MCS Manager API
 *
 * Usage:
 * import { startServer, stopServer, createServer } from './helper'
 *
 * await createServer({
 *   name: 'My Server',
 *   type: 'paper',
 *   version: '1.21.10',
 *   port: 0,
 *   memoryMb: 2048,
 *   eula: true
 * })
 */

import type { ServerInfo, CreateServerRequest, LogEvent } from "./types";

const API_URL = "http://localhost:8484";

/**
 * Make an API request
 */
async function apiRequest<T = any>(
  path: string,
  method: string = "GET",
  body?: any
): Promise<T> {
  const url = `${API_URL}${path}`;
  const options: RequestInit = { method };

  if (body) {
    options.headers = { "Content-Type": "application/json" };
    options.body = JSON.stringify(body);
  }

  const response = await fetch(url, options);
  if (!response.ok) {
    throw new Error(`API Error: ${response.status} ${response.statusText}`);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
}

/**
 * Get all servers
 */
export async function listServers(): Promise<ServerInfo[]> {
  return apiRequest("/servers");
}

/**
 * Get a specific server by ID
 */
export async function getServer(id: string): Promise<ServerInfo> {
  return apiRequest(`/servers/${id}`);
}

/**
 * Create a new server
 */
export async function createServer(
  config: CreateServerRequest
): Promise<ServerInfo> {
  return apiRequest("/servers", "POST", config);
}

/**
 * Start a server
 */
export async function startServer(id: string): Promise<ServerInfo> {
  return apiRequest(`/servers/${id}/start`, "POST");
}

/**
 * Stop a server
 */
export async function stopServer(id: string): Promise<ServerInfo> {
  return apiRequest(`/servers/${id}/stop`, "POST");
}

/**
 * Delete a server (must be stopped)
 */
export async function deleteServer(id: string): Promise<void> {
  return apiRequest(`/servers/${id}`, "DELETE");
}

/**
 * Send a command to a running server
 */
export async function sendCommand(id: string, command: string): Promise<void> {
  return apiRequest(`/servers/${id}/cmd`, "POST", { command });
}

/**
 * Get server logs (last 200 lines)
 */
export async function getServerLogs(id: string): Promise<string> {
  return apiRequest(`/servers/${id}/logs`, "GET");
}

/**
 * Subscribe to server logs via EventSource
 * Returns an EventSource - remember to close it when done!
 *
 * Usage:
 * const eventSource = subscribeToLogs('my-server-id')
 * eventSource.onmessage = (event) => {
 *   const log = JSON.parse(event.data)
 *   console.log(log)
 * }
 * // Later:
 * eventSource.close()
 */
export function subscribeToLogs(
  serverId: string,
  onMessage?: (event: LogEvent) => void,
  onError?: (error: Event) => void
): EventSource {
  const eventSource = new EventSource(`${API_URL}/events`);

  eventSource.onmessage = (event) => {
    try {
      const data: LogEvent = JSON.parse(event.data);
      if (data.serverId === serverId) {
        onMessage?.(data);
      }
    } catch (e) {
      console.error("Failed to parse log event:", e);
    }
  };

  eventSource.onerror = (error) => {
    onError?.(error);
  };

  return eventSource;
}

/**
 * Helper to format server uptime
 */
export function formatUptime(seconds: number): string {
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  if (hours > 0) return `${hours}h ${minutes}m`;
  if (minutes > 0) return `${minutes}m ${secs}s`;
  return `${secs}s`;
}

/**
 * Helper to get status emoji
 */
export function getStatusEmoji(state: string): string {
  switch (state) {
    case "running":
      return "🟢";
    case "stopped":
      return "🔴";
    case "starting":
      return "🟡";
    case "crashed":
      return "💥";
    default:
      return "❓";
  }
}

/**
 * Create a server with sensible defaults
 */
export async function quickCreateServer(
  name: string,
  version: string = "1.21.10",
  type: "vanilla" | "paper" = "paper"
): Promise<ServerInfo> {
  return createServer({
    name,
    type,
    version,
    port: 0, // Auto-assign
    memoryMb: 2048,
    eula: true,
  });
}

/**
 * Quick start/stop helper
 */
export async function quickStart(id: string): Promise<void> {
  await startServer(id);
}

export async function quickStop(id: string): Promise<void> {
  await stopServer(id);
}
