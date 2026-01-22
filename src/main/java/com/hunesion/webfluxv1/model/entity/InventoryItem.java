package com.hunesion.webfluxv1.model.entity;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Table;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Table("inventory_items")
public class InventoryItem {
    @Id
    private UUID id;

    private UUID inventoryId;
    private UUID itemId;
    private int quantity;
    private LocalDateTime acquiredAt;
}
