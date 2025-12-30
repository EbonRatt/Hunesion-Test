package com.hunesion.gucamole;

import org.springframework.boot.web.servlet.ServletRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class GuacServletConfig {

    @Bean
    public ServletRegistrationBean<GuacamoleController> guacTunnelServlet() {
        return new ServletRegistrationBean<>(
                new GuacamoleController(),
                "/tunnel",
                "/tunnel/*"
        );
    }
}

