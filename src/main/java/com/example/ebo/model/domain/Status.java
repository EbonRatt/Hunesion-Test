package com.example.ebo.model.domain;

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
public class Status {
    UUID playerId;
    Integer level;
    Long exp;
    Integer hp;
    Instant updatedAt;
}
