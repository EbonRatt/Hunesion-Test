"use client";

import { useEffect, useRef, useState } from "react";

export default function Page() {
  const ref = useRef(null);
  const [status, setStatus] = useState("Init...");

  useEffect(() => {
    let client, keyboard;

    async function start() {
      setStatus("Loading guacamole-common-js...");
      const Guacamole = (await import("guacamole-common-js")).default || (await import("guacamole-common-js"));

      setStatus("Connecting...");
      const tunnel = new Guacamole.HTTPTunnel("http://localhost:8080/tunnel", true);
      client = new Guacamole.Client(tunnel);

      ref.current.innerHTML = "";
      ref.current.appendChild(client.getDisplay().getElement());

      keyboard = new Guacamole.Keyboard(document);
      keyboard.onkeydown = (k) => client.sendKeyEvent(1, k);
      keyboard.onkeyup = (k) => client.sendKeyEvent(0, k);

      client.connect();
      setStatus("Connected (or connecting...)");
    }

    start().catch((e) => setStatus("Error: " + e.message));

    return () => {
      try {
        if (keyboard) keyboard.onkeydown = keyboard.onkeyup = null;
        if (client) client.disconnect();
      } catch {}
    };
  }, []);

  return (
      <div style={{ height: "100vh", background: "#000", color: "#fff" }}>
        <div style={{ padding: 12, background: "#1f2937" }}><b>Status:</b> {status}</div>
        <div ref={ref} style={{ height: "calc(100vh - 48px)" }} />
      </div>
  );
}
