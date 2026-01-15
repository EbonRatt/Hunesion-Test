"use client";

import { useEffect, useRef, useState } from "react";

export default function Page() {
  const [connecting, setConnecting] = useState(false);
  const [connected, setConnected] = useState(false);
  const [status, setStatus] = useState("Disconnected");
  const displayRef = useRef(null);
  const clientRef = useRef(null);
  const keyboardRef = useRef(null);
  const [formData, setFormData] = useState({
    protocol: "ssh",
    ip: "",
    username: "",
    password: "",
  });

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleConnect = async () => {
    if (!formData.ip || !formData.username || !formData.password) {
      setStatus("Please fill in all fields");
      return;
    }

    setConnecting(true);
    try {
      setStatus("Loading...");
      const Guacamole =
        (await import("guacamole-common-js")).default ||
        (await import("guacamole-common-js"));

      setStatus("Connecting...");
      const tunnel = new Guacamole.HTTPTunnel("http://localhost:8080/tunnel", true);
      const client = new Guacamole.Client(tunnel);
      clientRef.current = client;

      client.onerror = (err) => {
        setStatus("Error: " + (err?.message || "Connection failed"));
      };

      client.onstatechange = (state) => {
        const states = [
          "Idle",
          "Connecting",
          "Waiting",
          "Connected",
          "Disconnecting",
          "Disconnected",
        ];
        setStatus(states[state] || String(state));
      };

      setConnected(true);
      await new Promise((resolve) => setTimeout(resolve, 50));

      if (displayRef.current) {
        displayRef.current.innerHTML = "";
        const el = client.getDisplay().getElement();
        el.tabIndex = 0;
        el.style.outline = "none";
        el.addEventListener("click", () => el.focus());
        displayRef.current.appendChild(el);

        const keyboard = new Guacamole.Keyboard(document);
        keyboard.onkeydown = (k) => client.sendKeyEvent(1, k);
        keyboard.onkeyup = (k) => client.sendKeyEvent(0, k);
        keyboardRef.current = keyboard;
      }

      client.connect(
        `protocol=ssh&hostname=${encodeURIComponent(formData.ip)}&username=${encodeURIComponent(formData.username)}&password=${encodeURIComponent(formData.password)}&width=1280&height=720&dpi=96`
      );
    } catch (e) {
      setStatus("Error: " + e.message);
      setConnected(false);
    } finally {
      setConnecting(false);
    }
  };

  const handleDisconnect = () => {
    try {
      if (keyboardRef.current) {
        keyboardRef.current.onkeydown = null;
        keyboardRef.current.onkeyup = null;
        keyboardRef.current = null;
      }
      if (clientRef.current) {
        clientRef.current.disconnect();
        clientRef.current = null;
      }
      if (displayRef.current) {
        displayRef.current.innerHTML = "";
      }
    } catch {}
    setConnected(false);
    setStatus("Disconnected");
  };

  useEffect(() => {
    if (!connected) return;

    const el = displayRef.current;
    const client = clientRef.current;
    if (!el || !client || typeof client.sendSize !== "function") return;

    const ro = new ResizeObserver(() => {
      const rect = el.getBoundingClientRect();
      const w = Math.max(1, Math.floor(rect.width));
      const h = Math.max(1, Math.floor(rect.height));
      client.sendSize(w, h);
    });

    ro.observe(el);
    return () => ro.disconnect();
  }, [connected]);

  return connected ? (
    <div className="min-h-screen bg-gray-900 text-white flex flex-col">
      <div className="bg-gray-800 border-b border-gray-700 px-4 py-3 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-green-500 rounded-lg flex items-center justify-center">
            <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3" />
            </svg>
          </div>
          <div className="text-sm">
            <div className="font-semibold">SSH Session</div>
            <div className="text-gray-300">{status}</div>
          </div>
        </div>

        <button
          onClick={handleDisconnect}
          className="px-4 py-2 bg-red-500 hover:bg-red-600 text-white text-sm font-medium rounded-lg transition-colors"
          type="button"
        >
          Disconnect
        </button>
      </div>

      <div ref={displayRef} className="flex-1 w-full bg-black overflow-hidden" />
    </div>
  ) : (
    <div className="min-h-screen bg-gray-900 text-white flex items-center justify-center">
      <div className="bg-gray-800 rounded-xl shadow-2xl p-8 w-full max-w-md border border-gray-700">
        <div className="text-center mb-6">
          <div className="w-16 h-16 bg-gradient-to-br from-green-400 to-green-600 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </div>
          <h2 className="text-2xl font-bold text-white">Connect to Remote</h2>
          <p className="text-gray-400 mt-1">Enter your connection details below</p>
        </div>

        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1.5">
              Protocol
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg className="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <select
                name="protocol"
                value={formData.protocol}
                onChange={handleInputChange}
                className="w-full pl-10 pr-4 py-2.5 bg-gray-700 border border-gray-600 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent transition appearance-none cursor-pointer"
              >
                <option value="ssh">SSH (Secure Shell)</option>
              </select>
              <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
                <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              </div>
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1.5">
              IP Address / Hostname
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg className="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                </svg>
              </div>
              <input
                type="text"
                name="ip"
                value={formData.ip}
                onChange={handleInputChange}
                placeholder="192.168.1.100"
                className="w-full pl-10 pr-4 py-2.5 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent transition"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1.5">
              Username
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg className="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <input
                type="text"
                name="username"
                value={formData.username}
                onChange={handleInputChange}
                placeholder="admin"
                className="w-full pl-10 pr-4 py-2.5 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent transition"
              />
            </div>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-1.5">
              Password
            </label>
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg className="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <input
                type="password"
                name="password"
                value={formData.password}
                onChange={handleInputChange}
                placeholder="••••••••"
                className="w-full pl-10 pr-4 py-2.5 bg-gray-700 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent transition"
              />
            </div>
          </div>

          <button
            onClick={handleConnect}
            disabled={connecting}
            className="w-full mt-2 py-3 px-4 bg-gradient-to-r from-green-500 to-green-600 hover:from-green-600 hover:to-green-700 disabled:from-gray-600 disabled:to-gray-600 text-white font-semibold rounded-lg transition-all duration-200 flex items-center justify-center gap-2"
            type="button"
          >
            {connecting ? (
              <>
                <svg className="animate-spin h-5 w-5" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Connecting...
              </>
            ) : (
              <>
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                Connect
              </>
            )}
          </button>

          <div className="text-xs text-gray-400">{status}</div>
        </div>
      </div>
    </div>
  );
}
