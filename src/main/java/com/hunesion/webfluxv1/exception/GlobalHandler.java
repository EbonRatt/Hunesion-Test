package com.hunesion.webfluxv1.exception;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.util.Map;

@RestControllerAdvice
class GlobalHandler {
    @ExceptionHandler(ItemsNotFoundException.class)
    public ResponseEntity<?> handle(ItemsNotFoundException ex) {
        return ResponseEntity.status(404).body(Map.of("message", ex.getMessage()));
    }
}
