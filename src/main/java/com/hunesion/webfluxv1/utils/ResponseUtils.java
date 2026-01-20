package com.hunesion.webfluxv1.utils;

import com.hunesion.webfluxv1.model.response.ApiResponse;
import org.springframework.http.HttpStatusCode;
import org.springframework.http.ResponseEntity;
import reactor.core.publisher.Flux;
import reactor.core.publisher.Mono;

import java.util.List;

public final class ResponseUtils {

    public static  <T> Mono<ResponseEntity<ApiResponse<T>>> responseSingle(String message, HttpStatusCode httpStatusCode, Mono<T> payload) {
        return payload.map(data -> createResponseEntity(message,httpStatusCode,data));
    }

    public static <T> Mono<ResponseEntity<ApiResponse<List<T>>>> responseMulti(String message, HttpStatusCode httpStatusCode, Flux<T> data) {
        return data.collectList().map(list -> createResponseEntity(message,httpStatusCode,list));
    }

    private static <T> ResponseEntity<ApiResponse<T>> createResponseEntity(
            String message,
            HttpStatusCode httpStatusCode,
            T payload) {
        ApiResponse<T> apiResponse = ApiResponse.<T>builder()
                .message(message)
                .code(httpStatusCode)
                .payload(payload)
                .build();
        return new ResponseEntity<>(apiResponse, httpStatusCode);
    }
}
