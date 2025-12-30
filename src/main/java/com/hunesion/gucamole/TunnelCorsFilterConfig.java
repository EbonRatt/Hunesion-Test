package com.hunesion.gucamole;

import jakarta.servlet.*;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.io.IOException;

// Cuz Of GuacamoleHTTPTunnelServlet Not MVC so Config global Not work in
// So Create Filter To Handle CORS

@Configuration
public class TunnelCorsFilterConfig {

    @Bean
    public FilterRegistrationBean<Filter> tunnelCorsFilter() {
        FilterRegistrationBean<Filter> bean = new FilterRegistrationBean<>();

        bean.setFilter(new Filter() {
            @Override
            public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
                    throws IOException, ServletException {

                HttpServletRequest req = (HttpServletRequest) request;
                HttpServletResponse res = (HttpServletResponse) response;

                String origin = req.getHeader("Origin");
                if ("http://localhost:3000".equals(origin)) {
                    res.setHeader("Access-Control-Allow-Origin", origin);
                    res.setHeader("Vary", "Origin");
                    res.setHeader("Access-Control-Allow-Credentials", "true");
                    res.setHeader("Access-Control-Allow-Methods", "GET,POST,OPTIONS");
                    res.setHeader("Access-Control-Allow-Headers", req.getHeader("Access-Control-Request-Headers"));

                    // Guacamole-specific headers your JS must be allowed to read:
                    res.setHeader("Access-Control-Expose-Headers",
                            "Guacamole-Tunnel-Token,Guacamole-Status-Code,Guacamole-Error-Message");
                }

                // Preflight
                if ("OPTIONS".equalsIgnoreCase(req.getMethod())) {
                    res.setStatus(204);
                    return;
                }

                chain.doFilter(request, response);
            }
        });

        bean.addUrlPatterns("/tunnel", "/tunnel/*");
        bean.setOrder(0); // run early
        return bean;
    }
}
