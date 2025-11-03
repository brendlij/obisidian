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
 * Restart a server
 */
export async function restartServer(id: string): Promise<ServerInfo> {
  return apiRequest(`/servers/${id}/restart`, "POST");
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
 * Subscribe to server events (logs, state changes, etc.) via EventSource
 * Returns an EventSource - remember to close it when done!
 *
 * Usage:
 * const eventSource = subscribeToServerEvents('my-server-id', {
 *   onLog: (data) => console.log(data),
 *   onStateChange: () => refreshServerData(),
 *   onError: (err) => console.error(err)
 * })
 */
export function subscribeToServerEvents(
  serverId: string,
  options?: {
    onLog?: (event: any) => void;
    onStateChange?: (type: string, event: any) => void;
    onError?: (error: Event) => void;
  }
): EventSource {
  const eventSource = new EventSource(`${API_URL}/events`);

  console.log("[subscribeToServerEvents] Connected for server:", serverId);

  // Handle default message events
  eventSource.onmessage = (event) => {
    console.log("[subscribeToServerEvents] onmessage:", event.data);
    try {
      const data: any = JSON.parse(event.data);
      if (data.serverId === serverId) {
        options?.onLog?.(data);
      }
    } catch (e) {
      console.error("[subscribeToServerEvents] Parse error:", e);
    }
  };

  // Handle typed events (server.log, server.started, etc.)
  const eventHandler = (event: Event) => {
    const messageEvent = event as MessageEvent;
    console.log(
      `[subscribeToServerEvents] Event (${event.type}):`,
      messageEvent.data
    );
    try {
      const data: any = JSON.parse(messageEvent.data);

      if (data.serverId === serverId) {
        console.log("[subscribeToServerEvents] Server ID matched:", serverId);

        if (event.type === "server.log") {
          console.log("[subscribeToServerEvents] Calling onLog");
          options?.onLog?.(data);
        } else if (
          [
            "server.started",
            "server.stopped",
            "server.exited",
            "server.crashed",
          ].includes(event.type)
        ) {
          console.log(
            "[subscribeToServerEvents] Calling onStateChange:",
            event.type
          );
          options?.onStateChange?.(event.type, data);
        }
      } else {
        console.log(
          "[subscribeToServerEvents] Server ID mismatch:",
          data.serverId,
          "expected:",
          serverId
        );
      }
    } catch (e) {
      console.error("[subscribeToServerEvents] Parse error:", e);
    }
  };

  // Register listeners for all event types
  eventSource.addEventListener("server.log", eventHandler);
  eventSource.addEventListener("server.started", eventHandler);
  eventSource.addEventListener("server.stopped", eventHandler);
  eventSource.addEventListener("server.exited", eventHandler);
  eventSource.addEventListener("server.crashed", eventHandler);

  eventSource.onerror = (error) => {
    console.error("[subscribeToServerEvents] EventSource error:", error);
    options?.onError?.(error);
  };

  return eventSource;
}

/**
 * Subscribe to server logs via EventSource
 * Returns an EventSource - remember to close it when done!
 *
 * @deprecated Use subscribeToServerEvents instead
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

// Test function - can be called from browser console
export function testSSE(serverId: string = "53c2035c8641"): void {
  console.log("🧪 [TEST SSE] Starting SSE test for server:", serverId);
  subscribeToServerEvents(serverId, {
    onLog: (event) => {
      console.log("📝 [SSE TEST] Log:", event);
    },
    onStateChange: (type, event) => {
      console.log("🔄 [SSE TEST] State change:", type, event);
    },
    onError: (error) => {
      console.error("❌ [SSE TEST] Error:", error);
    },
  });
  console.log("✅ [TEST SSE] Subscription created - check your logs!");
}

// Make it globally available
(globalThis as any).testSSE = testSSE;
