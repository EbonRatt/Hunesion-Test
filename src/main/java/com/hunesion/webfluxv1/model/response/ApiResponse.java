package com.hunesion.webfluxv1.model.response;

import com.fasterxml.jackson.annotation.JsonInclude;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.http.HttpStatusCode;

import java.time.Instant;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class ApiResponse<T> {

    private HttpStatusCode code;
    private String message;

    @JsonInclude(JsonInclude.Include.NON_NULL)
    private T payload;

    @Builder.Default
    private Instant timestamps = Instant.now();

}
