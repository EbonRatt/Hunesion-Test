package com.hunesion.webfluxv1.model.response;

import com.hunesion.webfluxv1.model.enums.Rarity;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class ItemResponse {
    UUID id;
    String name;
    String description;
    double value;
    Rarity rarity;
    LocalDateTime createdAt;
    LocalDateTime updatedAt;
}
