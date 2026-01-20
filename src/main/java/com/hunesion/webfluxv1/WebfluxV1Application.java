package com.hunesion.webfluxv1;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.data.r2dbc.config.EnableR2dbcAuditing;

@SpringBootApplication
public class WebfluxV1Application {

    public static void main(String[] args) {
        SpringApplication.run(WebfluxV1Application.class, args);
    }

}
