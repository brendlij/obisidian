import { spawn } from "bun";
import process from "process";

console.log("🚀 Starting MCS Manager (Backend + Frontend)...\n");

// Spawn Backend (Go)
const backend = spawn({
  cmd: ["go", "run", "./cmd/mcs-manager"],
  cwd: "./mcs-manager",
  stdio: ["inherit", "inherit", "inherit"],
});

// Spawn Frontend (Vite)
const frontend = spawn({
  cmd: ["bun", "run", "dev"],
  cwd: "./ui",
  stdio: ["inherit", "inherit", "inherit"],
});

console.log("✅ Backend (Port 8484) & Frontend (Port 5173) started!");
console.log("📍 Frontend: http://localhost:5173");
console.log("📍 Backend: http://localhost:8484\n");

// Graceful Shutdown
process.on("SIGINT", async () => {
  console.log("\n\n🛑 Shutting down gracefully...");
  backend.kill();
  frontend.kill();
  process.exit(0);
});
