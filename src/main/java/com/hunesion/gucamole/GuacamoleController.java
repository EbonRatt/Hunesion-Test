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
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.time.LocalTime;

@RestController
public class GuacamoleController extends GuacamoleHTTPTunnelServlet {

    @RequestMapping(path = "/tunnel", method = {RequestMethod.GET, RequestMethod.POST, RequestMethod.OPTIONS})
    public void tunnel(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {
        super.handleTunnelRequest(request, response);
    }

    private static final Logger logger = LoggerFactory.getLogger(GuacamoleController.class);

    // Time-based access control (24-hour format)
    private static final LocalTime ALLOWED_START = LocalTime.of(13, 0);  // 1:00 PM
    private static final LocalTime ALLOWED_END = LocalTime.of(15, 0);    // 2:00 PM

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

        try {
            GuacamoleConfiguration config = new GuacamoleConfiguration();
            config.setProtocol("ssh");
            config.setParameter("hostname", "192.168.230.128");
            config.setParameter("port", "22");
            config.setParameter("username", "ebon");
            config.setParameter("password", "600");

            logger.info("Connecting to guacd at localhost:4822");

            GuacamoleSocket socket = new ConfiguredGuacamoleSocket(
                    new InetGuacamoleSocket("localhost", 4822),
                    config
            );

            logger.info("Connection successful!");
            GuacamoleTunnel tunnel = new SimpleGuacamoleTunnel(socket);
            return new LoggingGuacamoleTunnel(tunnel, logger);

        } catch (Exception e) {
            logger.error("Connection failed: ", e);
            throw new GuacamoleException("Failed to connect", e);
        }
    }

}
