package com.hunesion.gucamole;

import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.apache.guacamole.GuacamoleException;
import org.apache.guacamole.io.GuacamoleReader;
import org.apache.guacamole.io.GuacamoleWriter;
import org.apache.guacamole.net.GuacamoleSocket;
import org.apache.guacamole.net.GuacamoleTunnel;
import org.apache.guacamole.net.InetGuacamoleSocket;
import org.apache.guacamole.net.SimpleGuacamoleTunnel;
import org.apache.guacamole.protocol.ConfiguredGuacamoleSocket;
import org.apache.guacamole.protocol.GuacamoleConfiguration;
import org.apache.guacamole.protocol.GuacamoleInstruction;
import org.apache.guacamole.servlet.GuacamoleHTTPTunnelServlet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.time.LocalTime;

public class GuacamoleController extends GuacamoleHTTPTunnelServlet {

    private static final Logger logger = LoggerFactory.getLogger(GuacamoleController.class);

    // Time-based access control (24-hour format)
    private static final LocalTime ALLOWED_START = LocalTime.of(13, 0);  // 1:00 PM
    private static final LocalTime ALLOWED_END = LocalTime.of(17, 0);    // 2:00 PM

    @Override
    protected GuacamoleTunnel doConnect(HttpServletRequest request) throws GuacamoleException {

        logger.info("Tunnel connection requested!");
        logger.info("Remote address: {}", request.getRemoteAddr());

        // Check if current time is within allowed hours
        LocalTime now = LocalTime.now();
        if (now.isBefore(ALLOWED_START) || now.isAfter(ALLOWED_END)) {
            logger.warn("ACCESS DENIED: Current time {} is outside allowed hours ({} - {})",
                    now, ALLOWED_START, ALLOWED_END);
            throw new GuacamoleException("Access denied: Connection only allowed between " +
                    ALLOWED_START + " and " + ALLOWED_END);
        }
        logger.info("Time check passed: {} is within allowed hours", now);

        // Get dynamic connection parameters from frontend
        String protocol = request.getParameter("protocol");
        String hostname = request.getParameter("hostname");
        String username = request.getParameter("username");
        String password = request.getParameter("password");
        String port = request.getParameter("port");

        // Validate required parameters
        if (hostname == null || hostname.isEmpty()) {
            throw new GuacamoleException("Missing required parameter: hostname");
        }
        if (username == null || username.isEmpty()) {
            throw new GuacamoleException("Missing required parameter: username");
        }
        if (password == null || password.isEmpty()) {
            throw new GuacamoleException("Missing required parameter: password");
        }

        // Default protocol to SSH if not specified
        if (protocol == null || protocol.isEmpty()) {
            protocol = "ssh";
        }

        // Set default port based on protocol
        if (port == null || port.isEmpty()) {
            port = "rdp".equalsIgnoreCase(protocol) ? "3389" : "22";
        }

        logger.info("{} connection request - Host: {}, User: {}, Port: {}", protocol.toUpperCase(), hostname, username, port);

        try {
            GuacamoleConfiguration config = new GuacamoleConfiguration();
            config.setProtocol(protocol);
            config.setParameter("hostname", hostname);
            config.setParameter("port", port);
            config.setParameter("username", username);
            config.setParameter("password", password);

            // RDP-specific settings
            if ("rdp".equalsIgnoreCase(protocol)) {
                config.setParameter("security", "any");
                config.setParameter("ignore-cert", "true");
                config.setParameter("enable-wallpaper", "false");
                config.setParameter("enable-theming", "false");
                config.setParameter("enable-font-smoothing", "false");
            }

            logger.info("Connecting to guacd at localhost:4822");

            GuacamoleSocket socket = new ConfiguredGuacamoleSocket(
                    new InetGuacamoleSocket("localhost", 4822),
                    config
            );

            logger.info("Connection successful!");
            GuacamoleTunnel tunnel = new SimpleGuacamoleTunnel(socket);
//            return new LoggingGuacamoleTunnel(tunnel, logger);
            return tunnel;
        } catch (Exception e) {
            logger.error("Connection failed: ", e);
            throw new GuacamoleException("Failed to connect", e);
        }
    }

}
