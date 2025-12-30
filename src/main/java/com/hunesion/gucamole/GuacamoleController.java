package com.hunesion.gucamole;

import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.apache.guacamole.GuacamoleException;
import org.apache.guacamole.net.GuacamoleSocket;
import org.apache.guacamole.net.GuacamoleTunnel;
import org.apache.guacamole.net.InetGuacamoleSocket;
import org.apache.guacamole.net.SimpleGuacamoleTunnel;
import org.apache.guacamole.protocol.ConfiguredGuacamoleSocket;
import org.apache.guacamole.protocol.GuacamoleConfiguration;
import org.apache.guacamole.servlet.GuacamoleHTTPTunnelServlet;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;

@RestController
public class GuacamoleController extends GuacamoleHTTPTunnelServlet {

    @RequestMapping(path = "/tunnel", method = RequestMethod.OPTIONS)
    public void tunnel(HttpServletRequest request, HttpServletResponse response)
            throws ServletException, IOException {
        super.handleTunnelRequest(request, response);
    }

    private static final Logger logger = LoggerFactory.getLogger(GuacamoleController.class);

    @Override
    protected GuacamoleTunnel doConnect(HttpServletRequest request) throws GuacamoleException {

        logger.info("Tunnel connection requested!");
        logger.info("Remote address: {}", request.getRemoteAddr());

        try {
            GuacamoleConfiguration config = new GuacamoleConfiguration();
            config.setProtocol("ssh");
            config.setParameter("hostname", "192.168.230.128");
            config.setParameter("port", "22");
            config.setParameter("username", "ebon");
            config.setParameter("password", "123");

            logger.info("Connecting to guacd at localhost:4822");

            GuacamoleSocket socket = new ConfiguredGuacamoleSocket(
                    new InetGuacamoleSocket("localhost", 4822),
                    config
            );

            logger.info("Connection successful!");
            return new SimpleGuacamoleTunnel(socket);

        } catch (Exception e) {
            logger.error("Connection failed: ", e);
            throw new GuacamoleException("Failed to connect", e);
        }
    }

}
