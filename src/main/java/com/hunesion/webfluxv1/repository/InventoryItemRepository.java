package com.hunesion.webfluxv1.repository;

import com.hunesion.webfluxv1.model.entity.InventoryItem;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface InventoryItemRepository extends R2dbcRepository<InventoryItem, UUID> {
}
