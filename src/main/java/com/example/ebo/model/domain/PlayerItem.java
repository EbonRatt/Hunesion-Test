package com.example.ebo.model.domain;

import lombok.Data;

import java.time.Instant;
import java.util.UUID;

@Data
public class PlayerItem {
    UUID playerId;
    Long itemId;
    Integer quantity;
    Instant acquiredAt;
}
