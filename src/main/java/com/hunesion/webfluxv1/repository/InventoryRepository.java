package com.hunesion.webfluxv1.repository;

import com.hunesion.webfluxv1.model.entity.Inventory;
import org.springframework.data.r2dbc.repository.Query;
import org.springframework.data.r2dbc.repository.R2dbcRepository;
import org.springframework.stereotype.Repository;
import reactor.core.publisher.Mono;

import java.util.UUID;

@Repository
public interface InventoryRepository extends R2dbcRepository<Inventory, UUID> {

    // In InventoryRepository
//    @Query("INSERT INTO inventories (id, user_id, capacity) " +
//            "VALUES (:id, :userId, :capacity) " +
//            "RETURNING *")
//    Mono<Inventory> insertInventory(UUID id, UUID userId, Integer capacity);

}
