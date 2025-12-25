package com.example.ebo.model.domain;

import lombok.Data;

import java.time.Instant;
import java.util.UUID;

@Data
public class Item {
    UUID id;
    String code;
    String name;
    String rarity;
    Instant createdAt;
}
