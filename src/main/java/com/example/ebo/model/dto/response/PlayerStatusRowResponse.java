package com.example.ebo.model.dto.response;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.Instant;
import java.util.UUID;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class PlayerStatusRowResponse {
    private UUID id;
    private String username;
    private Instant createdAt;

    private Integer statusLevel;
    private Long statusExp;
    private Integer statusHp;
    private Instant statusUpdatedAt;

}
