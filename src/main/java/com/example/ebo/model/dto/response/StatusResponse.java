package com.example.ebo.model.dto.response;

import lombok.Data;

import java.time.Instant;
import java.util.UUID;

@Data
public class StatusResponse {
    private UUID playerId;
    private int level;
    private long exp;
    private int hp;
    private Instant updatedAt;
}
