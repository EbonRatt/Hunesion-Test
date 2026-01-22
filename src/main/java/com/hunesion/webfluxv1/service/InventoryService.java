package com.hunesion.webfluxv1.service;

import com.hunesion.webfluxv1.model.entity.Inventory;
import reactor.core.publisher.Mono;

import java.util.UUID;

public interface InventoryService {
    Mono<Inventory> createInventoryForUser(UUID userId);

}
