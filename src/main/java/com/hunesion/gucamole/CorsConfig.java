//package com.hunesion.gucamole;
//
//import org.springframework.context.annotation.Configuration;
//import org.springframework.web.servlet.config.annotation.CorsRegistry;
//import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
//
//@Configuration
//public class CorsConfig implements WebMvcConfigurer {
//
//    @Override
//    public void addCorsMappings(CorsRegistry registry) {
//        registry.addMapping("/tunnel/**")
//                .allowedOrigins("http://localhost:3000")
//                .allowedMethods("GET", "POST", "OPTIONS")
//                .allowedHeaders("*")
//                // Guacamole tunnel uses these response headers; JS needs to read them:
//                .exposedHeaders("Guacamole-Tunnel-Token",
//                        "Guacamole-Status-Code",
//                        "Guacamole-Error-Message")
//                .allowCredentials(true);
//    }
//}
