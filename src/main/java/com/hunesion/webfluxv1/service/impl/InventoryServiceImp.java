package com.hunesion.webfluxv1.service.impl;

import com.hunesion.webfluxv1.model.entity.Inventory;
import com.hunesion.webfluxv1.repository.InventoryRepository;
import com.hunesion.webfluxv1.service.InventoryService;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.r2dbc.core.R2dbcEntityTemplate;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;

import java.util.UUID;

@Service
@RequiredArgsConstructor
public class InventoryServiceImp implements InventoryService {

    private final InventoryRepository inventoryRepository;
    private final R2dbcEntityTemplate entityTemplate;

    @Override
    public Mono<Inventory> createInventoryForUser(UUID userId) {
        return entityTemplate.insert(
                Inventory.builder()
                        .id(UUID.randomUUID())
                        .userId(userId)
                        .capacity(10)
                        .build()
        );
    }
}
